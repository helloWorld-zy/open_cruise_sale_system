package notification

import (
	"backend/internal/domain"
	"backend/internal/service"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SMSSender sends notifications via SMS
type SMSSender struct {
	provider     string // aliyun, tencent, twilio
	accessKey    string
	accessSecret string
	signName     string
	templates    map[string]string // notification type -> template code
	httpClient   *http.Client
	endpoint     string
}

// AliyunSMSRequest represents Aliyun SMS API request
type AliyunSMSRequest struct {
	PhoneNumbers    string `json:"PhoneNumbers"`
	SignName        string `json:"SignName"`
	TemplateCode    string `json:"TemplateCode"`
	TemplateParam   string `json:"TemplateParam"`
	AccessKeyId     string `json:"AccessKeyId"`
	Timestamp       string `json:"Timestamp"`
	Signature       string `json:"Signature"`
	SignatureMethod string `json:"SignatureMethod"`
	SignatureNonce  string `json:"SignatureNonce"`
}

// AliyunSMSResponse represents Aliyun SMS API response
type AliyunSMSResponse struct {
	Code    string `json:"Code"`
	Message string `json:"Message"`
	BizId   string `json:"BizId"`
}

// SMSTemplateParam represents SMS template parameters
type SMSTemplateParam struct {
	Title   string `json:"title,omitempty"`
	Content string `json:"content,omitempty"`
	OrderNo string `json:"order_no,omitempty"`
	Amount  string `json:"amount,omitempty"`
	Count   string `json:"count,omitempty"`
	Days    string `json:"days,omitempty"`
}

// NewAliyunSMSSender creates a new Aliyun SMS sender
func NewAliyunSMSSender(accessKey, accessSecret, signName string) *SMSSender {
	templates := map[string]string{
		domain.NotificationTypeOrder:     "SMS_ORDER_TEMPLATE",
		domain.NotificationTypePayment:   "SMS_PAYMENT_TEMPLATE",
		domain.NotificationTypeRefund:    "SMS_REFUND_TEMPLATE",
		domain.NotificationTypeVoyage:    "SMS_VOYAGE_TEMPLATE",
		domain.NotificationTypeSystem:    "SMS_SYSTEM_TEMPLATE",
		domain.NotificationTypeInventory: "SMS_INVENTORY_TEMPLATE",
		domain.NotificationTypePromotion: "SMS_PROMOTION_TEMPLATE",
	}

	return &SMSSender{
		provider:     "aliyun",
		accessKey:    accessKey,
		accessSecret: accessSecret,
		signName:     signName,
		templates:    templates,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		endpoint:     "https://dysmsapi.aliyuncs.com",
	}
}

// NewTencentSMSSender creates a new Tencent Cloud SMS sender
func NewTencentSMSSender(accessKey, accessSecret, signName string) *SMSSender {
	templates := map[string]string{
		domain.NotificationTypeOrder:     "TENCENT_ORDER_TEMPLATE",
		domain.NotificationTypePayment:   "TENCENT_PAYMENT_TEMPLATE",
		domain.NotificationTypeRefund:    "TENCENT_REFUND_TEMPLATE",
		domain.NotificationTypeVoyage:    "TENCENT_VOYAGE_TEMPLATE",
		domain.NotificationTypeSystem:    "TENCENT_SYSTEM_TEMPLATE",
		domain.NotificationTypeInventory: "TENCENT_INVENTORY_TEMPLATE",
		domain.NotificationTypePromotion: "TENCENT_PROMOTION_TEMPLATE",
	}

	return &SMSSender{
		provider:     "tencent",
		accessKey:    accessKey,
		accessSecret: accessSecret,
		signName:     signName,
		templates:    templates,
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		endpoint:     "https://sms.tencentcloudapi.com",
	}
}

// IsAvailable checks if SMS sender is properly configured
func (s *SMSSender) IsAvailable() bool {
	return s.accessKey != "" && s.accessSecret != "" && s.signName != ""
}

// Send sends a notification via SMS
func (s *SMSSender) Send(ctx context.Context, notification *domain.Notification, user *domain.User) error {
	if user.Phone == "" {
		return fmt.Errorf("user has no phone number")
	}

	// Get template code
	templateCode := s.templates[notification.Type]
	if templateCode == "" {
		return fmt.Errorf("no SMS template for type: %s", notification.Type)
	}

	// Build template parameters
	params := s.buildTemplateParams(notification)

	// Send based on provider
	switch s.provider {
	case "aliyun":
		return s.sendAliyunSMS(ctx, user.Phone, templateCode, params)
	case "tencent":
		return s.sendTencentSMS(ctx, user.Phone, templateCode, params)
	default:
		return fmt.Errorf("unsupported SMS provider: %s", s.provider)
	}
}

// sendAliyunSMS sends SMS via Aliyun
func (s *SMSSender) sendAliyunSMS(ctx context.Context, phone, templateCode string, params SMSTemplateParam) error {
	// Build request body
	paramJSON, _ := json.Marshal(params)

	body := map[string]string{
		"PhoneNumbers":  phone,
		"SignName":      s.signName,
		"TemplateCode":  templateCode,
		"TemplateParam": string(paramJSON),
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", s.endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-acs-action", "SendSms")
	req.Header.Set("x-acs-version", "2017-05-25")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read SMS response: %w", err)
	}

	var result AliyunSMSResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return fmt.Errorf("failed to parse SMS response: %w", err)
	}

	if result.Code != "OK" {
		return fmt.Errorf("SMS API error: %s - %s", result.Code, result.Message)
	}

	return nil
}

// sendTencentSMS sends SMS via Tencent Cloud (simplified implementation)
func (s *SMSSender) sendTencentSMS(ctx context.Context, phone, templateCode string, params SMSTemplateParam) error {
	// Tencent Cloud SMS requires TC3-HMAC-SHA256 signature
	// This is a simplified implementation - full implementation would require
	// proper signature generation according to Tencent Cloud API specifications

	body := map[string]interface{}{
		"PhoneNumberSet": []string{phone},
		"SmsSdkAppId":    s.accessKey,
		"SignName":       s.signName,
		"TemplateId":     templateCode,
		"TemplateParamSet": []string{
			params.Title,
			params.Content,
		},
	}

	jsonBody, _ := json.Marshal(body)

	req, err := http.NewRequestWithContext(ctx, "POST", s.endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "sms.tencentcloudapi.com")
	req.Header.Set("X-TC-Action", "SendSms")
	req.Header.Set("X-TC-Version", "2021-01-11")
	req.Header.Set("X-TC-Timestamp", fmt.Sprintf("%d", time.Now().Unix()))

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Tencent SMS: %w", err)
	}
	defer resp.Body.Close()

	// Parse response
	respBody, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	// Check for errors in Response.Error
	if respObj, ok := result["Response"].(map[string]interface{}); ok {
		if errObj, ok := respObj["Error"].(map[string]interface{}); ok {
			return fmt.Errorf("Tencent SMS error: %v", errObj)
		}
	}

	return nil
}

// buildTemplateParams builds SMS template parameters from notification
func (s *SMSSender) buildTemplateParams(notification *domain.Notification) SMSTemplateParam {
	params := SMSTemplateParam{
		Title:   notification.Title,
		Content: notification.Content,
	}

	// Add type-specific parameters
	if notificationData, err := notification.GetData(); err == nil && notificationData != nil {
		if notificationData.OrderNo != "" {
			params.OrderNo = notificationData.OrderNo
		}
		if notificationData.Amount > 0 {
			params.Amount = fmt.Sprintf("%.2f", notificationData.Amount)
		}
		if notificationData.Count > 0 {
			params.Count = fmt.Sprintf("%d", notificationData.Count)
		}
		if notificationData.Days > 0 {
			params.Days = fmt.Sprintf("%d", notificationData.Days)
		}
	}

	return params
}

// Ensure SMSSender implements NotificationSender interface
var _ service.NotificationSender = (*SMSSender)(nil)

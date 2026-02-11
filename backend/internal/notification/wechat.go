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

// WechatTemplateSender sends notifications via WeChat template messages
type WechatTemplateSender struct {
	appID       string
	appSecret   string
	httpClient  *http.Client
	accessToken string
	tokenExpiry time.Time
}

// TemplateMessageRequest represents a WeChat template message request
type TemplateMessageRequest struct {
	Touser           string              `json:"touser"`
	TemplateID       string              `json:"template_id"`
	Page             string              `json:"page,omitempty"`
	Data             map[string]DataItem `json:"data"`
	MiniprogramState string              `json:"miniprogram_state,omitempty"`
	Lang             string              `json:"lang,omitempty"`
}

// DataItem represents a data item in template message
type DataItem struct {
	Value string `json:"value"`
	Color string `json:"color,omitempty"`
}

// WechatTemplateResponse represents WeChat API response
type WechatTemplateResponse struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// NewWechatTemplateSender creates a new WeChat template sender
func NewWechatTemplateSender(appID, appSecret string) *WechatTemplateSender {
	return &WechatTemplateSender{
		appID:      appID,
		appSecret:  appSecret,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// IsAvailable checks if WeChat sender is properly configured
func (s *WechatTemplateSender) IsAvailable() bool {
	return s.appID != "" && s.appSecret != ""
}

// Send sends a notification via WeChat template message
func (s *WechatTemplateSender) Send(ctx context.Context, notification *domain.Notification, user *domain.User) error {
	if user.WechatOpenID == "" {
		return fmt.Errorf("user has no WeChat OpenID")
	}

	// Get access token
	accessToken, err := s.getAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	// Get template ID based on notification type
	templateID := s.getTemplateID(notification.Type)
	if templateID == "" {
		return fmt.Errorf("no template ID for type: %s", notification.Type)
	}

	// Build template data
	data := s.buildTemplateData(notification)

	// Build request
	req := TemplateMessageRequest{
		Touser:           user.WechatOpenID,
		TemplateID:       templateID,
		Page:             notification.ActionURL,
		Data:             data,
		MiniprogramState: "formal",
		Lang:             "zh_CN",
	}

	// Send request
	url := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/subscribe/send?access_token=%s", accessToken)
	jsonBody, _ := json.Marshal(req)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	var result WechatTemplateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if result.ErrCode != 0 {
		return fmt.Errorf("WeChat API error: %d - %s", result.ErrCode, result.ErrMsg)
	}

	// Update sent time
	notification.SentAt = &[]time.Time{time.Now()}[0]

	return nil
}

// getAccessToken retrieves or refreshes WeChat access token
func (s *WechatTemplateSender) getAccessToken(ctx context.Context) (string, error) {
	// Check if current token is still valid (with 5 minute buffer)
	if s.accessToken != "" && time.Now().Before(s.tokenExpiry.Add(-5*time.Minute)) {
		return s.accessToken, nil
	}

	// Fetch new token
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		s.appID, s.appSecret,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	if tokenResp.ErrCode != 0 {
		return "", fmt.Errorf("WeChat API error: %d - %s", tokenResp.ErrCode, tokenResp.ErrMsg)
	}

	s.accessToken = tokenResp.AccessToken
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return s.accessToken, nil
}

// getTemplateID returns the template ID for a notification type
// In production, these would be configured via environment variables or database
func (s *WechatTemplateSender) getTemplateID(nType string) string {
	// Template IDs should be obtained from WeChat MP admin console
	templates := map[string]string{
		domain.NotificationTypeOrder:     "ORDER_CREATED_TEMPLATE_ID",
		domain.NotificationTypePayment:   "PAYMENT_SUCCESS_TEMPLATE_ID",
		domain.NotificationTypeRefund:    "REFUND_APPROVED_TEMPLATE_ID",
		domain.NotificationTypeVoyage:    "VOYAGE_REMINDER_TEMPLATE_ID",
		domain.NotificationTypeSystem:    "SYSTEM_NOTICE_TEMPLATE_ID",
		domain.NotificationTypeInventory: "INVENTORY_ALERT_TEMPLATE_ID",
		domain.NotificationTypePromotion: "PROMOTION_TEMPLATE_ID",
	}

	return templates[nType]
}

// buildTemplateData builds template message data from notification
func (s *WechatTemplateSender) buildTemplateData(notification *domain.Notification) map[string]DataItem {
	data := map[string]DataItem{
		"thing1": {Value: notification.Title},                    // Subject/title
		"thing2": {Value: notification.Content},                  // Content
		"time3":  {Value: time.Now().Format("2006-01-02 15:04")}, // Time
	}

	// Add specific data based on notification type
	if notificationData, err := notification.GetData(); err == nil && notificationData != nil {
		switch notification.Type {
		case domain.NotificationTypeOrder, domain.NotificationTypePayment:
			if notificationData.OrderNo != "" {
				data["character4"] = DataItem{Value: notificationData.OrderNo} // Order number
			}
			if notificationData.Amount > 0 {
				data["amount5"] = DataItem{Value: fmt.Sprintf("%.2f", notificationData.Amount)} // Amount
			}

		case domain.NotificationTypeRefund:
			if notificationData.RefundAmount > 0 {
				data["amount6"] = DataItem{Value: fmt.Sprintf("%.2f", notificationData.RefundAmount)}
			}

		case domain.NotificationTypeInventory:
			if notificationData.Count > 0 {
				data["number7"] = DataItem{Value: fmt.Sprintf("%d", notificationData.Count)}
			}
		}
	}

	return data
}

// Ensure WechatTemplateSender implements NotificationSender interface
var _ service.NotificationSender = (*WechatTemplateSender)(nil)

package payment

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"bytes"
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPaymentFailed       = errors.New("payment processing failed")
	ErrInvalidSignature    = errors.New("invalid payment signature")
	ErrPaymentNotFound     = errors.New("payment not found")
	ErrPaymentAlreadyPaid  = errors.New("payment already completed")
	ErrInsufficientPayment = errors.New("insufficient payment amount")
)

// PaymentProvider defines the interface for payment providers
type PaymentProvider interface {
	// CreatePayment creates a new payment request
	CreatePayment(ctx context.Context, order *domain.Order, description string) (*PaymentResult, error)

	// QueryPayment queries payment status
	QueryPayment(ctx context.Context, paymentNo string) (*PaymentQueryResult, error)

	// ProcessCallback processes payment callback/notification
	ProcessCallback(ctx context.Context, body []byte, signature string) (*CallbackResult, error)

	// VerifySignature verifies payment callback signature
	VerifySignature(body []byte, signature string) bool

	// Refund processes a refund
	Refund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*RefundResult, error)
}

// PaymentResult represents the result of creating a payment
type PaymentResult struct {
	PaymentNo    string `json:"payment_no"`
	ThirdPartyID string `json:"third_party_id,omitempty"`
	PayURL       string `json:"pay_url,omitempty"`
	PrepayID     string `json:"prepay_id,omitempty"`
	QRCodeURL    string `json:"qr_code_url,omitempty"`
	AppID        string `json:"app_id,omitempty"`
	Timestamp    string `json:"timestamp,omitempty"`
	NonceStr     string `json:"nonce_str,omitempty"`
	Package      string `json:"package,omitempty"`
	SignType     string `json:"sign_type,omitempty"`
	PaySign      string `json:"pay_sign,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
}

// PaymentQueryResult represents the result of querying payment
type PaymentQueryResult struct {
	Status       string  `json:"status"`
	Amount       float64 `json:"amount"`
	ThirdPartyID string  `json:"third_party_id,omitempty"`
	PaidAt       string  `json:"paid_at,omitempty"`
	ErrorMessage string  `json:"error_message,omitempty"`
}

// CallbackResult represents the result of processing a callback
type CallbackResult struct {
	PaymentNo    string  `json:"payment_no"`
	ThirdPartyID string  `json:"third_party_id"`
	Amount       float64 `json:"amount"`
	Status       string  `json:"status"`
	PaidAt       string  `json:"paid_at"`
}

// RefundResult represents the result of a refund
type RefundResult struct {
	RefundNo     string `json:"refund_no"`
	ThirdPartyID string `json:"third_party_id,omitempty"`
	Status       string `json:"status"`
}

// WechatPayConfig represents WeChat Pay V3 configuration
type WechatPayConfig struct {
	AppID       string
	MchID       string
	APIKey      string
	APIv3Key    string
	PrivateKey  string
	Certificate string
	SerialNo    string
	NotifyURL   string
	Sandbox     bool
}

// wechatPay implements PaymentProvider for WeChat Pay V3
type wechatPay struct {
	config      WechatPayConfig
	client      *http.Client
	paymentRepo repository.OrderRepository
}

// NewWechatPay creates a new WeChat Pay provider
func NewWechatPay(config WechatPayConfig, paymentRepo repository.OrderRepository) PaymentProvider {
	return &wechatPay{
		config:      config,
		client:      &http.Client{Timeout: 30 * time.Second},
		paymentRepo: paymentRepo,
	}
}

// CreatePayment creates a WeChat Pay native/App payment
func (w *wechatPay) CreatePayment(ctx context.Context, order *domain.Order, description string) (*PaymentResult, error) {
	paymentNo := generatePaymentNo()

	// Create payment record in database
	payment := &domain.Payment{
		OrderID:       order.ID.String(),
		PaymentNo:     paymentNo,
		PaymentMethod: domain.PaymentMethodWechat,
		Amount:        order.TotalAmount - order.DiscountAmount,
		Currency:      order.Currency,
		Status:        domain.PaymentStatusPending,
	}

	if err := w.paymentRepo.CreatePayment(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// Build payment request
	params := map[string]interface{}{
		"appid":        w.config.AppID,
		"mchid":        w.config.MchID,
		"description":  description,
		"out_trade_no": paymentNo,
		"notify_url":   w.config.NotifyURL,
		"amount": map[string]interface{}{
			"total":    int64(order.TotalAmount * 100), // Convert to cents
			"currency": order.Currency,
		},
		"time_expire": time.Now().Add(30 * time.Minute).Format(time.RFC3339),
	}

	// Make API request to WeChat Pay
	result, err := w.request(ctx, "POST", "/v3/pay/transactions/native", params)
	if err != nil {
		return nil, err
	}

	// Parse response
	codeURL, _ := result["code_url"].(string)
	prepayID, _ := result["prepay_id"].(string)

	return &PaymentResult{
		PaymentNo: paymentNo,
		PayURL:    codeURL,
		PrepayID:  prepayID,
		AppID:     w.config.AppID,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		NonceStr:  generateNonceStr(),
		Package:   "prepay_id=" + prepayID,
		SignType:  "RSA",
		ExpiresAt: time.Now().Add(30 * time.Minute).Format(time.RFC3339),
	}, nil
}

// QueryPayment queries WeChat Pay transaction status
func (w *wechatPay) QueryPayment(ctx context.Context, paymentNo string) (*PaymentQueryResult, error) {
	path := fmt.Sprintf("/v3/pay/transactions/out-trade-no/%s?mchid=%s", paymentNo, w.config.MchID)

	result, err := w.request(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	tradeState, _ := result["trade_state"].(string)
	amount, _ := result["amount"].(map[string]interface{})
	totalAmount := 0.0
	if amount != nil {
		total, _ := amount["total"].(float64)
		totalAmount = total / 100 // Convert from cents
	}

	status := w.mapTradeState(tradeState)

	return &PaymentQueryResult{
		Status:       status,
		Amount:       totalAmount,
		ThirdPartyID: result["transaction_id"].(string),
	}, nil
}

// ProcessCallback processes WeChat Pay notification
func (w *wechatPay) ProcessCallback(ctx context.Context, body []byte, signature string) (*CallbackResult, error) {
	// Verify signature
	if !w.VerifySignature(body, signature) {
		return nil, ErrInvalidSignature
	}

	// Parse callback body
	var notification struct {
		ID         string `json:"id"`
		CreateTime string `json:"create_time"`
		EventType  string `json:"event_type"`
		Resource   struct {
			Algorithm      string `json:"algorithm"`
			Ciphertext     string `json:"ciphertext"`
			AssociatedData string `json:"associated_data"`
			Nonce          string `json:"nonce"`
		} `json:"resource"`
	}

	if err := json.Unmarshal(body, &notification); err != nil {
		return nil, fmt.Errorf("failed to parse notification: %w", err)
	}

	// Decrypt resource
	plaintext, err := w.decrypt(notification.Resource.Ciphertext, notification.Resource.AssociatedData, notification.Resource.Nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt notification: %w", err)
	}

	var tradeData struct {
		AppID         string `json:"appid"`
		MchID         string `json:"mchid"`
		OutTradeNo    string `json:"out_trade_no"`
		TransactionID string `json:"transaction_id"`
		TradeType     string `json:"trade_type"`
		TradeState    string `json:"trade_state"`
		BankType      string `json:"bank_type"`
		Attach        string `json:"attach"`
		SuccessTime   string `json:"success_time"`
		Payer         struct {
			OpenID string `json:"openid"`
		} `json:"payer"`
		Amount struct {
			Total         int    `json:"total"`
			PayerTotal    int    `json:"payer_total"`
			Currency      string `json:"currency"`
			PayerCurrency string `json:"payer_currency"`
		} `json:"amount"`
	}

	if err := json.Unmarshal(plaintext, &tradeData); err != nil {
		return nil, fmt.Errorf("failed to parse trade data: %w", err)
	}

	// Build result
	status := w.mapTradeState(tradeData.TradeState)
	paidAt := tradeData.SuccessTime
	if paidAt == "" {
		paidAt = notification.CreateTime
	}

	return &CallbackResult{
		PaymentNo:    tradeData.OutTradeNo,
		ThirdPartyID: tradeData.TransactionID,
		Amount:       float64(tradeData.Amount.Total) / 100,
		Status:       status,
		PaidAt:       paidAt,
	}, nil
}

// VerifySignature verifies WeChat Pay callback signature
func (w *wechatPay) VerifySignature(body []byte, signature string) bool {
	// Parse certificate
	block, _ := pem.Decode([]byte(w.config.Certificate))
	if block == nil {
		return false
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return false
	}

	// Decode signature
	sig, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	// Verify
	hash := sha256.Sum256(body)
	err = rsa.VerifyPKCS1v15(cert.PublicKey.(*rsa.PublicKey), crypto.SHA256, hash[:], sig)
	return err == nil
}

// Refund processes a WeChat Pay refund
func (w *wechatPay) Refund(ctx context.Context, payment *domain.Payment, amount float64, reason string) (*RefundResult, error) {
	refundNo := generateRefundNo()

	params := map[string]interface{}{
		"out_trade_no":  payment.PaymentNo,
		"out_refund_no": refundNo,
		"reason":        reason,
		"amount": map[string]interface{}{
			"refund":   int64(amount * 100),
			"total":    int64(payment.Amount * 100),
			"currency": payment.Currency,
		},
		"notify_url": w.config.NotifyURL + "/refund",
	}

	result, err := w.request(ctx, "POST", "/v3/refund/domestic/refunds", params)
	if err != nil {
		return nil, err
	}

	status, _ := result["status"].(string)
	refundID, _ := result["refund_id"].(string)

	return &RefundResult{
		RefundNo:     refundNo,
		ThirdPartyID: refundID,
		Status:       status,
	}, nil
}

// Helper methods

func (w *wechatPay) request(ctx context.Context, method, path string, body interface{}) (map[string]interface{}, error) {
	var bodyReader io.Reader
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(jsonBody)
	}

	baseURL := "https://api.mch.weixin.qq.com"
	if w.config.Sandbox {
		baseURL = "https://api.mch.weixin.qq.com/sandboxnew"
	}

	req, err := http.NewRequestWithContext(ctx, method, baseURL+path, bodyReader)
	if err != nil {
		return nil, err
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", w.buildAuthorization(method, path, body))

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("WeChat Pay API error: %s", string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

func (w *wechatPay) buildAuthorization(method, path string, body interface{}) string {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := generateNonceStr()

	var bodyString string
	if body != nil {
		jsonBody, _ := json.Marshal(body)
		bodyString = string(jsonBody)
	}

	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", method, path, timestamp, nonce) + bodyString + "\n"

	// Sign with private key
	signature := w.sign(message)

	return fmt.Sprintf("WECHATPAY2-SHA256-RSA2048 mchid=\"%s\",nonce_str=\"%s\",signature=\"%s\",timestamp=\"%s\",serial_no=\"%s\"",
		w.config.MchID, nonce, signature, timestamp, w.config.SerialNo)
}

func (w *wechatPay) sign(message string) string {
	// Parse private key
	block, _ := pem.Decode([]byte(w.config.PrivateKey))
	if block == nil {
		return ""
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}

	// Sign
	hash := sha256.Sum256([]byte(message))
	signature, err := rsa.SignPKCS1v15(nil, privateKey.(*rsa.PrivateKey), crypto.SHA256, hash[:])
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(signature)
}

func (w *wechatPay) decrypt(ciphertext, associatedData, nonce string) ([]byte, error) {
	// Decode ciphertext from base64
	cipherBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// SEC-003: Implement AES-256-GCM decryption with APIv3 key
	key := []byte(w.config.APIv3Key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceBytes := []byte(nonce)
	additionalData := []byte(associatedData)

	plaintext, err := aesGCM.Open(nil, nonceBytes, cipherBytes, additionalData)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	return plaintext, nil
}

func (w *wechatPay) mapTradeState(state string) string {
	switch state {
	case "SUCCESS":
		return domain.PaymentStatusSuccess
	case "REFUND":
		return domain.PaymentStatusRefunded
	case "NOTPAY":
		return domain.PaymentStatusPending
	case "CLOSED":
		return domain.PaymentStatusCancelled
	case "USERPAYING":
		return domain.PaymentStatusProcessing
	case "PAYERROR":
		return domain.PaymentStatusFailed
	default:
		return domain.PaymentStatusPending
	}
}

func generatePaymentNo() string {
	return fmt.Sprintf("PAY%s%s", time.Now().Format("20060102"), uuid.New().String()[:12])
}

func generateRefundNo() string {
	return fmt.Sprintf("REF%s%s", time.Now().Format("20060102"), uuid.New().String()[:12])
}

func generateNonceStr() string {
	return uuid.New().String()[:32]
}

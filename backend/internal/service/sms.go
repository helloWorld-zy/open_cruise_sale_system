package service

import (
	"backend/internal/repository"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidPhone    = errors.New("invalid phone number")
	ErrInvalidCode     = errors.New("invalid verification code")
	ErrCodeExpired     = errors.New("verification code expired")
	ErrTooManyRequests = errors.New("too many verification requests")
	ErrCodeNotFound    = errors.New("verification code not found")
)

// SMSConfig holds SMS service configuration
type SMSConfig struct {
	Provider     string // aliyun, tencent, twilio
	AccessKey    string
	AccessSecret string
	SignName     string
	TemplateCode string
}

// SMSService handles SMS verification
type SMSService struct {
	config SMSConfig
	repo   repository.UserRepository
	// In-memory code storage (use Redis in production)
	codeStore map[string]*VerificationCode
}

// VerificationCode stores verification code info
type VerificationCode struct {
	Phone     string
	Code      string
	ExpiresAt time.Time
	Attempts  int
}

// NewSMSService creates a new SMS service
func NewSMSService(config SMSConfig, repo repository.UserRepository) *SMSService {
	return &SMSService{
		config:    config,
		repo:      repo,
		codeStore: make(map[string]*VerificationCode),
	}
}

// SendVerificationCode sends verification code to phone
func (s *SMSService) SendVerificationCode(ctx context.Context, phone string) error {
	// Validate phone format (simplified Chinese phone validation)
	if len(phone) != 11 || phone[0] != '1' {
		return ErrInvalidPhone
	}

	// Check rate limit
	if existing, ok := s.codeStore[phone]; ok {
		if time.Since(existing.ExpiresAt.Add(-10*time.Minute)) < 60*time.Second {
			return ErrTooManyRequests
		}
	}

	// Generate 6-digit code
	code := s.generateCode(6)

	// Store code
	s.codeStore[phone] = &VerificationCode{
		Phone:     phone,
		Code:      code,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		Attempts:  0,
	}

	// Send SMS (mock implementation)
	// In production, integrate with Aliyun SMS, Tencent SMS, etc.
	if err := s.sendSMS(phone, code); err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}

// VerifyCode verifies the SMS code
func (s *SMSService) VerifyCode(ctx context.Context, phone, code string) error {
	stored, ok := s.codeStore[phone]
	if !ok {
		return ErrCodeNotFound
	}

	// Check expiration
	if time.Now().After(stored.ExpiresAt) {
		delete(s.codeStore, phone)
		return ErrCodeExpired
	}

	// Check attempts
	if stored.Attempts >= 5 {
		delete(s.codeStore, phone)
		return ErrTooManyRequests
	}

	// Verify code
	stored.Attempts++
	if stored.Code != code {
		return ErrInvalidCode
	}

	// Success - remove code
	delete(s.codeStore, phone)
	return nil
}

// VerifyAndLogin verifies code and returns user
func (s *SMSService) VerifyAndLogin(ctx context.Context, phone, code string) (*UserLoginResult, error) {
	if err := s.VerifyCode(ctx, phone, code); err != nil {
		return nil, err
	}

	// Get or create user
	user, err := s.repo.GetByPhone(ctx, phone)
	if err != nil {
		// Create new user
		user = &domain.User{
			Phone:         phone,
			PhoneVerified: true,
			Status:        domain.UserStatusActive,
		}
		if err := s.repo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else {
		// Mark phone as verified if not already
		if !user.PhoneVerified {
			user.PhoneVerified = true
			s.repo.Update(ctx, user)
		}
	}

	return &UserLoginResult{
		User:      user,
		IsNewUser: false, // Set based on user creation time
	}, nil
}

// generateCode generates a random numeric code
func (s *SMSService) generateCode(length int) string {
	max := 1
	for i := 0; i < length; i++ {
		max *= 10
	}

	bytes := make([]byte, 4)
	rand.Read(bytes)
	num := int(bytes[0])<<24 | int(bytes[1])<<16 | int(bytes[2])<<8 | int(bytes[3])

	code := num % max
	format := fmt.Sprintf("%%0%dd", length)
	return fmt.Sprintf(format, code)
}

// sendSMS sends SMS using configured provider
func (s *SMSService) sendSMS(phone, code string) error {
	// Mock implementation - log to console
	fmt.Printf("[SMS] Sending code %s to %s\n", code, phone)

	// In production, implement actual SMS sending:
	// - Aliyun SMS
	// - Tencent Cloud SMS
	// - Twilio

	return nil
}

// UserLoginResult represents login result
type UserLoginResult struct {
	User      *domain.User
	IsNewUser bool
}

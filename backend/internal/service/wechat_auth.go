package service

import (
	"backend/internal/domain"
	"backend/internal/repository"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var (
	ErrWechatAuthFailed  = errors.New("wechat authentication failed")
	ErrWechatCodeInvalid = errors.New("invalid wechat code")
	ErrUserNotFound      = errors.New("user not found")
	ErrPhoneNotBound     = errors.New("phone number not bound")
)

// WechatAuthConfig holds WeChat mini program configuration
type WechatAuthConfig struct {
	AppID     string
	AppSecret string
}

// WechatAuthService handles WeChat authentication
type WechatAuthService struct {
	config WechatAuthConfig
	repo   repository.UserRepository
	jwtSvc *JWTService
}

// NewWechatAuthService creates a new WeChat auth service
func NewWechatAuthService(config WechatAuthConfig, repo repository.UserRepository, jwtSvc *JWTService) *WechatAuthService {
	return &WechatAuthService{
		config: config,
		repo:   repo,
		jwtSvc: jwtSvc,
	}
}

// WechatLoginResult represents the result of WeChat login
type WechatLoginResult struct {
	User         *domain.User `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	IsNewUser    bool         `json:"is_new_user"`
}

// WechatLogin handles WeChat mini program login
func (s *WechatAuthService) WechatLogin(ctx context.Context, code string) (*WechatLoginResult, error) {
	// Exchange code for session info
	session, err := s.code2Session(code)
	if err != nil {
		return nil, ErrWechatAuthFailed
	}

	// Check if user exists
	user, err := s.repo.GetByWechatOpenID(ctx, session.OpenID)
	if err != nil {
		// Create new user
		user = &domain.User{
			WechatOpenID:  session.OpenID,
			WechatUnionID: session.UnionID,
			Status:        domain.UserStatusActive,
		}
		if err := s.repo.Create(ctx, user); err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}

		tokens, err := s.jwtSvc.GenerateTokenPair(user.ID)
		if err != nil {
			return nil, err
		}

		return &WechatLoginResult{
			User:         user,
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			IsNewUser:    true,
		}, nil
	}

	// Existing user
	if !user.CanLogin() {
		return nil, errors.New("user account is disabled")
	}

	// Update last login
	user.UpdateLastLogin("") // IP will be set by handler
	s.repo.Update(ctx, user)

	tokens, err := s.jwtSvc.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, err
	}

	return &WechatLoginResult{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IsNewUser:    false,
	}, nil
}

// WechatPhoneLogin handles login with encrypted phone data
func (s *WechatAuthService) WechatPhoneLogin(ctx context.Context, code, encryptedData, iv string) (*WechatLoginResult, error) {
	// First get session
	session, err := s.code2Session(code)
	if err != nil {
		return nil, ErrWechatAuthFailed
	}

	// Decrypt phone data
	phoneData, err := s.decryptWechatData(session.SessionKey, encryptedData, iv)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt phone data: %w", err)
	}

	phoneNumber := phoneData.PhoneNumber

	// Check if user exists by phone
	user, err := s.repo.GetByPhone(ctx, phoneNumber)
	if err != nil {
		// Check by WeChat OpenID
		user, err = s.repo.GetByWechatOpenID(ctx, session.OpenID)
		if err != nil {
			// Create new user with phone
			user = &domain.User{
				Phone:         phoneNumber,
				PhoneVerified: true,
				WechatOpenID:  session.OpenID,
				WechatUnionID: session.UnionID,
				Status:        domain.UserStatusActive,
			}
			if err := s.repo.Create(ctx, user); err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
		} else {
			// Bind phone to existing WeChat user
			user.Phone = phoneNumber
			user.PhoneVerified = true
			s.repo.Update(ctx, user)
		}
	} else {
		// Bind WeChat to existing phone user
		user.WechatOpenID = session.OpenID
		user.WechatUnionID = session.UnionID
		s.repo.Update(ctx, user)
	}

	tokens, err := s.jwtSvc.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, err
	}

	return &WechatLoginResult{
		User:         user,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		IsNewUser:    false,
	}, nil
}

// WechatSessionInfo holds session info from WeChat
type WechatSessionInfo struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid,omitempty"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

// code2Session exchanges code for session info
func (s *WechatAuthService) code2Session(code string) (*WechatSessionInfo, error) {
	url := fmt.Sprintf(
		"https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		s.config.AppID, s.config.AppSecret, code,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var session WechatSessionInfo
	if err := json.Unmarshal(body, &session); err != nil {
		return nil, err
	}

	if session.ErrCode != 0 {
		return nil, fmt.Errorf("wechat error %d: %s", session.ErrCode, session.ErrMsg)
	}

	return &session, nil
}

// WechatPhoneData holds decrypted phone data
type WechatPhoneData struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
	Watermark       struct {
		AppID     string `json:"appid"`
		Timestamp int64  `json:"timestamp"`
	} `json:"watermark"`
}

// decryptWechatData decrypts WeChat encrypted data
func (s *WechatAuthService) decryptWechatData(sessionKey, encryptedData, iv string) (*WechatPhoneData, error) {
	// Decode base64
	sessionKeyBytes, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, err
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	ivBytes, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, err
	}

	// AES-128-CBC decryption
	block, err := aes.NewCipher(sessionKeyBytes)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, ivBytes)
	mode.CryptBlocks(encryptedBytes, encryptedBytes)

	// Remove PKCS7 padding
	data := s.pkcs7Unpad(encryptedBytes)

	var phoneData WechatPhoneData
	if err := json.Unmarshal(data, &phoneData); err != nil {
		return nil, err
	}

	// Verify watermark
	if phoneData.Watermark.AppID != s.config.AppID {
		return nil, errors.New("invalid watermark appid")
	}

	return &phoneData, nil
}

// pkcs7Unpad removes PKCS7 padding
func (s *WechatAuthService) pkcs7Unpad(data []byte) []byte {
	length := len(data)
	if length == 0 {
		return data
	}
	unpadding := int(data[length-1])
	if unpadding > length {
		return data
	}
	return data[:(length - unpadding)]
}

// UpdateWechatUserInfo updates user info from WeChat
func (s *WechatAuthService) UpdateWechatUserInfo(ctx context.Context, userID string, nickname, avatarURL string) error {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if nickname != "" {
		user.WechatNickname = nickname
		if user.Nickname == "" {
			user.Nickname = nickname
		}
	}

	if avatarURL != "" {
		user.WechatAvatarURL = avatarURL
		if user.AvatarURL == "" {
			user.AvatarURL = avatarURL
		}
	}

	return s.repo.Update(ctx, user)
}

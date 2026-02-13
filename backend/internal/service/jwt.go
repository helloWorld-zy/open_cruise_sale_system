package service

import (
	"backend/internal/config"
	"backend/internal/middleware"
)

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type JWTService struct {
	config *config.JWTConfig
}

func NewJWTService(config *config.JWTConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}

func (s *JWTService) GenerateTokenPair(userID string) (*TokenPair, error) {
	accessToken, err := middleware.GenerateToken(userID, "", "user", s.config)
	if err != nil {
		return nil, err
	}

	refreshToken, err := middleware.GenerateRefreshToken(userID, s.config)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

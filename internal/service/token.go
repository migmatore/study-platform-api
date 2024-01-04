package service

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/migmatore/study-platform-api/config"
	"github.com/migmatore/study-platform-api/internal/core"
	"time"
)

type TokenService struct {
	config *config.Config
}

func NewTokenService(config *config.Config) *TokenService {
	return &TokenService{config: config}
}

func (s TokenService) GetToken(userId int, role string) (core.TokenWithClaims, error) {
	expires := time.Now().Add(time.Minute * time.Duration(s.config.Server.JwtExpTimeMin)).Unix()

	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.config.Server.JwtSecretKey))
	if err != nil {
		return core.TokenWithClaims{}, err
	}

	return core.TokenWithClaims{
		Token:  t,
		UserId: userId,
		Role:   role,
	}, nil
}

func (s TokenService) GetRefreshToken(userId int) (core.RefreshTokenWithClaims, error) {
	expires := time.Now().Add(time.Minute * time.Duration(s.config.Server.JwtExpTimeMin)).Unix()

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.config.Server.JwtSecretKey))
	if err != nil {
		return core.RefreshTokenWithClaims{}, err
	}

	return core.RefreshTokenWithClaims{
		Token:  t,
		UserId: userId,
	}, nil
}

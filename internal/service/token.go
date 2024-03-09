package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/migmatore/study-platform-api/config"
	"github.com/migmatore/study-platform-api/internal/apperrors"
	"github.com/migmatore/study-platform-api/internal/core"
	"time"
)

type TokenService struct {
	config *config.Config
}

func NewTokenService(config *config.Config) *TokenService {
	return &TokenService{config: config}
}

func (s TokenService) Token(userId int, role string) (core.TokenWithClaims, error) {
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

func (s TokenService) WSToken(userId int, role string) (core.TokenWithClaims, error) {
	expires := time.Now().Add(time.Hour * time.Duration(s.config.Server.WSJwtExpTimeHour)).Unix()

	claims := jwt.MapClaims{
		"user_id": userId,
		"role":    role,
		"exp":     expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.config.Server.WSJwtSecretKey))
	if err != nil {
		return core.TokenWithClaims{}, err
	}

	return core.TokenWithClaims{
		Token:  t,
		UserId: userId,
		Role:   role,
	}, nil
}

func (s TokenService) RefreshToken(userId int) (core.RefreshTokenWithClaims, error) {
	expires := time.Now().Add(time.Hour * time.Duration(s.config.Server.JwtRefreshExpTimeHour)).Unix()

	claims := jwt.MapClaims{
		"user_id": userId,
		"exp":     expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(s.config.Server.JwtRefreshSecretKey))
	if err != nil {
		return core.RefreshTokenWithClaims{}, err
	}

	return core.RefreshTokenWithClaims{
		Token:  t,
		UserId: userId,
	}, nil
}

func (s TokenService) ExtractTokenMetadata(tokenString string) (core.TokenMetadata, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Server.JwtSecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return core.TokenMetadata{}, apperrors.InvalidToken
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return core.TokenMetadata{}, apperrors.ExpiredToken
		}

		return core.TokenMetadata{}, err
	}

	if !token.Valid {
		return core.TokenMetadata{}, apperrors.InvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Expires time.
		expires := int64(claims["exp"].(float64))
		userId := int(claims["user_id"].(float64))
		role := claims["role"].(string)

		return core.TokenMetadata{
			Expires: expires,
			UserId:  userId,
			Role:    role,
		}, nil
	}

	return core.TokenMetadata{}, err
}

func (s TokenService) ExtractWSTokenMetadata(tokenString string) (core.TokenMetadata, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Server.WSJwtSecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return core.TokenMetadata{}, apperrors.InvalidToken
		}

		if errors.Is(err, jwt.ErrTokenExpired) {
			return core.TokenMetadata{}, apperrors.ExpiredToken
		}

		return core.TokenMetadata{}, err
	}

	if !token.Valid {
		return core.TokenMetadata{}, apperrors.InvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		expires := int64(claims["exp"].(float64))
		userId := int(claims["user_id"].(float64))
		role := claims["role"].(string)

		return core.TokenMetadata{
			Expires: expires,
			UserId:  userId,
			Role:    role,
		}, nil
	}

	return core.TokenMetadata{}, err
}

func (s TokenService) ExtractRefreshTokenMetadata(tokenString string) (core.RefreshTokenMetadata, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Server.JwtRefreshSecretKey), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) || errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return core.RefreshTokenMetadata{}, apperrors.InvalidToken
		}

		return core.RefreshTokenMetadata{}, err
	}

	if !token.Valid {
		return core.RefreshTokenMetadata{}, apperrors.InvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		// Expires time.
		expires := int64(claims["exp"].(float64))
		userId := int(claims["user_id"].(float64))

		return core.RefreshTokenMetadata{
			Expires: expires,
			UserId:  userId,
		}, nil
	}

	return core.RefreshTokenMetadata{}, err
}

package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/magmaheat/merchStore/internal/repo"
	"time"
)

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type Auth interface {
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(accessToken string) (int, error)
}

type AuthService struct {
	repo     repo.Repo
	signKey  string
	tokenTTL time.Duration
}

func NewAuthService(repo repo.Repo, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		repo:     repo,
		signKey:  signKey,
		tokenTTL: tokenTTL,
	}
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	return "", nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, ErrCannotParseToken
	}

	return claims.UserId, nil
}

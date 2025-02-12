package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/magmaheat/merchStore/internal/repo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
	const fn = "service.AuthService.GenerateToken"

	userID, hash, err := s.repo.GetUser(ctx, input.Username)
	if err != nil {
		return "", err
	}

	if userID == 0 {
		hash, _ = hashPassword(input.Password)
		userID, err = s.repo.CreateUserWithBalance(ctx, input.Username, hash)
		if err != nil {
			return "", err
		}
	} else {
		if !checkPassword(input.Password, hash) {
			log.Errorf("%s.checkPassword: %s", fn, ErrInvalidPassword.Error())
			return "", ErrInvalidPassword
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: userID,
	})

	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		log.Errorf("%s.SignedString: %s: %v", fn, ErrCannotSignToken.Error(), err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
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

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

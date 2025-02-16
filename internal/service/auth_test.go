package service

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/magmaheat/merchStore/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	TokenTTL = 15 * time.Minute
)

func TestAuthService_GenerateToken(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		password      string
		mockSetup     func(mockRepo *mocks.Repo, mockHasher *mocks.PasswordHasher)
		expectedError error
	}{
		{
			name:     "New user",
			username: "newUser",
			password: "password",
			mockSetup: func(mockRepo *mocks.Repo, mockHasher *mocks.PasswordHasher) {
				mockRepo.On("GetUserIdWithPassword", mock.Anything, "newUser").Return(0, "", nil)
				mockRepo.On("CreateUserWithBalance", mock.Anything, "newUser", mock.Anything).Return(1, nil)
				mockHasher.On("Hash", "password").Return("hashedPassword", nil)
			},
			expectedError: nil,
		},
		{
			name:     "Invalid password",
			username: "user",
			password: "wrongPassword",
			mockSetup: func(mockRepo *mocks.Repo, mockHasher *mocks.PasswordHasher) {
				mockRepo.On("GetUserIdWithPassword", mock.Anything, "user").Return(1, "123", nil)
				mockHasher.On("Check", mock.Anything, mock.Anything).Return(false)
			},
			expectedError: ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.Repo{}
			mockHasher := &mocks.PasswordHasher{}

			tt.mockSetup(mockRepo, mockHasher)

			authService := NewAuthService(mockRepo, mockHasher, "secret", TokenTTL)

			token, err := authService.GenerateToken(context.Background(), tt.username, tt.password)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token, "Expected a token but got an empty string")
			}

			mockRepo.AssertExpectations(t)
			mockHasher.AssertExpectations(t)
		})
	}
}

func TestAuthService_ParseToken(t *testing.T) {
	secretKey := "test-secret"
	authService := &AuthService{signKey: secretKey}

	tests := []struct {
		name          string
		tokenSetup    func() string
		expectedID    int
		expectedError error
	}{
		{
			name: "Valid token",
			tokenSetup: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
						IssuedAt:  time.Now().Unix(),
					},
					UserId: 42,
				})
				tokenString, err := token.SignedString([]byte(secretKey))
				require.NoError(t, err)
				return tokenString
			},
			expectedID:    42,
			expectedError: nil,
		},
		{
			name: "Invalid signature",
			tokenSetup: func() string {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
					StandardClaims: jwt.StandardClaims{
						ExpiresAt: time.Now().Add(time.Hour).Unix(),
						IssuedAt:  time.Now().Unix(),
					},
					UserId: 42,
				})
				tokenString, err := token.SignedString([]byte("wrong-secret"))
				require.NoError(t, err)
				return tokenString
			},
			expectedID:    0,
			expectedError: ErrCannotParseToken,
		},
		{
			name: "Malformed token",
			tokenSetup: func() string {
				return "invalid-token"
			},
			expectedID:    0,
			expectedError: ErrCannotParseToken,
		},
		//{
		//	name: "Wrong signing method",
		//	tokenSetup: func() string {
		//		token := jwt.NewWithClaims(jwt.SigningMethodRS256, &TokenClaims{
		//			StandardClaims: jwt.StandardClaims{
		//				ExpiresAt: time.Now().Add(time.Hour).Unix(),
		//				IssuedAt:  time.Now().Unix(),
		//			},
		//			UserId: 42,
		//		})
		//		tokenString, err := token.SignedString([]byte(secretKey))
		//		require.NoError(t, err)
		//		return tokenString
		//	},
		//	expectedID:    0,
		//	expectedError: ErrCannotParseToken,
		//},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString := tt.tokenSetup()

			userID, err := authService.ParseToken(tokenString)

			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
				assert.Zero(t, userID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, userID)
			}
		})
	}
}

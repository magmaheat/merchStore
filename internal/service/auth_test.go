package service

import (
	"context"
	"testing"
	"time"

	"github.com/magmaheat/merchStore/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_GenerateToken_NewUser(t *testing.T) {
	mockRepo := &mocks.Repo{}

	mockRepo.On("GetUserIdWithPassword", mock.Anything, "newUser").Return(0, "", nil)
	mockRepo.On("CreateUserWithBalance", mock.Anything, "newUser", mock.Anything).Return(1, nil)

	authService := NewAuthService(mockRepo, "secret", 15*time.Minute)

	token, err := authService.GenerateToken(context.Background(), "newUser", "password")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	mockRepo.AssertExpectations(t)
}

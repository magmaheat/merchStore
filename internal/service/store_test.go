package service

import (
	"context"
	"errors"
	"github.com/magmaheat/merchStore/internal/mocks"
	"github.com/magmaheat/merchStore/internal/models"
	"github.com/magmaheat/merchStore/internal/repo"
	"github.com/magmaheat/merchStore/internal/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestStoreService_BuyItem(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mockRepo *mocks.Repo)
		nameItem      string
		expectedError error
	}{
		{
			name: "Success",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetPriceItem", mock.Anything, "item1").Return(10, nil)
				mockRepo.On("UpdateBalance", mock.Anything, mock.Anything, -10).Return(nil)
				mockRepo.On("AddItem", mock.Anything, mock.Anything, "item1").Return(nil)
			},
			nameItem:      "item1",
			expectedError: nil,
		},
		{
			name: "Item not found",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetPriceItem", mock.Anything, "item2").Return(0, nil)
			},
			nameItem:      "item2",
			expectedError: ErrItemNotFound,
		},
		{
			name: "Error updating balance",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetPriceItem", mock.Anything, "item1").Return(10, nil)
				mockRepo.On("UpdateBalance", mock.Anything, mock.Anything, -10).Return(errors.New("db error"))
			},
			nameItem:      "item1",
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.Repo{}

			tt.setupMock(mockRepo)

			storeService := NewStoreService(mockRepo)

			ctx := context.WithValue(context.Background(), types.UserIdCtx, 1)
			err := storeService.BuyItem(ctx, tt.nameItem)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStoreService_SendCoin(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mockRepo *mocks.Repo)
		toUser        string
		amount        int
		expectedError error
	}{
		{
			name: "Success",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetUserIdWithPassword", mock.Anything, "user2").Return(2, "", nil)
				mockRepo.On("TransferCoins", mock.Anything, 1, 2, 10).Return(nil)
			},
			toUser:        "user2",
			amount:        10,
			expectedError: nil,
		},
		{
			name: "Error getting user id",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetUserIdWithPassword", mock.Anything, "user2").Return(0, "", errors.New("user not found"))
			},
			toUser:        "user2",
			amount:        10,
			expectedError: errors.New("user not found"),
		},
		{
			name: "Error transferring coins",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetUserIdWithPassword", mock.Anything, "user2").Return(2, "", nil)
				mockRepo.On("TransferCoins", mock.Anything, 1, 2, 10).Return(errors.New("db error"))
			},
			toUser:        "user2",
			amount:        10,
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.Repo{}

			tt.setupMock(mockRepo)

			storeService := &StoreService{repo: mockRepo}

			ctx := context.WithValue(context.Background(), types.UserIdCtx, 1)
			err := storeService.SendCoin(ctx, tt.toUser, tt.amount)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestStoreService_GetInfo(t *testing.T) {
	tests := []struct {
		name          string
		setupMock     func(mockRepo *mocks.Repo)
		expectedInfo  *models.Info
		expectedError error
	}{
		{
			name: "Success",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetBalance", mock.Anything, 1).Return(100, nil)
				mockRepo.On("GetItems", mock.Anything, 1).Return([]repo.Item{{Type: "item1", Quantity: 1}}, nil)
				mockRepo.On("GetReceivedTransactions", mock.Anything, 1).Return([]repo.GetReceivedTransactionOutput{}, nil)
				mockRepo.On("GetSentTransactions", mock.Anything, 1).Return([]repo.GetSentTransactionOutput{}, nil)
			},
			expectedInfo: &models.Info{
				Coins: 100,
				Inventory: []models.Item{
					{Type: "item1", Quantity: 1},
				},
				CoinHistory: models.CoinHistory{
					Received: []models.ReceivedTransaction{},
					Sent:     []models.SentTransaction{},
				},
			},
			expectedError: nil,
		},
		{
			name: "Error getting balance",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetBalance", mock.Anything, 1).Return(0, errors.New("db error"))
			},
			expectedInfo:  nil,
			expectedError: errors.New("db error"),
		},
		{
			name: "Error getting items",
			setupMock: func(mockRepo *mocks.Repo) {
				mockRepo.On("GetBalance", mock.Anything, 1).Return(100, nil)
				mockRepo.On("GetItems", mock.Anything, 1).Return(nil, errors.New("db error"))
			},
			expectedInfo:  nil,
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mocks.Repo{}

			tt.setupMock(mockRepo)

			storeService := NewStoreService(mockRepo)

			ctx := context.WithValue(context.Background(), types.UserIdCtx, 1)
			info, err := storeService.GetInfo(ctx)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedInfo, info)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

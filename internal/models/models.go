package models

import "github.com/magmaheat/merchStore/internal/repo"

type Info struct {
	Coins       int         `json:"coins"`
	Inventory   []Item      `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

type Item struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

func NewInfo(coins int, inventory []Item, received []ReceivedTransaction, sent []SentTransaction) *Info {
	return &Info{
		Coins:       coins,
		Inventory:   inventory,
		CoinHistory: CoinHistory{Received: received, Sent: sent},
	}
}

func ConvertReceivedTransactions(transactions []repo.GetReceivedTransactionOutput) []ReceivedTransaction {
	result := make([]ReceivedTransaction, len(transactions))
	for i, t := range transactions {
		result[i] = ReceivedTransaction{
			FromUser: t.FromUser,
			Amount:   t.Amount,
		}
	}
	return result
}

func ConvertSentTransactions(transactions []repo.GetSentTransactionOutput) []SentTransaction {
	result := make([]SentTransaction, len(transactions))
	for i, t := range transactions {
		result[i] = SentTransaction{
			ToUser: t.ToUser,
			Amount: t.Amount,
		}
	}
	return result
}

func ConvertInventory(items []repo.Item) []Item {
	result := make([]Item, len(items))
	for i, item := range items {
		result[i] = Item{
			Type:     item.Type,
			Quantity: item.Quantity,
		}
	}
	return result
}

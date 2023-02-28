package fintoc

import (
	"time"
)

const (
	// checking_account
	CheckingAccount = "checking_account"
	// savings_account
	SavingsACcount = "savings_account"
	// sight_account
	SightAccount = "sight_account"
	// rut_account	Cuenta RUT. Only available for Chile
	RutAccount = "rut_account"
	// line_of_credit.
	LineOfCredit = "line_of_credit"
	// credit_card
	CreditCard = "credit_card"
)

// FinancialInstitution represents a financia institution inside fintoc's API
type FinancialInstitution struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
}

// TransferAccount represents a transfer account inside fintoc's API
type TransferAccount struct {
	HolderID    string                `json:"holder_id"` // chilean DNI
	HolderName  string                `json:"holder_name"`
	Number      *string               `json:"number,omitempty"`
	Institution *FinancialInstitution `json:"institution"`
}

type Movement struct {
	ID               string           `json:"id"`
	Object           string           `json:"object"`
	Amount           int64            `json:"amount"`
	PostDate         time.Time        `json:"post_date"`
	Description      string           `json:"description"`
	TransactionDate  time.Time        `json:"transaction_date,omitempty"`
	Currency         string           `json:"currency"`
	ReferenceID      string           `json:"reference_id,omitempty"`
	Type             string           `json:"type"`
	Pending          bool             `json:"pending"`
	RecipientAccount *TransferAccount `json:"recipient_account,omitempty"`
	SenderAccount    *TransferAccount `json:"sender_account,omitempty"`
	Comment          string           `json:"comment,omitempty"`
}

type Pages struct {
	Currenct string
	Next     string
	Last     string
}

type Account struct {
	// Bank account unique identifier
	ID string `json:"id,omitempty"`
	// Object identifier, in this case is: account
	Object string `json:"object,omitempty"`
	// Standard account name
	Name string `json:"name,omitempty"`
	// Institution bank account name
	OfficialName string `json:"official_name,omitempty"`
	// Account number, without hiphens and zeros prefix
	Number string `json:"number,omitempty"`
	// Owner identifier. For Chile is the owner RUT
	HolderID string `json:"holder_id,omitempty"`
	// Owner name
	HolderName string `json:"holder_name,omitempty"`
	// Account type
	AccountType string `json:"type,omitempty"`
	// Currency ISO code (3 letters)
	Currency string `json:"currency,omitempty"`
	// Account balance
	Balance *AccountBalance `json:"balance,omitempty"`
	// Last refreshed time
	RefreshedAt string `json:"refreshed_at,omitempty"`
}

type AccountBalance struct {
	Available int64 `json:"available,omitempty"`
	Current   int64 `json:"current,omitempty"`
	Limit     int64 `json:"limit,omitempty"`
}

type RefreshIntent struct {
	ID                string `json:"id,omitempty"`
	RefreshedObject   string `json:"refreshed_object,omitempty"`
	RefreshedObjectID string `json:"refreshed_object_id,omitempty"`
	Object            string `json:"object,omitempty"`
	Status            string `json:"status,omitempty"`
	NewMovements      int    `json:"new_movements,omitempty"`
	CreatedAt         string `json:"created_at,omitempty"`
	Type              string `json:"type,omitempty"`
}

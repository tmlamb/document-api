package model

import "github.com/google/uuid"

// Account data...
type Account struct {
	UserName   string    `json:"userName,omitempty"`
	AccountID  uuid.UUID `json:"accountId,omitempty"`
}

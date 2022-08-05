package model

import "github.com/google/uuid"

// Document data...
type Document struct {
	DocumentID          uuid.UUID `json:"documentId,omitempty"`
	Description 		string    `json:"description,omitempty"`
	HREF        		string    `json:"href,omitempty"`
	Filename    		uuid.UUID `json:"fileName,omitempty"`
	AccountID       	uuid.UUID `json:"accountId,omitempty"`
}

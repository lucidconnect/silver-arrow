package models

import "github.com/google/uuid"

type Token struct {
	ID          uuid.UUID
	Name        string
	Chain       int64
	Address     string
	MinorFactor int64
}

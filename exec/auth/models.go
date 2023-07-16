//	authd.go - Kahlehla project
//	Copyright 2023 by Kurt Duncan
//	Models for facilitating authentication

package auth

import (
	"github.com/google/uuid"
)

type Group struct {
	ID            uuid.UUID
	ParentGroupID uuid.UUID
	GroupName     string
	Members       []*uuid.UUID
}

type Account struct {
	ID           uuid.UUID
	LoginName    string
	PasswordHash string
	Disabled     bool
	UserName     string
	EmailAddress string
}

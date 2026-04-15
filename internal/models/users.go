package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`

	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

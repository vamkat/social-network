package models

import (
	"time"
)

// Session represents a user session based on JWT tokens.
// It can be used to store metadata for auditing or blacklisting tokens.
type Session struct {
	ID           int64      `db:"id" json:"id"`
	UserID       int64      `db:"user_id" json:"user_id"`
	AccessToken  string     `db:"access_token" json:"access_token"`
	RefreshToken string     `db:"refresh_token" json:"refresh_token"`
	ExpiresAt    time.Time  `db:"expires_at" json:"expires_at"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	RevokedAt    *time.Time `db:"revoked_at,omitempty" json:"revoked_at,omitempty"`
}

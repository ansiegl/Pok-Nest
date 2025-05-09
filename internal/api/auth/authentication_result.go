package auth

import (
	"time"

	"github.com/ansiegl/Pok-Nest.git/internal/models"
)

type AuthenticationResult struct {
	Token      string
	User       *models.User
	ValidUntil time.Time
	Scopes     []string
}

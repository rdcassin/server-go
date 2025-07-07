package auth

import (
	"time"
)

const (
	JWTTokenExpiration time.Duration = time.Hour
	RefreshTokenExpiration time.Duration = 60 * 24 * time.Hour
)
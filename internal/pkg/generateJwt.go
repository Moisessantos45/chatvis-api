package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Id        string `json:"id,omitempty"`
	Username  string `json:"username"`
	Useremail string `json:"useremail"`
	Token     string `json:"token,omitempty"`
	IsAdmin   bool   `json:"isAdmin"`
	jwt.RegisteredClaims
}

// GenerateJWT genera un nuevo token JWT para un usuario dado.
func GenerateJWT(id string, username string, email string, isAdmin bool) (string, error) {
	var JwtSecret = []byte(os.Getenv("SECRET_KEY_JWT"))

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &MyClaims{
		Id:        id,
		Username:  username,
		Useremail: email,
		IsAdmin:   isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "chatvist_chat-api",
			Subject:   email,
		},
	}

	newtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := newtoken.SignedString(JwtSecret)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err)
	}

	return tokenString, nil
}

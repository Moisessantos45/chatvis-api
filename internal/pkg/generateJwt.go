package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Define una clave secreta para firmar y verificar tus JWTs.
// ¡IMPORTANTE: En un entorno de producción, esto debería ser una variable de entorno segura, no hardcodeada!

// Claims personalizados que se añadirán a los claims estándar de JWT.
type MyClaims struct {
	Id        string `json:"id,omitempty"` // Opcional, si necesitas incluir el ID del usuario
	Username  string `json:"username"`
	Useremail string `json:"useremail"`
	Token     string `json:"token,omitempty"` // Opcional, si necesitas incluir un token específico
	IsAdmin   bool   `json:"isAdmin"`         // Indica si el usuario es administrador
	jwt.RegisteredClaims
}

// GenerateJWT genera un nuevo token JWT para un usuario dado.
func GenerateJWT(id string, username string, email string, isAdmin bool) (string, error) {
	var JwtSecret = []byte(os.Getenv("SECRET_KEY_JWT"))

	// Define la fecha de expiración del token (ej. 24 horas a partir de ahora)
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
			Issuer:    "chatvist_chat-api", // Quién emite el token
			Subject:   email,               // A quién pertenece el token
		},
	}

	// Crea el token con el algoritmo de firma y los claims
	newtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firma el token con la clave secreta
	tokenString, err := newtoken.SignedString(JwtSecret)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err)
	}

	return tokenString, nil
}

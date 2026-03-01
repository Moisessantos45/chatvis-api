package middleware

import (
	"chatvis-chat/internal/pkg"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

// JWTAuthMiddleware is a Fiber middleware to validate JWT tokens.

func JWTAuthMiddleware() fiber.Handler {

	var JwtSecret = []byte(os.Getenv("SECRET_KEY_JWT"))

	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")

		if authHeader == "" {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authentication token not provided"})

		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format. Expected 'Bearer <token>'"})

		}

		tokenString := parts[1]

		claims := &pkg.MyClaims{}

		// Parse and validate the token with the custom claims struct

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])

			}

			return JwtSecret, nil

		})

		if err != nil {

			log.Printf("Error parsing or validating the token: %v", err)

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})

		}

		// Check if the token is valid

		if !token.Valid {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})

		}

		// log.Println("Token is valid")

		// log.Println("UserID from claims:", claims.Id)

		// log.Println("Useremail from claims:", claims.Useremail)

		// log.Println("Username from claims:", claims.Username)

		// Token is valid, now we can access the claims

		c.Locals("userId", claims.Id)

		c.Locals("email", claims.Useremail)

		c.Locals("username", claims.Username)

		c.Locals("isAdmin", claims.IsAdmin)

		return c.Next()

	}

}

// AdminAuthMiddleware is a Fiber middleware to check if the user is an admin
func AdminAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		isAdmin, ok := c.Locals("isAdmin").(bool)
		if !ok || !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Access denied. Admin privileges required."})
		}
		return c.Next()
	}
}

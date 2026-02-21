package routers

import (
	"chatvis-chat/internal/app/auth/handleauth"

	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(group fiber.Router, handle *handleauth.AuthController) {
	group.Post("/login", handle.Login)
}

func RegisterAuthRoutesWithMiddleware(group fiber.Router, handle *handleauth.AuthController) {
	group.Post("/logout", handle.Logout)
}

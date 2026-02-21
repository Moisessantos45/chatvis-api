package http

import (
	"chatvis-chat/internal/models"
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	AUsecase domain.AuthUseCase
}

func NewAuthHandler(group fiber.Router, au domain.AuthUseCase) {
	handler := &AuthHandler{
		AUsecase: au,
	}

	group.Post("/login", handler.Login)
}

func NewAuthProtectedHandler(group fiber.Router, au domain.AuthUseCase) {
	handler := &AuthHandler{
		AUsecase: au,
	}

	group.Post("/logout", handler.Logout)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	// Reutilizamos el struct UsuarioLogin anterior si existe, o definimos uno local:
	var reqBody models.UsuarioLogin

	if err := c.BodyParser(&reqBody); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al iniciar sesión", "Error de parseo", err.Error())
	}

	usuario, err := h.AUsecase.Authenticate(reqBody.Email, reqBody.Password)

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusUnauthorized, "Error al iniciar sesión", "Credenciales inválidas", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Inicio de sesión exitoso", "", usuario)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {

	email := c.Locals("email").(string)

	if len(strings.TrimSpace(email)) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al cerrar sesión", "Error de parámetro", "Email no proporcionado")
	}

	if err := h.AUsecase.Logout(email); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al cerrar sesión", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Cierre de sesión exitoso", "", nil)
}

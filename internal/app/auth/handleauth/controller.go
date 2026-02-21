package handleauth

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/auth/serviceauth"
	"chatvis-chat/internal/pkg"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	Service *serviceauth.ServiceAuth
}

func (h *AuthController) Login(c *fiber.Ctx) error {
	var reqBody app.UsuarioLogin

	if err := c.BodyParser(&reqBody); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al iniciar sesión", "Error de parseo", err.Error())
	}

	usuario, err := h.Service.Authenticate(reqBody.Email, reqBody.Password)

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusUnauthorized, "Error al iniciar sesión", "Credenciales inválidas", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Inicio de sesión exitoso", "", usuario)
}

func (h *AuthController) Logout(c *fiber.Ctx) error {

	email := c.Locals("email").(string)

	if len(strings.TrimSpace(email)) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al cerrar sesión", "Error de parámetro", "Email no proporcionado")
	}

	if err := h.Service.Logout(email); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al cerrar sesión", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Cierre de sesión exitoso", "", nil)
}

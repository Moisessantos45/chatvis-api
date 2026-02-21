package handleusuario

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/usuario/serviceusuario"
	"chatvis-chat/internal/pkg"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UsuarioController struct {
	Service *serviceusuario.UsuarioUseCase
}

func (h *UsuarioController) GetUsuarioByID(c *fiber.Ctx) error {

	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario", "Error parametro", err.Error())
	}

	usuario, err := h.Service.GetById(id)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", usuario)
}

func (h *UsuarioController) GetUsuarioByEmail(c *fiber.Ctx) error {

	email := c.Params("email")
	if len(strings.TrimSpace(email)) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario", "Error parametro", "Email no proporcionado")
	}

	usuario, err := h.Service.GetByEmail(email)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", usuario)
}

func (h *UsuarioController) GetUsuarioByEmailToken(c *fiber.Ctx) error {

	localEmail := c.Locals("email")

	// This log will show if the value is nil or a string
	log.Printf("Raw value from c.Locals('email'): %+v", localEmail)

	if localEmail == nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error", "Error", "Email not found in context. Is the JWT middleware applied?")
	}

	email, ok := localEmail.(string)

	if !ok {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error", "Error", "Context email value is not a string.")
	}

	if len(strings.TrimSpace(email)) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario", "Error parametro", "Email no proporcionados")
	}

	usuario, err := h.Service.GetByEmail(email)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", usuario)
}

func (h *UsuarioController) CreateUsuario(c *fiber.Ctx) error {

	var usuario app.Usuarios
	if err := c.BodyParser(&usuario); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al crear usuario", "Error de parseo", err.Error())
	}

	if err := h.Service.Create(&usuario); err != nil {
		log.Println("Error al crear usuario:", err)
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	return pkg.ResponseJson(c, fiber.StatusCreated, "Usuario creado correctamente", "", usuario)
}

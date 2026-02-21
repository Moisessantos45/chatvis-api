package http

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type UsuarioHandler struct {
	UUsecase domain.UsuarioUseCase
}

// NewUsuarioHandler registra todos los endpoints para el módulo de usuario en fiber.Router
func NewUsuarioHandler(group fiber.Router, uu domain.UsuarioUseCase) {
	handler := &UsuarioHandler{
		UUsecase: uu,
	}

	group.Get("/correo/:email", handler.GetUsuarioByEmail) // /api/usuario/correo/:email
	group.Get("/token", handler.GetUsuarioByEmailToken)    // /api/usuario/token
	group.Get("/:id", handler.GetUsuarioByID)              // /api/usuario/:id
	group.Post("/", handler.CreateUsuario)                 // /api/usuario/
}

// NewUsuarioPublicHandler registra endpoints de registro/login sin auth requerida
func NewUsuarioPublicHandler(group fiber.Router, uu domain.UsuarioUseCase) {
	handler := &UsuarioHandler{
		UUsecase: uu,
	}
	group.Post("/register", handler.CreateUsuario)
}

func (h *UsuarioHandler) GetUsuarioByID(c *fiber.Ctx) error {

	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario", "Error parametro", err.Error())
	}

	usuario, err := h.UUsecase.GetById(id)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", usuario)
}

func (h *UsuarioHandler) GetUsuarioByEmail(c *fiber.Ctx) error {

	email := c.Params("email")
	if len(strings.TrimSpace(email)) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario", "Error parametro", "Email no proporcionado")
	}

	usuario, err := h.UUsecase.GetByEmail(email)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", usuario)
}

func (h *UsuarioHandler) GetUsuarioByEmailToken(c *fiber.Ctx) error {

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

	usuario, err := h.UUsecase.GetByEmail(email)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", usuario)
}

func (h *UsuarioHandler) CreateUsuario(c *fiber.Ctx) error {

	var usuario domain.Usuario
	if err := c.BodyParser(&usuario); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al crear usuario", "Error de parseo", err.Error())
	}

	if err := h.UUsecase.Create(&usuario); err != nil {
		log.Println("Error al crear usuario:", err)
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	// Regresar status 201 Created
	return pkg.ResponseJson(c, fiber.StatusCreated, "Usuario creado correctamente", "", usuario)
}

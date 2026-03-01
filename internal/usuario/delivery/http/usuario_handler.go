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

	group.Get("/correo/:email", handler.GetUsuarioByEmail)
	group.Get("/token", handler.GetUsuarioByEmailToken)
	group.Get("/:id", handler.GetUsuarioByID)
	group.Post("/", handler.CreateUsuario)
}

// NewUsuarioPublicHandler registra endpoints de registro/login sin auth requerida
func NewUsuarioPublicHandler(group fiber.Router, uu domain.UsuarioUseCase) {
	handler := &UsuarioHandler{
		UUsecase: uu,
	}
	group.Post("/register", handler.CreateUsuario)
}

// NewAdminUsuarioHandler registra endpoints solo para admin
func NewAdminUsuarioHandler(group fiber.Router, uu domain.UsuarioUseCase) {
	handler := &UsuarioHandler{
		UUsecase: uu,
	}
	group.Get("/user", handler.GetAllUsuarios)
	group.Post("/user", handler.CreateAdminUsuario)
	group.Patch("/user/:id/deactivate", handler.DeactivateUsuario)
	group.Post("/user/:id/remove-sessions", handler.RemoveSesionesUsuario)
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

	// En el registro normal no se permite ser admin ni llm
	usuario.IsAdmin = false
	usuario.IsLlm = false

	if err := h.UUsecase.Create(&usuario); err != nil {
		log.Println("Error al crear usuario:", err)
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear usuario", "Error interno", err.Error())
	}

	usuario.Password = "" // No enviar la contraseña en la respuesta

	// Regresar status 201 Created
	return pkg.ResponseJson(c, fiber.StatusCreated, "Usuario creado correctamente", "", usuario)
}

func (h *UsuarioHandler) GetAllUsuarios(c *fiber.Ctx) error {
	usuarios, err := h.UUsecase.GetAllUsuarios()
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuarios", "Error interno", err.Error())
	}
	for i := range usuarios {
		usuarios[i].Password = ""
	}
	return pkg.ResponseJson(c, fiber.StatusOK, "Usuarios obtenidos correctamente", "", usuarios)
}

func (h *UsuarioHandler) CreateAdminUsuario(c *fiber.Ctx) error {
	var usuario domain.Usuario
	if err := c.BodyParser(&usuario); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al crear usuario", "Error de parseo", err.Error())
	}

	if err := h.UUsecase.Create(&usuario); err != nil {
		log.Println("Error al crear usuario (admin):", err)
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear usuario", "Error interno", err.Error())
	}
	usuario.Password = ""
	return pkg.ResponseJson(c, fiber.StatusCreated, "Usuario creado correctamente por admin", "", usuario)
}

func (h *UsuarioHandler) DeactivateUsuario(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error id", "Error parametro", err.Error())
	}
	if err := h.UUsecase.UpdateIsActive(id, false); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al desactivar", "Error interno", err.Error())
	}
	_ = h.UUsecase.ClearToken(id) // Clear sessions
	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario desactivado correctamente", "", nil)
}

func (h *UsuarioHandler) RemoveSesionesUsuario(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error id", "Error parametro", err.Error())
	}
	if err := h.UUsecase.ClearToken(id); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al remover sesion", "Error interno", err.Error())
	}
	return pkg.ResponseJson(c, fiber.StatusOK, "Sesiones removidas", "", nil)
}

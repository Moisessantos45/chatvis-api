package http

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type GrupoHandler struct {
	GUsecase domain.GrupoUseCase
}

func NewGrupoHandler(group fiber.Router, gu domain.GrupoUseCase) {
	handler := &GrupoHandler{
		GUsecase: gu,
	}

	group.Get("", handler.GetAllGrupos)
	group.Get("/key/:key", handler.GetGrupoByClave)
	group.Get("/user/:id", handler.GetAllGruposByUsuarioId)
	group.Get("/:id", handler.GetGrupoById)
	group.Post("", handler.CreateGrupo)
	group.Post("/generate-code/:id", handler.CreateInvitationUrl)
}

// NewAdminGrupoHandler registra endpoints de grupos para administradores
func NewAdminGrupoHandler(group fiber.Router, gu domain.GrupoUseCase) {
	handler := &GrupoHandler{
		GUsecase: gu,
	}
	group.Get("/group", handler.GetAllGrupos)
}

func (h *GrupoHandler) GetAllGrupos(c *fiber.Ctx) error {
	result, err := h.GUsecase.GetAll()
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupos", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupos obtenidos correctamente", "", result)
}

func (h *GrupoHandler) GetGrupoById(c *fiber.Ctx) error {
	grupoId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupo", "Error parametro", err.Error())
	}

	result, err := h.GUsecase.GetById(grupoId)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupo obtenido correctamente", "", result)
}

func (h *GrupoHandler) GetGrupoByClave(c *fiber.Ctx) error {
	clave := c.Params("key")
	if len(strings.TrimSpace(clave)) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupo", "Error parametro", "Clave no proporcionada")
	}

	result, err := h.GUsecase.GetByClave(clave)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupo obtenido correctamente", "", result)
}

func (h *GrupoHandler) GetAllGruposByUsuarioId(c *fiber.Ctx) error {
	usuarioId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupos", "Error parametro", err.Error())
	}

	result, err := h.GUsecase.GetAllByUsuarioId(usuarioId)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupos", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupos obtenidos correctamente", "", result)
}

func (h *GrupoHandler) CreateInvitationUrl(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al generar URL de invitación", "Error parametro", err.Error())
	}

	url, clave, err := h.GUsecase.CreateInvitationUrl(id)
	return pkg.ResponseJson(c, fiber.StatusOK, "URL de invitación generada correctamente", "", map[string]string{
		"url":   url,
		"clave": clave,
	})
}

func (h *GrupoHandler) CreateGrupo(c *fiber.Ctx) error {
	var grupo domain.Grupo
	if err := c.BodyParser(&grupo); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al crear grupo", "Error de parseo", err.Error())
	}

	err := h.GUsecase.Create(&grupo)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Grupo creado correctamente", "", grupo)
}

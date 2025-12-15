package handlegrupo

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/grupo/servicegrupo"
	"chatvis-chat/internal/pkg"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type GrupoController struct {
	Service *servicegrupo.GrupoService
}

func (h *GrupoController) GetAllGrupos(c *fiber.Ctx) error {
	result, err := h.Service.GetAll()

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupos", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupos obtenidos correctamente", "", result)
}

func (h *GrupoController) GetGrupoById(c *fiber.Ctx) error {
	grupoId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupo", "Error parametro", err.Error())
	}

	result, err := h.Service.GetById(grupoId)

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupo obtenido correctamente", "", result)
}

func (h *GrupoController) GetGrupoByClave(c *fiber.Ctx) error {
	clave := c.Params("clave")
	if len(strings.TrimSpace(clave)) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupo", "Error parametro", "Clave no proporcionada")
	}

	result, err := h.Service.GetByClave(clave)

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupo obtenido correctamente", "", result)
}

func (h *GrupoController) GetAllGruposByUsuarioId(c *fiber.Ctx) error {
	usuarioId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupos", "Error parametro", err.Error())
	}

	result, err := h.Service.GetAllByUsuarioId(usuarioId)

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupos", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupos obtenidos correctamente", "", result)
}

func (h *GrupoController) CreateGrupo(c *fiber.Ctx) error {

	var grupo app.Grupos
	if err := c.BodyParser(&grupo); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al crear grupo", "Error de parseo", err.Error())
	}

	err := h.Service.Create(&grupo)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Grupo creado correctamente", "", grupo)
}

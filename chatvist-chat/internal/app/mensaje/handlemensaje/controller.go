package handlemensaje

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/mensaje/servicemensaje"
	"chatvis-chat/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

type MensajeController struct {
	Service *servicemensaje.MensajeService
}

func (h *MensajeController) GetMensajeByID(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener mensaje", "Error parametro", err.Error())
	}

	mensaje, err := h.Service.GetById(id)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener mensaje", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Mensaje obtenido correctamente", "", mensaje)
}

func (h *MensajeController) GetMensajesByChatID(c *fiber.Ctx) error {
	grupoId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener mensajes", "Error parametro", err.Error())
	}

	mensajes, err := h.Service.GetAllByGrupoId(grupoId)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener mensajes", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Mensajes obtenidos correctamente", "", mensajes)
}

func (h *MensajeController) CreateMensaje(c *fiber.Ctx) error {
	var mensaje app.Mensajes
	if err := c.BodyParser(&mensaje); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al crear mensaje", "Error de parseo", err.Error())
	}

	newMensaje, err := h.Service.Create(&mensaje)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear mensaje", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Mensaje creado correctamente", "", newMensaje)
}

func (h *MensajeController) UpdateMensaje(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al actualizar mensaje", "Error parametro", err.Error())
	}

	var mensaje app.Mensajes
	if err := c.BodyParser(&mensaje); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al actualizar mensaje", "Error de parseo", err.Error())
	}

	if err := h.Service.Update(id, &mensaje); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al actualizar mensaje", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Mensaje actualizado correctamente", "", mensaje)
}

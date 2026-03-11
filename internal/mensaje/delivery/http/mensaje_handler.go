package http

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

type MensajeHandler struct {
	MUsecase domain.MensajeUseCase
}

func NewMensajeHandler(group fiber.Router, mu domain.MensajeUseCase) {
	handler := &MensajeHandler{
		MUsecase: mu,
	}

	group.Get("/group/:id", handler.GetMensajesByChatID)
	group.Get("/group/clave/:clave", handler.GetMensajesByChatClave)
	group.Get("/:id", handler.GetMensajeByID)
	group.Post("/", handler.CreateMensaje)
	group.Put("/:id", handler.UpdateMensaje)
}

func (h *MensajeHandler) GetMensajeByID(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener mensaje", "Error parametro", err.Error())
	}

	mensaje, err := h.MUsecase.GetById(id)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener mensaje", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Mensaje obtenido correctamente", "", mensaje)
}

func (h *MensajeHandler) GetMensajesByChatID(c *fiber.Ctx) error {
	grupoId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener mensajes", "Error parametro", err.Error())
	}

	fechaInicioStr := c.Query("fechaInicio", "")
	fechaFinStr := c.Query("fechaFin", "")

	var fechaInicio, fechaFin time.Time
	if fechaInicioStr != "" {
		fechaInicio, err = time.Parse(time.RFC3339, fechaInicioStr)
		if err != nil {
			return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener mensajes", "Error parametro", "La fecha de inicio debe estar en formato RFC3339")
		}
	}

	if fechaFinStr != "" {
		fechaFin, err = time.Parse(time.RFC3339, fechaFinStr)
		if err != nil {
			return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener mensajes", "Error parametro", "La fecha de fin debe estar en formato RFC3339")
		}
	}

	log.Println("Grupo ID:", grupoId, "Filtro:", fechaInicio.Format(time.RFC3339), fechaFin.Format(time.RFC3339))

	mensajes, err := h.MUsecase.GetAllByGrupoId(grupoId, fechaInicio, fechaFin)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener mensajes", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Mensajes obtenidos correctamente", "", mensajes)
}

func (h *MensajeHandler) GetMensajesByChatClave(c *fiber.Ctx) error {
	clave := c.Params("clave")
	if len(clave) == 0 {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener mensajes", "Error parametro", "La clave del grupo es requerida")
	}

	mensajes, err := h.MUsecase.GetAllByGrupoClave(clave)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener mensajes", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Mensajes obtenidos correctamente", "", mensajes)
}

func (h *MensajeHandler) CreateMensaje(c *fiber.Ctx) error {
	var mensaje domain.Mensaje
	if err := c.BodyParser(&mensaje); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al crear mensaje", "Error de parseo", err.Error())
	}

	newMensaje, err := h.MUsecase.Create(&mensaje)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear mensaje", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Mensaje creado correctamente", "", newMensaje)
}

func (h *MensajeHandler) UpdateMensaje(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al actualizar mensaje", "Error parametro", err.Error())
	}

	var mensaje domain.Mensaje
	if err := c.BodyParser(&mensaje); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al actualizar mensaje", "Error de parseo", err.Error())
	}

	if err := h.MUsecase.Update(id, &mensaje); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al actualizar mensaje", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Mensaje actualizado correctamente", "", mensaje)
}

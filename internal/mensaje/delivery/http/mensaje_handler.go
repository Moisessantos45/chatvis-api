package http

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

type MensajeHandler struct {
	MUsecase domain.MensajeUseCase
}

func NewMensajeHandler(group fiber.Router, mu domain.MensajeUseCase) {
	handler := &MensajeHandler{
		MUsecase: mu,
	}

	group.Get("/grupo/:id", handler.GetMensajesByChatID) // Obtener mensajes por ID de grupo
	group.Get("/:id", handler.GetMensajeByID)            // Obtener mensaje por ID
	group.Post("/", handler.CreateMensaje)               // Crear un nuevo mensaje
	group.Put("/:id", handler.UpdateMensaje)             // Actualizar un mensaje existente
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

	mensajes, err := h.MUsecase.GetAllByGrupoId(grupoId)
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

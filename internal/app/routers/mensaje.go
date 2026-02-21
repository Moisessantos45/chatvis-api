package routers

import (
	"chatvis-chat/internal/app/mensaje/handlemensaje"

	"github.com/gofiber/fiber/v2"
)

func RegisterMensajeRoutes(group fiber.Router, handler *handlemensaje.MensajeController) {

	group.Get("/grupo/:id", handler.GetMensajesByChatID) // Obtener mensajes por ID de grupo
	group.Get("/:id", handler.GetMensajeByID)            // Obtener mensaje por ID
	group.Post("/", handler.CreateMensaje)               // Crear un nuevo mensaje
	group.Put("/:id", handler.UpdateMensaje)             // Actualizar un mensaje existente
}

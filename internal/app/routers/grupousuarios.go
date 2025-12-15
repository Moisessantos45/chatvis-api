package routers

import (
	"chatvis-chat/internal/app/grupousario/handlegrupousuario"

	"github.com/gofiber/fiber/v2"
)

func RegisterGrupoUsuarioRoutes(group fiber.Router, handle *handlegrupousuario.GruposUsuariosController) {
	group.Get("/:id", handle.GetByGrupoId)
	group.Get("/groups/:userId", handle.GetByUsuarioId)
	group.Post("/add", handle.AddUserToGroup)
}

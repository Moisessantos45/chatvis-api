package routers

import (
	"chatvis-chat/internal/app/grupo/handlegrupo"

	"github.com/gofiber/fiber/v2"
)

func RegisterGrupoRoutes(group fiber.Router, handle *handlegrupo.GrupoController) {

	group.Get("", handle.GetAllGrupos)                        // Obtener todos los grupos
	group.Get("/clave/:clave", handle.GetGrupoByClave)        // Obtener grupo por clave
	group.Get("/usuario/:id", handle.GetAllGruposByUsuarioId) // Obtener grupos por ID de usuario
	group.Get("/:id", handle.GetGrupoById)                    // Obtener grupo por ID
	group.Post("", handle.CreateGrupo)                        // Crear un nuevo grupo
}

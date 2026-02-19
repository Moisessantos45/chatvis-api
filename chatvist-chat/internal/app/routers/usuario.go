package routers

import (
	"chatvis-chat/internal/app/usuario/handleusuario"

	"github.com/gofiber/fiber/v2"
)

func RegisterUsuarioRoutes(group fiber.Router, handler *handleusuario.UsuarioController) {
	group.Get("/correo", handler.GetUsuarioByEmail)     // Obtener usuario por email
	group.Get("/token", handler.GetUsuarioByEmailToken) // Obtener usuario por email desde el token
	group.Get("/:id", handler.GetUsuarioByID)           // Obtener usuario por ID
	group.Post("/", handler.CreateUsuario)              // Crear un nuevo usuario
}

func RegisterUsuarioRoutesBasicNotAuth(group fiber.Router, handler *handleusuario.UsuarioController) {
	group.Post("/register", handler.CreateUsuario) // Crear un nuevo usuario
}

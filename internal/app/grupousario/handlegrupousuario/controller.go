package handlegrupousuario

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/grupousario/servicegrupousuario"
	"chatvis-chat/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

type GruposUsuariosController struct {
	Service *servicegrupousuario.GruposUsuariosService
}

func (h *GruposUsuariosController) GetByGrupoId(c *fiber.Ctx) error {

	grupoId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupo", "Error parametro", err.Error())
	}

	result, err := h.Service.GetByGrupoId(grupoId)

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupo obtenido correctamente", "", result)
}

func (h *GruposUsuariosController) GetByUsuarioId(c *fiber.Ctx) error {
	usuarioId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario", "Error parametro", err.Error())
	}

	result, err := h.Service.GetByGrupoUsuarioId(usuarioId)

	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", result)
}

func (h *GruposUsuariosController) AddUserToGroup(c *fiber.Ctx) error {
	var grupoUsuario app.GrupoWithUsuario

	if err := c.BodyParser(&grupoUsuario); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al parsear el cuerpo de la solicitud", "Error de formato", err.Error())
	}

	if err := h.Service.AddUserToGroup(grupoUsuario.Clave, grupoUsuario.UsuarioId); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear grupo de usuario", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Grupo de usuario creado correctamente", "", grupoUsuario)
}

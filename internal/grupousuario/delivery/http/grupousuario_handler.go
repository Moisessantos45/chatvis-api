package http

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"

	"github.com/gofiber/fiber/v2"
)

type GrupoUsuarioHandler struct {
	Usecase domain.GrupoUsuarioUseCase
}

func NewGrupoUsuarioHandler(group fiber.Router, uc domain.GrupoUsuarioUseCase) {
	handler := &GrupoUsuarioHandler{
		Usecase: uc,
	}

	group.Get("/:id", handler.GetByGrupoId)
	group.Get("/group/:userId", handler.GetByUsuarioId)
	group.Post("/add", handler.AddUserToGroup)
}

func NewAdminGrupoUsuarioHandler(group fiber.Router, uc domain.GrupoUsuarioUseCase) {
	handler := &GrupoUsuarioHandler{
		Usecase: uc,
	}

	group.Post("/group-users", handler.AdminAddUserToGroup)
	group.Post("/group/:id/user", handler.AdminAddUserToGroup)
}

func (h *GrupoUsuarioHandler) GetByGrupoId(c *fiber.Ctx) error {
	grupoId, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener grupo", "Error parametro", err.Error())
	}

	result, err := h.Usecase.GetUsersByGroupId(grupoId)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener grupo", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Grupo obtenido correctamente", "", result)
}

func (h *GrupoUsuarioHandler) GetByUsuarioId(c *fiber.Ctx) error {
	usuarioId, err := c.ParamsInt("userId")
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario", "Error parametro", err.Error())
	}

	result, err := h.Usecase.GetByUsuarioId(uint64(usuarioId))
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al obtener usuario", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusOK, "Usuario obtenido correctamente", "", result)
}

func (h *GrupoUsuarioHandler) AddUserToGroup(c *fiber.Ctx) error {
	var body struct {
		Clave     string `json:"clave"`
		UsuarioId uint64 `json:"usuarioId"`
	}

	if err := c.BodyParser(&body); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al parsear el cuerpo de la solicitud", "Error de formato", err.Error())
	}

	if err := h.Usecase.JoinGroup(body.UsuarioId, body.Clave); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al crear grupo de usuario", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Grupo de usuario creado correctamente", "", body)
}

func (h *GrupoUsuarioHandler) AdminAddUserToGroup(c *fiber.Ctx) error {
	id, err := pkg.ValidateParamsId(c)
	if err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al obtener usuario id", "Error parametro", err.Error())
	}

	var body struct {
		Clave string `json:"clave"`
	}

	if err := c.BodyParser(&body); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al parsear el cuerpo de la solicitud", "Error de formato", err.Error())
	}

	if err := h.Usecase.JoinGroup(id, body.Clave); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al asignar grupo a usuario", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Grupo asignado correctamente al usuario", "", nil)
}

func (h *GrupoUsuarioHandler) AdminAddUsersToGroup(c *fiber.Ctx) error {

	var body struct {
		GroupsIds []uint64 `json:"groupsIds"`
		UsersIds  []uint64 `json:"usersIds"`
	}

	if err := c.BodyParser(&body); err != nil {
		return pkg.ResponseJson(c, fiber.StatusBadRequest, "Error al parsear el cuerpo de la solicitud", "Error de formato", err.Error())
	}

	if err := h.Usecase.JoinGroups(body.UsersIds, body.GroupsIds); err != nil {
		return pkg.ResponseJson(c, fiber.StatusInternalServerError, "Error al asignar grupo a usuario", "Error interno", err.Error())
	}

	return pkg.ResponseJson(c, fiber.StatusCreated, "Grupos asignados correctamente a los usuarios", "", nil)
}

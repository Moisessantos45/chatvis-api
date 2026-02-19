package servicegrupousuario

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/grupo"
	"chatvis-chat/internal/app/grupousario"
	"errors"
	"log"
	"strings"
)

type GruposUsuariosService struct {
	repo      grupousario.GruposUsuariosRepository
	repoGroup grupo.GrupoRepository
}

func NewGruposUsuariosServiceClient(repo grupousario.GruposUsuariosRepository, repoG grupo.GrupoRepository) *GruposUsuariosService {
	return &GruposUsuariosService{
		repo:      repo,
		repoGroup: repoG,
	}
}

func (s *GruposUsuariosService) GetByGrupoId(grupoId uint64) ([]app.GruposUsuarios, error) {
	if grupoId <= 0 {
		return nil, errors.New("El ID del grupo debe ser mayor que cero")
	}

	return s.repo.GetByGrupoId(grupoId)
}

func (s *GruposUsuariosService) GetByGrupoUsuarioId(usuarioId uint64) (*app.GruposUsuarios, error) {
	if usuarioId <= 0 {
		return nil, errors.New("El ID del usuario debe ser mayor que cero")
	}

	return s.repo.GetByGrupoUsuarioId(usuarioId)
}

func (s *GruposUsuariosService) Create(grupoUsuario *app.GruposUsuarios) error {
	if grupoUsuario == nil {
		return errors.New("El grupo de usuario no puede ser nulo")
	}

	if grupoUsuario.IdGrupo <= 0 || grupoUsuario.IdUsuario <= 0 {
		return errors.New("Los IDs del grupo y del usuario deben ser mayores que cero")
	}

	return s.repo.Create(grupoUsuario)
}

func (s *GruposUsuariosService) AddUserToGroup(clave string, id uint64) error {
	log.Println("AddUserToGroup - Clave:", clave, "ID Usuario:", id)
	if len(strings.TrimSpace(clave)) == 0 {
		return errors.New("La clave no puede estar vacía")
	}

	if id <= 0 {
		return errors.New("El ID del usuario debe ser mayor que cero")
	}

	existsGroup, code, err := s.repoGroup.GetByClave(clave)
	log.Println("Grupo encontrado:", existsGroup, "Código:", code, "Error:", err)

	// Manejar errores de la base de datos
	if err != nil {
		return err // Esto capturará cualquier error real de la DB, como una conexión fallida
	}

	// Si no se encontró el grupo (código 404)
	if code == 404 {
		return errors.New("El grupo no existe con la clave proporcionada")
	}

	newGroupWithUser := &app.GruposUsuarios{
		IdGrupo:   existsGroup.Id,
		IdUsuario: id,
	}

	return s.repo.Create(newGroupWithUser)
}

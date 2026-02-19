package servicegrupo

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/grupo"
	"chatvis-chat/internal/pkg"
	"errors"
	"strings"
	"time"
)

type GrupoService struct {
	repo grupo.GrupoRepository
}

func NewGrupoService(repo grupo.GrupoRepository) *GrupoService {
	return &GrupoService{
		repo: repo,
	}
}

func (s *GrupoService) GetAll() ([]app.Grupos, error) {

	grupos, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.New("Error al obtener los grupos: " + err.Error())
	}

	return grupos, nil
}

func (s *GrupoService) GetById(id uint64) (*app.Grupos, error) {

	if id <= 0 {
		return nil, errors.New("El ID del grupo debe ser mayor que cero")
	}

	return s.repo.GetById(id)
}

func (s *GrupoService) GetByClave(clave string) (*app.Grupos, error) {

	if len(strings.TrimSpace(clave)) == 0 {
		return nil, errors.New("La clave del grupo no puede estar vacía")
	}

	group, _, err := s.repo.GetByClave(clave)

	return group, err
}

func (s *GrupoService) GetAllByUsuarioId(usuarioId uint64) ([]app.Grupos, error) {

	if usuarioId <= 0 {
		return nil, errors.New("El ID del usuario debe ser mayor que cero")
	}

	return s.repo.GetAllByUsuarioId(usuarioId)
}

func (s *GrupoService) GetAllGruposByUsuarioIdToIds(usuarioId uint64) ([]string, error) {

	if usuarioId <= 0 {
		return nil, errors.New("El ID del usuario debe ser mayor que cero")
	}

	claves, err := s.repo.GetAllGruposByUsuarioIdToClaves(usuarioId)
	if err != nil {
		return nil, errors.New("Error al obtener los grupos por ID de usuario: " + err.Error())
	}

	return claves, nil
}

func (s *GrupoService) Create(grupo *app.Grupos) error {
	if grupo == nil {
		return errors.New("El grupo no puede ser nulo")
	}

	if len(strings.TrimSpace(grupo.Nombre)) == 0 {
		return errors.New("El nombre del grupo no puede estar vacío")
	}

	if grupo.CreatedById <= 0 {
		return errors.New("El ID del creador del grupo debe ser mayor que cero")
	}

	grupo.Clave = pkg.GenerateUUID()
	grupo.Fecha = time.Now()

	return s.repo.Create(grupo)
}

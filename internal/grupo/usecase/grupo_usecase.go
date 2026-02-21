package usecase

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
	"errors"
	"strings"
	"time"
)

type grupoUseCase struct {
	repo domain.GrupoRepository
}

func NewGrupoUseCase(repo domain.GrupoRepository) domain.GrupoUseCase {
	return &grupoUseCase{
		repo: repo,
	}
}

func (s *grupoUseCase) GetAll() ([]domain.Grupo, error) {
	grupos, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.New("error al obtener los grupos: " + err.Error())
	}
	return grupos, nil
}

func (s *grupoUseCase) GetById(id uint64) (*domain.Grupo, error) {
	if id <= 0 {
		return nil, errors.New("el ID del grupo debe ser mayor que cero")
	}
	return s.repo.GetById(id)
}

func (s *grupoUseCase) GetByClave(clave string) (*domain.Grupo, error) {
	if len(strings.TrimSpace(clave)) == 0 {
		return nil, errors.New("la clave del grupo no puede estar vacía")
	}

	group, _, err := s.repo.GetByClave(clave)
	return group, err
}

func (s *grupoUseCase) GetAllByUsuarioId(usuarioId uint64) ([]domain.Grupo, error) {
	if usuarioId <= 0 {
		return nil, errors.New("el ID del usuario debe ser mayor que cero")
	}
	return s.repo.GetAllByUsuarioId(usuarioId)
}

func (s *grupoUseCase) GetAllGruposByUsuarioIdToClaves(usuarioId uint64) ([]string, error) {
	if usuarioId <= 0 {
		return nil, errors.New("el ID del usuario debe ser mayor que cero")
	}

	claves, err := s.repo.GetAllGruposByUsuarioIdToClaves(usuarioId)
	if err != nil {
		return nil, errors.New("error al obtener los grupos por ID de usuario: " + err.Error())
	}

	return claves, nil
}

func (s *grupoUseCase) GetByName(name string) (*domain.Grupo, error) {
	if len(strings.TrimSpace(name)) == 0 {
		return nil, errors.New("el nombre no puede estar vacío")
	}

	group, _, err := s.repo.GetByName(name)
	return group, err
}

func (s *grupoUseCase) Create(grupo *domain.Grupo) error {
	if grupo == nil {
		return errors.New("el grupo no puede ser nulo")
	}

	if len(strings.TrimSpace(grupo.Nombre)) == 0 {
		return errors.New("el nombre del grupo no puede estar vacío")
	}

	if grupo.CreatedById <= 0 {
		return errors.New("el ID del creador del grupo debe ser mayor que cero")
	}

	grupo.Clave = pkg.GenerateUUID()
	grupo.Fecha = time.Now()

	return s.repo.Create(grupo)
}

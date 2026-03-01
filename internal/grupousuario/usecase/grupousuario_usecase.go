package usecase

import (
	"chatvis-chat/internal/domain"
	"errors"
	"log"
	"strings"
)

type grupoUsuarioUseCase struct {
	repo      domain.GrupoUsuarioRepository
	repoGrupo domain.GrupoRepository
}

func NewGrupoUsuarioUseCase(repo domain.GrupoUsuarioRepository, repoGrupo domain.GrupoRepository) domain.GrupoUsuarioUseCase {
	return &grupoUsuarioUseCase{
		repo:      repo,
		repoGrupo: repoGrupo,
	}
}

func (u *grupoUsuarioUseCase) GetUsersByGroupId(grupoId uint64) ([]domain.GrupoUsuario, error) {
	if grupoId <= 0 {
		return nil, errors.New("el ID del grupo debe ser mayor que cero")
	}

	return u.repo.GetByGrupoId(grupoId)
}

func (u *grupoUsuarioUseCase) GetByUsuarioId(usuarioId uint64) (*domain.GrupoUsuario, error) {
	if usuarioId <= 0 {
		return nil, errors.New("el ID del usuario debe ser mayor que cero")
	}

	return u.repo.GetByUsuarioId(usuarioId)
}

func (u *grupoUsuarioUseCase) JoinGroup(userId uint64, clave string) error {
	log.Println("JoinGroup - Clave:", clave, "ID Usuario:", userId)

	if len(strings.TrimSpace(clave)) == 0 {
		return errors.New("la clave no puede estar vacía")
	}

	if userId <= 0 {
		return errors.New("el ID del usuario debe ser mayor que cero")
	}

	// 1. Validar que el grupo exista con esa clave
	existsGroup, statusCode, err := u.repoGrupo.GetByClave(clave)

	if err != nil {
		return err // Error de DB
	}

	if statusCode == 404 || existsGroup == nil {
		return errors.New("el grupo no existe con la clave proporcionada")
	}

	// 2. Crear relación
	newGroupUser := &domain.GrupoUsuario{
		IdGrupo:   existsGroup.Id,
		IdUsuario: userId,
	}

	return u.repo.Create(newGroupUser)
}

func (u *grupoUsuarioUseCase) JoinGroups(usersIds []uint64, gruposIds []uint64) error {
	log.Println("JoinGroups - IDs Usuarios:", usersIds)

	if len(usersIds) == 0 {
		return errors.New("la lista de IDs de usuarios debe tener al menos un elemento")
	}


	if len(gruposIds) == 0 {
		return errors.New("la lista de IDs de grupos debe tener al menos un elemento")
	}

	// 2. Crear relaciones
	for _, userId := range usersIds {
		for _, groupId := range gruposIds {
			newGroupUser := &domain.GrupoUsuario{
				IdGrupo:   groupId,
				IdUsuario: userId,
			}

			if err := u.repo.Create(newGroupUser); err != nil {
				return err
			}
		}
	}

	return nil
}

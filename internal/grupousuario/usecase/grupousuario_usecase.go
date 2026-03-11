package usecase

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
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

func (u *grupoUsuarioUseCase) VerifyMembership(userId uint64, clave string) (bool, error) {
	log.Println("VerifyMembership - Clave:", clave, "ID Usuario:", userId)

	if len(strings.TrimSpace(clave)) == 0 {
		return false, errors.New("la clave no puede estar vacía")
	}

	extractClave := pkg.Base58ToUuid(clave)
	log.Println("Clave extraída:", extractClave)

	if userId <= 0 {
		return false, errors.New("el ID del usuario debe ser mayor que cero")
	}

	return u.repo.VerifyMembership(userId, extractClave)
}

func (u *grupoUsuarioUseCase) JoinGroup(userId uint64, clave string) error {
	log.Println("JoinGroup - Clave:", clave, "ID Usuario:", userId)

	if len(strings.TrimSpace(clave)) == 0 {
		return errors.New("la clave no puede estar vacía")
	}

	extractClave := pkg.Base58ToUuid(clave)
	log.Println("Clave extraída:", extractClave)

	if userId <= 0 {
		return errors.New("el ID del usuario debe ser mayor que cero")
	}

	existsGroup, statusCode, err := u.repoGrupo.GetByClave(extractClave)

	if err != nil {
		return err
	}

	if statusCode == 404 || existsGroup == nil {
		return errors.New("el grupo no existe con la clave proporcionada")
	}

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

package repository

import (
	"chatvis-chat/internal/models"
	"chatvis-chat/internal/domain"

	"gorm.io/gorm"
)

type postgresGrupoUsuarioRepository struct {
	db *gorm.DB
}

func NewPostgresGrupoUsuarioRepository(db *gorm.DB) domain.GrupoUsuarioRepository {
	return &postgresGrupoUsuarioRepository{db: db}
}

func mapGormToDomainGrupoUsuario(gu *models.GruposUsuarios) *domain.GrupoUsuario {
	if gu == nil {
		return nil
	}
	return &domain.GrupoUsuario{
		IdGrupo:   gu.IdGrupo,
		IdUsuario: gu.IdUsuario,
	}
}

func mapDomainToGormGrupoUsuario(gu *domain.GrupoUsuario) *models.GruposUsuarios {
	if gu == nil {
		return nil
	}
	return &models.GruposUsuarios{
		IdGrupo:   gu.IdGrupo,
		IdUsuario: gu.IdUsuario,
	}
}

func (r *postgresGrupoUsuarioRepository) GetByGrupoId(grupoId uint64) ([]domain.GrupoUsuario, error) {
	var gormGruposUsuarios []models.GruposUsuarios
	if err := r.db.Where("id_grupo = ?", grupoId).Find(&gormGruposUsuarios).Error; err != nil {
		return nil, err
	}

	var res []domain.GrupoUsuario
	for _, gu := range gormGruposUsuarios {
		res = append(res, *mapGormToDomainGrupoUsuario(&gu))
	}
	return res, nil
}

func (r *postgresGrupoUsuarioRepository) GetByUsuarioId(usuarioId uint64) (*domain.GrupoUsuario, error) {
	var grupoUsuario models.GruposUsuarios
	if err := r.db.Where("id_usuario = ?", usuarioId).First(&grupoUsuario).Error; err != nil {
		return nil, err
	}
	return mapGormToDomainGrupoUsuario(&grupoUsuario), nil
}

func (r *postgresGrupoUsuarioRepository) Create(grupoUsuario *domain.GrupoUsuario) error {
	gormGu := mapDomainToGormGrupoUsuario(grupoUsuario)
	if err := r.db.Create(gormGu).Error; err != nil {
		return err
	}
	return nil
}

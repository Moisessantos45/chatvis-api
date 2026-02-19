package grupousario

import (
	"chatvis-chat/internal/app"

	"gorm.io/gorm"
)

type GruposUsuariosRepository struct {
	DB *gorm.DB
}

func (r *GruposUsuariosRepository) GetByGrupoId(grupoId uint64) ([]app.GruposUsuarios, error) {
	var gruposUsuarios []app.GruposUsuarios
	if err := r.DB.Where("id_grupo = ?", grupoId).Find(&gruposUsuarios).Error; err != nil {
		return nil, err
	}
	return gruposUsuarios, nil
}

func (r *GruposUsuariosRepository) GetByGrupoUsuarioId(usuarioId uint64) (*app.GruposUsuarios, error) {
	var grupoUsuario app.GruposUsuarios
	if err := r.DB.Where("id_usuario = ?", usuarioId).First(&grupoUsuario).Error; err != nil {
		return nil, err
	}
	return &grupoUsuario, nil
}

func (r *GruposUsuariosRepository) Create(grupoUsuario *app.GruposUsuarios) error {
	if err := r.DB.Create(grupoUsuario).Error; err != nil {
		return err
	}
	return nil
}

func (r *GruposUsuariosRepository) CreateTx(grupoUsuario *app.GruposUsuarios, tx *gorm.DB) error {
	if err := tx.Create(grupoUsuario).Error; err != nil {
		return err
	}
	return nil
}

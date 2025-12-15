package grupo

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/grupousario"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type GrupoRepository struct {
	DB             *gorm.DB
	RepoRelaciones *grupousario.GruposUsuariosRepository
}

func (r *GrupoRepository) GetAll() ([]app.Grupos, error) {

	var grupos []app.Grupos

	if err := r.DB.Find(&grupos).Error; err != nil {
		return nil, err
	}

	return grupos, nil
}

func (r *GrupoRepository) GetById(id uint64) (*app.Grupos, error) {

	var grupo app.Grupos

	if err := r.DB.First(&grupo, id).Error; err != nil {
		return nil, err
	}

	return &grupo, nil
}

func (r *GrupoRepository) GetByClave(clave string) (*app.Grupos, int, error) {

	var grupo app.Grupos

	err := r.DB.Model(&app.Grupos{}).Where("clave = ?", clave).First(&grupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No se encontró el registro, esto es esperado cuando creamos un nuevo usuario
			return nil, 404, nil
		}

		// Otro tipo de error en la base de datos
		return nil, 500, fmt.Errorf("error al buscar grupo por clave: %w", err)
	}

	return &grupo, 200, nil
}

func (r *GrupoRepository) GetAllByUsuarioId(usuarioId uint64) ([]app.Grupos, error) {
	var grupos []app.Grupos

	err := r.DB.
		Joins("JOIN grupos_usuarios ON grupos_usuarios.id_grupo = grupos.id").
		Where("grupos_usuarios.id_usuario = ?", usuarioId).
		Find(&grupos).Error

	if err != nil {
		return nil, err
	}

	if len(grupos) == 0 {
		return grupos, nil
	}

	grupoIDs := make([]uint64, len(grupos))
	for i, grupo := range grupos {
		grupoIDs[i] = grupo.Id
	}

	var ultimosMensajes []app.Mensajes

	// Usando DISTINCT ON de PostgreSQL
	err = r.DB.
		Preload("Respuesta").
		Where("id IN (SELECT DISTINCT ON (id_grupo) id FROM mensajes WHERE id_grupo IN (?) ORDER BY id_grupo, fecha DESC)", grupoIDs).
		Find(&ultimosMensajes).Error

	if err != nil {
		return nil, err
	}

	mensajesPorGrupo := make(map[uint64]app.Mensajes)
	for _, mensaje := range ultimosMensajes {
		mensajesPorGrupo[mensaje.GrupoId] = mensaje
	}

	for i := range grupos {
		if mensaje, exists := mensajesPorGrupo[grupos[i].Id]; exists {
			grupos[i].Mensajes = []app.Mensajes{mensaje}
		}
	}

	return grupos, nil
}

func (r *GrupoRepository) GetAllGruposByUsuarioIdToClaves(usuarioId uint64) ([]string, error) {

	var grupos []app.Grupos

	err := r.DB.Model(&app.Grupos{}).
		Joins("JOIN grupos_usuarios ON grupos_usuarios.id_grupo = grupos.id").
		Where("grupos_usuarios.id_usuario = ?", usuarioId).
		Find(&grupos).Error

	if err != nil {
		return nil, err
	}

	var grupoIds []string

	for _, grupo := range grupos {
		grupoIds = append(grupoIds, grupo.Clave)
	}

	return grupoIds, nil
}

func (r *GrupoRepository) GetByName(name string) (*app.Grupos, int, error) {

	var grupo app.Grupos

	err := r.DB.Where("LOWER(nombre) = LOWER(?)", name).First(&grupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No se encontró el registro, esto es esperado cuando creamos un nuevo usuario
			return nil, 404, nil
		}
		// Otro tipo de error en la base de datos
		return nil, 500, fmt.Errorf("error al buscar grupo por nombre: %w", err)
	}

	return &grupo, 200, nil
}

func (r *GrupoRepository) Create(grupo *app.Grupos) error {

	err := r.DB.Transaction(func(tx *gorm.DB) error {

		// Check if the group already exists by name
		existingGrupo, code, err := r.GetByName(grupo.Nombre)
		if err != nil && code != 404 {
			return fmt.Errorf("error checking existing group: %w", err)
		}

		if existingGrupo != nil {
			return fmt.Errorf("el grupo con nombre '%s' ya existe", grupo.Nombre)
		}

		if err := tx.Create(grupo).Error; err != nil {
			return err
		}

		// Create the relationship with the user
		newRel := &app.GruposUsuarios{
			IdGrupo:   grupo.Id,
			IdUsuario: grupo.CreatedById,
		}

		if err := r.RepoRelaciones.CreateTx(newRel, tx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

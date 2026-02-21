package repository

import (
	"chatvis-chat/internal/models"
	"chatvis-chat/internal/domain"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type postgresGrupoRepository struct {
	db *gorm.DB
}

func NewPostgresGrupoRepository(db *gorm.DB) domain.GrupoRepository {
	return &postgresGrupoRepository{db: db}
}

// mapGormToDomain convierte un modelo de GORM a la entidad de dominio
func mapGormToDomain(gormGrupo *models.Grupos) *domain.Grupo {
	if gormGrupo == nil {
		return nil
	}
	domainGrupo := &domain.Grupo{
		Id:          gormGrupo.Id,
		Clave:       gormGrupo.Clave,
		Nombre:      gormGrupo.Nombre,
		Fecha:       gormGrupo.Fecha,
		CreatedById: gormGrupo.CreatedById,
	}

	// Mapear el último mensaje opcionalmente (basado en la lógica anterior)
	if len(gormGrupo.Mensajes) > 0 {
		m := gormGrupo.Mensajes[0]
		domainMensaje := domain.Mensaje{
			Id:         m.Id,
			Contenido:  m.Contenido,
			Fecha:      m.Fecha,
			GrupoId:    m.GrupoId,
			UsuarioId:  m.UsuarioId,
			ResponseId: m.ResponseId,
		}
		// Si hubiese respuesta adjunta, se mapea aquí (dependiendo de preloads)
		if m.Respuesta != nil {
			domainMensaje.Respuesta = &domain.Mensaje{
				Id:        m.Respuesta.Id,
				Contenido: m.Respuesta.Contenido,
			}
		}
		domainGrupo.Mensajes = []domain.Mensaje{domainMensaje}
	}

	return domainGrupo
}

// mapDomainToGorm convierte la entidad de dominio al modelo de GORM
func mapDomainToGorm(domainGrupo *domain.Grupo) *models.Grupos {
	if domainGrupo == nil {
		return nil
	}
	return &models.Grupos{
		Id:          domainGrupo.Id,
		Clave:       domainGrupo.Clave,
		Nombre:      domainGrupo.Nombre,
		Fecha:       domainGrupo.Fecha,
		CreatedById: domainGrupo.CreatedById,
	}
}

func (r *postgresGrupoRepository) GetAll() ([]domain.Grupo, error) {
	var gormGrupos []models.Grupos

	if err := r.db.Find(&gormGrupos).Error; err != nil {
		return nil, err
	}

	var grupos []domain.Grupo
	for _, g := range gormGrupos {
		grupos = append(grupos, *mapGormToDomain(&g))
	}

	return grupos, nil
}

func (r *postgresGrupoRepository) GetById(id uint64) (*domain.Grupo, error) {
	var gormGrupo models.Grupos

	if err := r.db.First(&gormGrupo, id).Error; err != nil {
		return nil, err
	}

	return mapGormToDomain(&gormGrupo), nil
}

func (r *postgresGrupoRepository) GetByClave(clave string) (*domain.Grupo, int, error) {
	var gormGrupo models.Grupos

	err := r.db.Model(&models.Grupos{}).Where("clave = ?", clave).First(&gormGrupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 404, nil
		}
		return nil, 500, fmt.Errorf("error al buscar grupo por clave: %w", err)
	}

	return mapGormToDomain(&gormGrupo), 200, nil
}

func (r *postgresGrupoRepository) GetAllByUsuarioId(usuarioId uint64) ([]domain.Grupo, error) {
	var gormGrupos []models.Grupos

	err := r.db.
		Joins("JOIN grupos_usuarios ON grupos_usuarios.id_grupo = grupos.id").
		Where("grupos_usuarios.id_usuario = ?", usuarioId).
		Find(&gormGrupos).Error

	if err != nil {
		return nil, err
	}

	if len(gormGrupos) == 0 {
		return []domain.Grupo{}, nil
	}

	grupoIDs := make([]uint64, len(gormGrupos))
	for i, g := range gormGrupos {
		grupoIDs[i] = g.Id
	}

	var ultimosMensajes []models.Mensajes

	// Usando DISTINCT ON de PostgreSQL para obtener el último mensaje por grupo
	err = r.db.
		Preload("Respuesta").
		Where("id IN (SELECT DISTINCT ON (id_grupo) id FROM mensajes WHERE id_grupo IN (?) ORDER BY id_grupo, fecha DESC)", grupoIDs).
		Find(&ultimosMensajes).Error

	if err != nil {
		return nil, err
	}

	mensajesPorGrupo := make(map[uint64]models.Mensajes)
	for _, m := range ultimosMensajes {
		mensajesPorGrupo[m.GrupoId] = m
	}

	var grupos []domain.Grupo
	for i := range gormGrupos {
		if mensaje, exists := mensajesPorGrupo[gormGrupos[i].Id]; exists {
			gormGrupos[i].Mensajes = []models.Mensajes{mensaje}
		}
		grupos = append(grupos, *mapGormToDomain(&gormGrupos[i]))
	}

	return grupos, nil
}

func (r *postgresGrupoRepository) GetAllGruposByUsuarioIdToClaves(usuarioId uint64) ([]string, error) {
	var gormGrupos []models.Grupos

	err := r.db.Model(&models.Grupos{}).
		Joins("JOIN grupos_usuarios ON grupos_usuarios.id_grupo = grupos.id").
		Where("grupos_usuarios.id_usuario = ?", usuarioId).
		Find(&gormGrupos).Error

	if err != nil {
		return nil, err
	}

	var grupoIds []string
	for _, grupo := range gormGrupos {
		grupoIds = append(grupoIds, grupo.Clave)
	}

	return grupoIds, nil
}

func (r *postgresGrupoRepository) GetByName(name string) (*domain.Grupo, int, error) {
	var gormGrupo models.Grupos

	err := r.db.Where("LOWER(nombre) = LOWER(?)", name).First(&gormGrupo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 404, nil
		}
		return nil, 500, fmt.Errorf("error al buscar grupo por nombre: %w", err)
	}

	return mapGormToDomain(&gormGrupo), 200, nil
}

func (r *postgresGrupoRepository) Create(grupo *domain.Grupo) error {

	err := r.db.Transaction(func(tx *gorm.DB) error {

		// Controladores anteriores revisaban si existía
		existingGrupo, code, err := r.GetByName(grupo.Nombre)
		if err != nil && code != 404 {
			return fmt.Errorf("error checking existing group: %w", err)
		}

		if existingGrupo != nil {
			return fmt.Errorf("el grupo con nombre '%s' ya existe", grupo.Nombre)
		}

		gormGrupo := mapDomainToGorm(grupo)
		if err := tx.Create(gormGrupo).Error; err != nil {
			return err
		}

		// Inject back the generated ID
		grupo.Id = gormGrupo.Id

		// Create the relationship with the user directly mapping to GroupsUsers table
		newRel := models.GruposUsuarios{
			IdGrupo:   gormGrupo.Id,
			IdUsuario: gormGrupo.CreatedById,
		}

		if err := tx.Create(&newRel).Error; err != nil {
			return err
		}

		return nil
	})

	return err
}

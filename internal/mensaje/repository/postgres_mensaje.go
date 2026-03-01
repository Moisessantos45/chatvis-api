package repository

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/models"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type postgresMensajeRepository struct {
	db *gorm.DB
}

func NewPostgresMensajeRepository(db *gorm.DB) domain.MensajeRepository {
	return &postgresMensajeRepository{db: db}
}

func mapGormToDomainMensaje(gormMsg *models.Mensajes) *domain.Mensaje {
	if gormMsg == nil {
		return nil
	}

	domainMsg := &domain.Mensaje{
		Id:         gormMsg.Id,
		Contenido:  gormMsg.Contenido,
		Fecha:      gormMsg.Fecha,
		GrupoId:    gormMsg.GrupoId,
		UsuarioId:  gormMsg.UsuarioId,
		ResponseId: gormMsg.ResponseId,
	}

	if gormMsg.Usuario.Id != 0 {
		domainMsg.Usuario = &domain.Usuario{
			Id:     gormMsg.Usuario.Id,
			Nombre: gormMsg.Usuario.Nombre,
			Apodo:  gormMsg.Usuario.Apodo,
			Email:  gormMsg.Usuario.Email,
			IsLlm:  gormMsg.Usuario.IsLlm,
		}
	}

	if gormMsg.Respuesta != nil {
		domainMsg.Respuesta = mapGormToDomainMensaje(gormMsg.Respuesta)
	}

	return domainMsg
}

func mapDomainToGormMensaje(domainMsg *domain.Mensaje) *models.Mensajes {
	if domainMsg == nil {
		return nil
	}
	return &models.Mensajes{
		Id:         domainMsg.Id,
		Contenido:  domainMsg.Contenido,
		Fecha:      domainMsg.Fecha,
		GrupoId:    domainMsg.GrupoId,
		UsuarioId:  domainMsg.UsuarioId,
		ResponseId: domainMsg.ResponseId,
	}
}

func (r *postgresMensajeRepository) GetAll() ([]domain.Mensaje, error) {
	var gormMensajes []models.Mensajes

	if err := r.db.Find(&gormMensajes).Error; err != nil {
		return nil, err
	}

	var mensajes []domain.Mensaje
	for _, gm := range gormMensajes {
		mensajes = append(mensajes, *mapGormToDomainMensaje(&gm))
	}

	return mensajes, nil
}

func (r *postgresMensajeRepository) GetById(id uint64) (*domain.Mensaje, error) {
	var gormMensaje models.Mensajes

	if err := r.db.First(&gormMensaje, id).Error; err != nil {
		return nil, err
	}
	return mapGormToDomainMensaje(&gormMensaje), nil
}

func (r *postgresMensajeRepository) GetAllByGrupoId(grupoId uint64) ([]domain.Mensaje, error) {
	var gormMensajes []models.Mensajes

	err := r.db.
		Preload("Usuario", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nombre", "apodo")
		}).
		Where("id_grupo = ?", grupoId).
		Find(&gormMensajes).Error

	if err != nil {
		return nil, err
	}

	if len(gormMensajes) == 0 {
		return []domain.Mensaje{}, nil
	}

	var mensajes []domain.Mensaje
	for _, gm := range gormMensajes {
		log.Println("Mensaje: ", gm)
		mensajes = append(mensajes, *mapGormToDomainMensaje(&gm))
	}

	return mensajes, nil
}

func (r *postgresMensajeRepository) GetAllByGrupoClave(clave string) ([]domain.Mensaje, error) {
	var gormMensajes []models.Mensajes

	var grupo models.Grupos
	err := r.db.Select("id").Where("clave = ?", clave).First(&grupo).Error
	if err != nil {
		return nil, err
	}

	err = r.db.
		Preload("Usuario", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nombre", "apodo")
		}).
		Preload("Respuesta").
		Preload("Respuesta.Usuario", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nombre", "apodo")
		}).
		Where("id_grupo = ?", grupo.Id).
		Find(&gormMensajes).Error

	if err != nil {
		return nil, err
	}

	var mensajes []domain.Mensaje
	for _, gm := range gormMensajes {
		mensajes = append(mensajes, *mapGormToDomainMensaje(&gm))
	}

	return mensajes, nil
}

func (r *postgresMensajeRepository) Create(mensaje *domain.Mensaje) (*domain.Mensaje, error) {
	gormMensaje := mapDomainToGormMensaje(mensaje)
	if err := r.db.Create(gormMensaje).Error; err != nil {
		return nil, err
	}
	mensaje.Id = gormMensaje.Id
	return mensaje, nil
}

func (r *postgresMensajeRepository) Update(id uint64, mensaje *domain.Mensaje) error {
	if id <= 0 {
		return errors.New("el ID del mensaje debe ser mayor que cero")
	}

	if mensaje == nil {
		return errors.New("el mensaje no puede ser nulo")
	}

	var existingGormMensaje models.Mensajes
	if err := r.db.First(&existingGormMensaje, id).Error; err != nil {
		return err
	}

	existingGormMensaje.Contenido = mensaje.Contenido
	existingGormMensaje.ResponseId = mensaje.ResponseId

	if err := r.db.Save(&existingGormMensaje).Error; err != nil {
		return err
	}

	return nil
}

func (r *postgresMensajeRepository) GetNuevosMensajesParaIA(aiID uint64, grupoID uint64) ([]domain.Mensaje, error) {
	var checkpoint models.ModelSyncCheckpoint
	var gormMensajes []models.Mensajes

	r.db.Where("id_usuario = ? AND id_grupo = ?", aiID, grupoID).First(&checkpoint)

	err := r.db.
		Where("id_grupo = ? AND id > ? AND id_usuario != ?",
			grupoID, checkpoint.UltimoMensajeId, aiID).
		Order("id asc").
		Find(&gormMensajes).Error

	var mensajes []domain.Mensaje
	for _, gm := range gormMensajes {
		mensajes = append(mensajes, *mapGormToDomainMensaje(&gm))
	}

	return mensajes, err
}

func (r *postgresMensajeRepository) ActualizarPuntoControl(aiID uint64, grupoID uint64, ultimoID uint64) error {
	checkpoint := models.ModelSyncCheckpoint{
		UsuarioId:       aiID,
		GrupoId:         grupoID,
		UltimoMensajeId: ultimoID,
		UpdatedAt:       time.Now(),
	}
	return r.db.Save(&checkpoint).Error
}

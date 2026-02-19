package mensaje

import (
	"chatvis-chat/internal/app"
	"errors"

	"gorm.io/gorm"
)

type MensajeRepository struct {
	DB *gorm.DB
}

func (r *MensajeRepository) GetAll() ([]app.Mensajes, error) {
	var mensajes []app.Mensajes

	if err := r.DB.Find(&mensajes).Error; err != nil {
		return nil, err
	}

	return mensajes, nil
}

func (r *MensajeRepository) GetById(id uint64) (*app.Mensajes, error) {
	var mensaje app.Mensajes

	if err := r.DB.First(&mensaje, id).Error; err != nil {
		return nil, err
	}
	return &mensaje, nil
}

func (r *MensajeRepository) GetAllByGrupoId(grupoId uint64) ([]app.Mensajes, error) {
	var mensajes []app.Mensajes

	err := r.DB.
		Preload("Usuario", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nombre", "apodo")
		}).
		Where("id_grupo = ?", grupoId).
		Find(&mensajes).Error

	if err != nil {
		return nil, err
	}

	return mensajes, nil
}

func (r *MensajeRepository) GetAllByGrupoClave(clave string) ([]app.Mensajes, error) {
	var mensajes []app.Mensajes

	// Primero obtenemos el ID del grupo basado en la clave
	var grupo app.Grupos
	err := r.DB.Select("id").Where("clave = ?", clave).First(&grupo).Error
	if err != nil {
		return nil, err
	}

	// Luego obtenemos los mensajes usando el ID del grupo
	err = r.DB.
		Preload("Usuario", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nombre", "apodo")
		}).
		Preload("Respuesta").
		Preload("Respuesta.Usuario", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "nombre", "apodo")
		}).
		Where("id_grupo = ?", grupo.Id).
		Find(&mensajes).Error

	if err != nil {
		return nil, err
	}

	return mensajes, nil
}

func (r *MensajeRepository) Create(mensaje *app.Mensajes) (*app.Mensajes, error) {
	if err := r.DB.Create(mensaje).Error; err != nil {
		return nil, err
	}

	return mensaje, nil
}

func (r *MensajeRepository) Update(id uint64, mensaje *app.Mensajes) error {

	if id <= 0 {
		return errors.New("El ID del mensaje debe ser mayor que cero")
	}

	if mensaje == nil {
		return errors.New("El mensaje no puede ser nulo")
	}

	existingMensaje, err := r.GetById(id)
	if err != nil {
		return err
	}

	existingMensaje.Contenido = mensaje.Contenido
	existingMensaje.ResponseId = mensaje.ResponseId

	if err := r.DB.Save(existingMensaje).Error; err != nil {
		return err
	}

	return nil
}

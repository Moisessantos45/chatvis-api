package servicemensaje

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/mensaje"
	"errors"
	"fmt"
	"strings"
	"time"
)

type MensajeService struct {
	repo mensaje.MensajeRepository
}

func NewMensajeService(repo mensaje.MensajeRepository) *MensajeService {
	return &MensajeService{
		repo: repo,
	}
}

func (s *MensajeService) GetAll() ([]app.Mensajes, error) {

	mensajes, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	return mensajes, nil
}

func (s *MensajeService) GetAllByGrupoClave(grupoClave string) ([]app.Mensajes, error) {
	if len(strings.TrimSpace(grupoClave)) == 0 {
		return nil, errors.New("La clave del grupo no puede estar vacía")
	}

	mensajes, err := s.repo.GetAllByGrupoClave(grupoClave)
	if err != nil {
		return nil, err
	}

	return mensajes, nil
}

func (s *MensajeService) GetById(id uint64) (*app.Mensajes, error) {

	if id <= 0 {
		return nil, errors.New("El ID del mensaje debe ser mayor que cero")
	}

	mensaje, err := s.repo.GetById(id)
	if err != nil {
		return nil, err
	}

	return mensaje, nil
}

func (s *MensajeService) GetAllByGrupoId(grupoId uint64) ([]app.Mensajes, error) {
	if grupoId <= 0 {
		return nil, errors.New("El ID del grupo debe ser mayor que cero")
	}

	mensajes, err := s.repo.GetAllByGrupoId(grupoId)
	if err != nil {
		return nil, err
	}

	return mensajes, nil
}

func (s *MensajeService) Create(mensaje *app.Mensajes) (*app.Mensajes, error) {
	if mensaje == nil {
		return nil, errors.New("El mensaje no puede ser nulo")
	}

	if mensaje.GrupoId <= 0 {
		return nil, errors.New("El ID del grupo debe ser mayor que cero")
	}

	if mensaje.UsuarioId <= 0 {
		return nil, errors.New("El ID del usuario debe ser mayor que cero")
	}

	if len(strings.TrimSpace(mensaje.Contenido)) == 0 {
		return nil, errors.New("El contenido del mensaje no puede estar vacío")
	}

	mensaje.Fecha = time.Now()

	mensaje, err := s.repo.Create(mensaje)
	if err != nil {
		return nil, fmt.Errorf("Error al crear el mensaje: %w", err)
	}

	return mensaje, nil
}

func (s *MensajeService) Update(id uint64, mensaje *app.Mensajes) error {
	if id <= 0 {
		return errors.New("El ID del mensaje debe ser mayor que cero")
	}

	if mensaje == nil {
		return errors.New("El mensaje no puede ser nulo")
	}

	if len(strings.TrimSpace(mensaje.Contenido)) == 0 {
		return errors.New("El contenido del mensaje no puede estar vacío")
	}

	// Validar que no haya pasado más de 1 minuto desde la creación
	tiempoTranscurrido := time.Since(mensaje.Fecha)
	if tiempoTranscurrido > time.Minute {
		return fmt.Errorf("No se puede actualizar el mensaje: ha pasado más de 1 minuto desde su creación")
	}

	err := s.repo.Update(id, mensaje)
	if err != nil {
		return fmt.Errorf("Error al actualizar el mensaje: %w", err)
	}

	return nil
}

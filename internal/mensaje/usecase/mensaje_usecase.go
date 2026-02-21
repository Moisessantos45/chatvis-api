package usecase

import (
	"chatvis-chat/internal/domain"
	"errors"
	"fmt"
	"strings"
	"time"
)

type mensajeUseCase struct {
	repo domain.MensajeRepository
}

func NewMensajeUseCase(r domain.MensajeRepository) domain.MensajeUseCase {
	return &mensajeUseCase{repo: r}
}

func (s *mensajeUseCase) GetAll() ([]domain.Mensaje, error) {
	return s.repo.GetAll()
}

func (s *mensajeUseCase) GetAllByGrupoClave(grupoClave string) ([]domain.Mensaje, error) {
	if len(strings.TrimSpace(grupoClave)) == 0 {
		return nil, errors.New("la clave del grupo no puede estar vacía")
	}

	return s.repo.GetAllByGrupoClave(grupoClave)
}

func (s *mensajeUseCase) GetById(id uint64) (*domain.Mensaje, error) {
	if id <= 0 {
		return nil, errors.New("el ID del mensaje debe ser mayor que cero")
	}

	return s.repo.GetById(id)
}

func (s *mensajeUseCase) GetAllByGrupoId(grupoId uint64) ([]domain.Mensaje, error) {
	if grupoId <= 0 {
		return nil, errors.New("el ID del grupo debe ser mayor que cero")
	}

	return s.repo.GetAllByGrupoId(grupoId)
}

func (s *mensajeUseCase) Create(mensaje *domain.Mensaje) (*domain.Mensaje, error) {
	if mensaje == nil {
		return nil, errors.New("el mensaje no puede ser nulo")
	}

	if mensaje.GrupoId <= 0 {
		return nil, errors.New("el ID del grupo debe ser mayor que cero")
	}

	if mensaje.UsuarioId <= 0 {
		return nil, errors.New("el ID del usuario debe ser mayor que cero")
	}

	if len(strings.TrimSpace(mensaje.Contenido)) == 0 {
		return nil, errors.New("el contenido del mensaje no puede estar vacío")
	}

	mensaje.Fecha = time.Now()

	mensajeCreado, err := s.repo.Create(mensaje)
	if err != nil {
		return nil, fmt.Errorf("error al crear el mensaje: %w", err)
	}

	return mensajeCreado, nil
}

func (s *mensajeUseCase) Update(id uint64, mensaje *domain.Mensaje) error {
	if id <= 0 {
		return errors.New("el ID del mensaje debe ser mayor que cero")
	}

	if mensaje == nil {
		return errors.New("el mensaje no puede ser nulo")
	}

	if len(strings.TrimSpace(mensaje.Contenido)) == 0 {
		return errors.New("el contenido del mensaje no puede estar vacío")
	}

	// Validar que no haya pasado más de 1 minuto desde la creación
	// NOTA: en una arquitectura estrica, el GetById lo haríamos aquí para validar el tiempo de creación.
	existingMensaje, err := s.repo.GetById(id)
	if err != nil {
		return err
	}

	tiempoTranscurrido := time.Since(existingMensaje.Fecha)
	if tiempoTranscurrido > time.Minute {
		return fmt.Errorf("no se puede actualizar el mensaje: ha pasado más de 1 minuto desde su creación")
	}

	err = s.repo.Update(id, mensaje)
	if err != nil {
		return fmt.Errorf("error al actualizar el mensaje: %w", err)
	}

	return nil
}

func (s *mensajeUseCase) GetNuevosMensajesParaIA(aiID uint64, grupoID uint64) ([]domain.Mensaje, error) {
	return s.repo.GetNuevosMensajesParaIA(aiID, grupoID)
}

func (s *mensajeUseCase) ActualizarPuntoControl(aiID uint64, grupoID uint64, ultimoID uint64) error {
	return s.repo.ActualizarPuntoControl(aiID, grupoID, ultimoID)
}

package domain

import "time"

// Mensaje representa la entidad de dominio pura de un mensaje.
// Al igual que Usuario, no contiene etiquetas de GORM u otros frameworks.
type Mensaje struct {
	Id         uint64    `json:"id"`
	Contenido  string    `json:"contenido"`
	Fecha      time.Time `json:"fecha"`
	GrupoId    uint64    `json:"grupoId"`
	UsuarioId  uint64    `json:"usuarioId"`
	ResponseId *uint64   `json:"respuestaId,omitempty"`

	// Relaciones opcionales para cuando se hace fetch con Joins
	Respuesta *Mensaje `json:"respuesta,omitempty"`
	Usuario   *Usuario `json:"usuario,omitempty"`
}

// MensajeRepository define los métodos que cualquier implementación de DB debe cumplir
type MensajeRepository interface {
	GetAll() ([]Mensaje, error)
	GetById(id uint64) (*Mensaje, error)
	GetAllByGrupoId(grupoId uint64) ([]Mensaje, error)
	GetAllByGrupoClave(clave string) ([]Mensaje, error)
	Create(mensaje *Mensaje) (*Mensaje, error)
	Update(id uint64, mensaje *Mensaje) error

	// IA Checkpoints
	GetNuevosMensajesParaIA(aiID uint64, grupoID uint64) ([]Mensaje, error)
	ActualizarPuntoControl(aiID uint64, grupoID uint64, ultimoID uint64) error
}

// MensajeUseCase define la lógica de negocio expuesta a los controladores
type MensajeUseCase interface {
	GetAll() ([]Mensaje, error)
	GetById(id uint64) (*Mensaje, error)
	GetAllByGrupoId(grupoId uint64) ([]Mensaje, error)
	GetAllByGrupoClave(clave string) ([]Mensaje, error)
	Create(mensaje *Mensaje) (*Mensaje, error)
	Update(id uint64, mensaje *Mensaje) error

	// IA Checkpoints
	GetNuevosMensajesParaIA(aiID uint64, grupoID uint64) ([]Mensaje, error)
	ActualizarPuntoControl(aiID uint64, grupoID uint64, ultimoID uint64) error
}

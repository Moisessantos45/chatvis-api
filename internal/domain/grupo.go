package domain

import "time"

// Grupo representa la entidad principal para los chats de grupo
type Grupo struct {
	Id          uint64    `json:"id"`
	Clave       string    `json:"clave"`
	Nombre      string    `json:"nombre"`
	Fecha       time.Time `json:"fecha"`
	CreatedById uint64    `json:"createdById"`

	// Relaciones (Opcionales dependiendo del fetch)
	UsuarioCreatedBy *Usuario  `json:"usuarioCreatedBy,omitempty"`
	Usuarios         []Usuario `json:"usuarios,omitempty"`
	Mensajes         []Mensaje `json:"mensajes,omitempty"`
}

// GrupoRepository define los métodos requeridos para acceso a datos del grupo
type GrupoRepository interface {
	GetAll() ([]Grupo, error)
	GetById(id uint64) (*Grupo, error)
	GetByClave(clave string) (*Grupo, int, error)
	GetAllByUsuarioId(usuarioId uint64) ([]Grupo, error)
	GetAllGruposByUsuarioIdToClaves(usuarioId uint64) ([]string, error)
	GetByName(name string) (*Grupo, int, error)
	Create(grupo *Grupo) error
}

// GrupoUseCase define las reglas de negocio para los grupos
type GrupoUseCase interface {
	GetAll() ([]Grupo, error)
	GetById(id uint64) (*Grupo, error)
	GetByClave(clave string) (*Grupo, error)
	GetAllByUsuarioId(usuarioId uint64) ([]Grupo, error)
	GetAllGruposByUsuarioIdToClaves(usuarioId uint64) ([]string, error)
	GetByName(name string) (*Grupo, error)
	Create(grupo *Grupo) error
}

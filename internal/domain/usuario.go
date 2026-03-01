package domain

import "time"

// Usuario representa la entidad de dominio pura de un usuario.
// Nota: No contiene etiquetas de base de datos (GORM) ni de enrutador (JSON).
type Usuario struct {
	Id       uint64    `json:"id"`
	Nombre   string    `json:"nombre"`
	Apodo    string    `json:"apodo"`
	Email    string    `json:"email"`
	Password string    `json:"password,omitempty"`
	Fecha    time.Time `json:"fecha"`
	Token    string    `json:"token,omitempty"`
	IsLlm    bool      `json:"isLlm"`
	IsAdmin  bool      `json:"isAdmin"`
	IsActive bool      `json:"isActive"`
}

// UsuarioRepository define la interfaz que cualquier implementación de base de datos debe cumplir.
type UsuarioRepository interface {
	GetById(id uint64) (*Usuario, error)
	GetByEmail(email string) (*Usuario, int, error)
	GetAllUsuarios() ([]Usuario, error)
	Create(usuario *Usuario) error
	Update(id uint64, usuario Usuario) error
	UpdateToken(id uint64, token string) error
	UpdateIsActive(id uint64, isActive bool) error
}

// UsuarioUseCase define los métodos expuestos a la capa de entrega (HTTP).
type UsuarioUseCase interface {
	GetById(id uint64) (*Usuario, error)
	GetByEmail(email string) (*Usuario, error)
	GetAllUsuarios() ([]Usuario, error)
	Create(usuario *Usuario) error
	Update(id uint64, usuario Usuario) error
	UpdateIsActive(id uint64, isActive bool) error
	ClearToken(id uint64) error
}

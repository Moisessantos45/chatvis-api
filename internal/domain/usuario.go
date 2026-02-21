package domain

import "time"

// Usuario representa la entidad de dominio pura de un usuario.
// Nota: No contiene etiquetas de base de datos (GORM) ni de enrutador (JSON).
type Usuario struct {
	Id       uint64
	Nombre   string
	Apodo    string
	Email    string
	Password string
	Fecha    time.Time
	Token    string
	IsLlm    bool
}

// UsuarioRepository define la interfaz que cualquier implementación de base de datos debe cumplir.
type UsuarioRepository interface {
	GetById(id uint64) (*Usuario, error)
	GetByEmail(email string) (*Usuario, int, error)
	Create(usuario *Usuario) error
	Update(id uint64, usuario Usuario) error
	UpdateToken(id uint64, token string) error
}

// UsuarioUseCase define los métodos expuestos a la capa de entrega (HTTP).
type UsuarioUseCase interface {
	GetById(id uint64) (*Usuario, error)
	GetByEmail(email string) (*Usuario, error)
	Create(usuario *Usuario) error
	Update(id uint64, usuario Usuario) error
}

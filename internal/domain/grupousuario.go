package domain

// GrupoUsuario representa la relación muchos a muchos entre grupos y usuarios
type GrupoUsuario struct {
	IdGrupo   uint64
	IdUsuario uint64
}

// GrupoUsuarioRepository permite operar sobre la relación de grupos y usuarios
type GrupoUsuarioRepository interface {
	GetByGrupoId(grupoId uint64) ([]GrupoUsuario, error)
	GetByUsuarioId(usuarioId uint64) (*GrupoUsuario, error) // Nombrado de forma más lógica que "GetByGrupoUsuarioId" antiguo
	Create(grupoUsuario *GrupoUsuario) error
}

// GrupoUsuarioUseCase define las reglas de negocio para la membresía de grupos
type GrupoUsuarioUseCase interface {
	JoinGroup(userId uint64, claveGrupo string) error
	GetUsersByGroupId(grupoId uint64) ([]GrupoUsuario, error)
	GetByUsuarioId(usuarioId uint64) (*GrupoUsuario, error)
}

package domain

// GrupoUsuario representa la relación muchos a muchos entre grupos y usuarios
type GrupoUsuario struct {
	IdGrupo   uint64 `json:"grupoId"`
	IdUsuario uint64 `json:"usuarioId"`
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
	JoinGroups(usersIds []uint64, groupsIds []uint64) error
	GetUsersByGroupId(grupoId uint64) ([]GrupoUsuario, error)
	GetByUsuarioId(usuarioId uint64) (*GrupoUsuario, error)
}

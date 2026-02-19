package app

import "time"

type GruposUsuarios struct {
	IdGrupo   uint64 `json:"grupoId" gorm:"primaryKey;not null;column:id_grupo"`
	IdUsuario uint64 `json:"usuarioId" gorm:"primaryKey;not null;column:id_usuario"`
}

type Grupos struct {
	Id          uint64    `json:"id" gorm:"primaryKey"`
	Clave       string    `json:"clave" gorm:"type:varchar(100);not null;unique;index"`
	Nombre      string    `json:"nombre" gorm:"type:varchar(100);not null"`
	Fecha       time.Time `json:"fecha" gorm:"type:date;not null"`
	CreatedById uint64    `json:"createdById" gorm:"not null;column:created_by_id"`

	UsuarioCreatedBy Usuarios   `json:"usuario_created_by" gorm:"foreignKey:CreatedById;references:Id"`
	Usuarios         []Usuarios `json:"usuarios" gorm:"many2many:grupos_usuarios;foreignKey:Id;joinForeignKey:IdGrupo;References:Id;JoinReferences:IdUsuario"`
	Mensajes         []Mensajes `json:"mensajes" gorm:"foreignKey:GrupoId;references:Id"`
}

type Usuarios struct {
	Id       uint64    `json:"id" gorm:"primaryKey"`
	Nombre   string    `json:"nombre" gorm:"type:varchar(100);not null"`
	Apodo    string    `json:"apodo" gorm:"type:varchar(100);not null"`
	Email    string    `json:"email" gorm:"type:varchar(100);not null;unique;index"`
	Password string    `json:"password" gorm:"type:varchar(100);not null"`
	Fecha    time.Time `json:"fecha" gorm:"type:date;not null"`
	Token    string    `json:"token" gorm:"type:text"`

	GrupoCreatedBy []Grupos   `json:"gruposCreatedBy" gorm:"foreignKey:CreatedById;references:Id"`
	Grupos         []Grupos   `json:"grupos" gorm:"many2many:grupos_usuarios;foreignKey:Id;joinForeignKey:IdUsuario;References:Id;JoinReferences:IdGrupo"`
	Mensajes       []Mensajes `json:"mensajes" gorm:"foreignKey:UsuarioId;references:Id"`
}

type Mensajes struct {
	Id        uint64    `json:"id" gorm:"primaryKey"`
	Contenido string    `json:"contenido" gorm:"type:text;not null"`
	Fecha     time.Time `json:"fecha" gorm:"type:date;not null"`
	GrupoId   uint64    `json:"grupoId" gorm:"not null;column:id_grupo"`
	UsuarioId uint64    `json:"usuarioId" gorm:"not null;column:id_usuario"`

	ResponseId *uint64   `json:"respuestaId,omitempty" gorm:"column:respuesta_id;default:null"`
	Respuesta  *Mensajes `json:"respuesta,omitempty" gorm:"foreignKey:ResponseId;references:Id"`

	Grupo   Grupos   `json:"grupo" gorm:"foreignKey:GrupoId;references:Id"`
	Usuario Usuarios `json:"usuario" gorm:"foreignKey:UsuarioId;references:Id"`
}

var Models = []any{
	&GruposUsuarios{},
	&Grupos{},
	&Usuarios{},
	&Mensajes{},
}

type UsuarioLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GrupoWithUsuario struct {
	Clave     string `json:"clave"`
	UsuarioId uint64 `json:"usuarioId"`
}

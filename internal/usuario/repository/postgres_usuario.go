package repository

import (
	"chatvis-chat/internal/models"
	"chatvis-chat/internal/domain"
	"errors"

	"gorm.io/gorm"
)

type postgresUsuarioRepository struct {
	db *gorm.DB
}

// NewPostgresUsuarioRepository crea una nueva instancia del repositorio de usuario usando GORM
func NewPostgresUsuarioRepository(db *gorm.DB) domain.UsuarioRepository {
	return &postgresUsuarioRepository{db: db}
}

// mapGormToDomain convierte un modelo de GORM (models.Usuarios) a la entidad de dominio pura (domain.Usuario)
func mapGormToDomain(gormUser *models.Usuarios) *domain.Usuario {
	if gormUser == nil {
		return nil
	}
	return &domain.Usuario{
		Id:       gormUser.Id,
		Nombre:   gormUser.Nombre,
		Apodo:    gormUser.Apodo,
		Email:    gormUser.Email,
		Password: gormUser.Password,
		Fecha:    gormUser.Fecha,
		Token:    gormUser.Token,
		IsLlm:    gormUser.IsLlm,
	}
}

// mapDomainToGorm convierte la entidad de dominio (domain.Usuario) al modelo de GORM (models.Usuarios)
func mapDomainToGorm(domainUser *domain.Usuario) *models.Usuarios {
	if domainUser == nil {
		return nil
	}
	return &models.Usuarios{
		Id:       domainUser.Id,
		Nombre:   domainUser.Nombre,
		Apodo:    domainUser.Apodo,
		Email:    domainUser.Email,
		Password: domainUser.Password,
		Fecha:    domainUser.Fecha,
		Token:    domainUser.Token,
		IsLlm:    domainUser.IsLlm,
	}
}

func (r *postgresUsuarioRepository) GetById(id uint64) (*domain.Usuario, error) {
	var gormUser models.Usuarios
	if err := r.db.First(&gormUser, id).Error; err != nil {
		return nil, err
	}
	return mapGormToDomain(&gormUser), nil
}

func (r *postgresUsuarioRepository) GetByEmail(email string) (*domain.Usuario, int, error) {
	var gormUser models.Usuarios
	if err := r.db.Where("email = ?", email).First(&gormUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 404, err
		}
		return nil, 500, err
	}
	return mapGormToDomain(&gormUser), 200, nil
}

func (r *postgresUsuarioRepository) Create(usuario *domain.Usuario) error {
	gormUser := mapDomainToGorm(usuario)
	if err := r.db.Create(gormUser).Error; err != nil {
		return err
	}
	// Asignamos el ID generado por GORM de vuelta al objeto de dominio
	usuario.Id = gormUser.Id
	return nil
}

func (r *postgresUsuarioRepository) Update(id uint64, usuario domain.Usuario) error {
	var existingGormUser models.Usuarios
	if err := r.db.First(&existingGormUser, id).Error; err != nil {
		return err
	}

	existingGormUser.Nombre = usuario.Nombre
	existingGormUser.Token = usuario.Token

	if err := r.db.Save(&existingGormUser).Error; err != nil {
		return err
	}

	return nil
}

func (r *postgresUsuarioRepository) UpdateToken(id uint64, token string) error {
	var existingGormUser models.Usuarios
	if err := r.db.First(&existingGormUser, id).Error; err != nil {
		return err
	}

	existingGormUser.Token = token

	if err := r.db.Save(&existingGormUser).Error; err != nil {
		return err
	}

	return nil
}

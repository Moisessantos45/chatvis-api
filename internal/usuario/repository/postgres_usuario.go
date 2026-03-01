package repository

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/models"
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
		IsAdmin:  gormUser.IsAdmin,
		IsActive: gormUser.IsActive,
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
		IsAdmin:  domainUser.IsAdmin,
		IsActive: domainUser.IsActive,
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

func (r *postgresUsuarioRepository) GetAllUsuarios() ([]domain.Usuario, error) {
	var gormUsers []models.Usuarios
	if err := r.db.Find(&gormUsers).Error; err != nil {
		return nil, err
	}
	var users []domain.Usuario
	for _, gu := range gormUsers {
		users = append(users, *mapGormToDomain(&gu))
	}
	return users, nil
}

func (r *postgresUsuarioRepository) UpdateIsActive(id uint64, isActive bool) error {
	var existingGormUser models.Usuarios
	if err := r.db.First(&existingGormUser, id).Error; err != nil {
		return err
	}

	existingGormUser.IsActive = isActive

	if err := r.db.Save(&existingGormUser).Error; err != nil {
		return err
	}

	return nil
}

package usuario

import (
	"chatvis-chat/internal/app"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type UsuarioRepository struct {
	DB *gorm.DB
}

func (r *UsuarioRepository) GetById(id uint64) (*app.Usuarios, error) {
	var usuario app.Usuarios

	if err := r.DB.First(&usuario, id).Error; err != nil {
		return nil, err
	}

	return &usuario, nil
}

func (r *UsuarioRepository) GetByEmail(email string) (*app.Usuarios, int, error) {
	var usuario app.Usuarios

	log.Printf("Buscando usuario por email: %s", email)

	err := r.DB.Model(&app.Usuarios{}).Where("email = ?", email).First(&usuario).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No se encontró el registro, esto es esperado cuando creamos un nuevo usuario
			return nil, 404, nil
		}
		// Otro tipo de error en la base de datos
		return nil, 500, fmt.Errorf("error al buscar usuario por email: %w", err)
	}

	return &usuario, 200, nil
}

func (r *UsuarioRepository) Create(usuario *app.Usuarios) error {
	if err := r.DB.Create(usuario).Error; err != nil {
		return err
	}
	return nil
}

func (r *UsuarioRepository) Update(id uint64, usuario app.Usuarios) error {
	existingUsuario, err := r.GetById(id)
	if err != nil {
		return err
	}

	existingUsuario.Nombre = usuario.Nombre
	existingUsuario.Apodo = usuario.Apodo
	existingUsuario.Token = usuario.Token

	if err := r.DB.Save(existingUsuario).Error; err != nil {
		return err
	}

	return nil
}

func (r *UsuarioRepository) UpdatePassword(id uint64, password string) error {
	existingUsuario, err := r.GetById(id)
	if err != nil {
		return err
	}

	existingUsuario.Password = password

	if err := r.DB.Save(existingUsuario).Error; err != nil {
		return err
	}

	return nil
}

func (r *UsuarioRepository) UpdateToken(id uint64, token string) error {

	existingUsuario, err := r.GetById(id)
	if err != nil {
		return err
	}

	existingUsuario.Token = token

	if err := r.DB.Save(existingUsuario).Error; err != nil {
		return err
	}

	return nil
}

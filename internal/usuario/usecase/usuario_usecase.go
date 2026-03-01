package usecase

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
	"errors"
	"fmt"
	"strings"
)

type usuarioUseCase struct {
	repo domain.UsuarioRepository
}

// NewUsuarioUseCase crea una nueva instancia del servicio de usuario inyectando su dependencia
func NewUsuarioUseCase(r domain.UsuarioRepository) domain.UsuarioUseCase {
	return &usuarioUseCase{repo: r}
}

func (uc *usuarioUseCase) GetById(id uint64) (*domain.Usuario, error) {
	if id == 0 {
		return nil, errors.New("id no puede ser 0")
	}

	return uc.repo.GetById(id)
}

func (uc *usuarioUseCase) GetByEmail(email string) (*domain.Usuario, error) {
	if len(strings.TrimSpace(email)) == 0 {
		return nil, errors.New("email no puede ser vacio")
	}

	usuario, code, err := uc.repo.GetByEmail(email)
	if err != nil && code != 404 {
		return nil, fmt.Errorf("error al obtener usuario por email: %w", err)
	}

	if usuario == nil {
		return nil, errors.New("usuario no encontrado")
	}

	return usuario, err
}

func (uc *usuarioUseCase) Create(usuario *domain.Usuario) error {
	if usuario == nil {
		return errors.New("el objeto usuario no puede ser nulo")
	}

	// Validación de nombre
	if len(strings.TrimSpace(usuario.Nombre)) == 0 {
		return errors.New("el nombre no puede estar vacío")
	}

	// Validación de apodo
	if len(strings.TrimSpace(usuario.Apodo)) == 0 {
		return errors.New("el apodo no puede estar vacío")
	}

	// Validación de email
	if len(strings.TrimSpace(usuario.Email)) == 0 {
		return errors.New("el email no puede estar vacío")
	}

	// Validar formato de email
	if !pkg.IsValidEmail(usuario.Email) {
		return errors.New("el formato del email no es válido")
	}

	// Validación de contraseña
	if len(strings.TrimSpace(usuario.Password)) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}

	// Validar fortaleza de la contraseña
	if err := pkg.ValidatePasswordStrength(usuario.Password); err != nil {
		return err
	}

	// Verificar si el usuario ya existe por email
	existingUsuario, code, err := uc.repo.GetByEmail(usuario.Email)
	if err != nil && code != 404 {
		// Solo retornar error si no es el error de "registro no encontrado"
		return fmt.Errorf("error al verificar existencia de usuario: %w", err)
	}

	if existingUsuario != nil {
		return errors.New("ya existe un usuario registrado con ese email")
	}

	// Hash de la contraseña
	hashedPassword, err := pkg.HashPassword(usuario.Password)
	if err != nil {
		return errors.New("error al encriptar la contraseña")
	}

	usuario.Password = hashedPassword

	return uc.repo.Create(usuario)
}

func (uc *usuarioUseCase) Update(id uint64, usuario domain.Usuario) error {
	if id == 0 {
		return errors.New("id no puede ser 0")
	}

	if len(strings.TrimSpace(usuario.Nombre)) == 0 {
		return errors.New("El nombre no puede ser vacio")
	}

	if len(strings.TrimSpace(usuario.Apodo)) == 0 {
		return errors.New("El apodo no puede ser vacio")
	}

	if len(strings.TrimSpace(usuario.Token)) == 0 {
		return errors.New("El token no puede ser vacio")
	}

	existingUsuario, err := uc.repo.GetById(id)
	if err != nil {
		return err
	}

	existingUsuario.Nombre = usuario.Nombre
	existingUsuario.Token = usuario.Token

	return uc.repo.Update(id, *existingUsuario)
}

func (uc *usuarioUseCase) GetAllUsuarios() ([]domain.Usuario, error) {
	return uc.repo.GetAllUsuarios()
}

func (uc *usuarioUseCase) UpdateIsActive(id uint64, isActive bool) error {
	if id == 0 {
		return errors.New("id no puede ser 0")
	}
	return uc.repo.UpdateIsActive(id, isActive)
}

func (uc *usuarioUseCase) ClearToken(id uint64) error {
	if id == 0 {
		return errors.New("id no puede ser 0")
	}
	return uc.repo.UpdateToken(id, "")
}

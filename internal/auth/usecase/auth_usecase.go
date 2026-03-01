package usecase

import (
	"chatvis-chat/internal/domain"
	"chatvis-chat/internal/pkg"
	"errors"
	"fmt"
)

type authUseCase struct {
	repo domain.UsuarioRepository
}

func NewAuthUseCase(repo domain.UsuarioRepository) domain.AuthUseCase {
	return &authUseCase{
		repo: repo,
	}
}

func (s *authUseCase) Authenticate(email, password string) (*domain.Usuario, error) {

	user, code, err := s.repo.GetByEmail(email)
	if err != nil && code != 404 {
		// Solo retornar error si no es el error de "registro no encontrado"
		return nil, fmt.Errorf("error al obtener el usuario: %w", err)
	}

	if user == nil {
		return nil, errors.New("Las credenciales son incorrectas")
	}

	checkPassword := pkg.CheckPasswordHash(password, user.Password)
	if !checkPassword {
		return nil, errors.New("Las credenciales son incorrectas")
	}

	// Generate a new token for the user
	token, err := pkg.GenerateJWT(fmt.Sprint(user.Id), user.Nombre, user.Email, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	user.Token = token
	// actualizar el token en la base de datos
	err = s.repo.UpdateToken(user.Id, token)
	if err != nil {
		return nil, err
	}

	user.Password = "" // Limpiar el campo de la contraseña antes de devolver el usuario

	return user, nil
}

func (s *authUseCase) Logout(email string) error {
	existingUser, code, err := s.repo.GetByEmail(email)
	if err != nil && code != 404 {
		return err
	}

	if existingUser == nil {
		return nil // Si no existe, no hay nada que cerrar
	}

	existingUser.Token = ""
	err = s.repo.UpdateToken(existingUser.Id, "")
	if err != nil {
		return err
	}

	return nil
}

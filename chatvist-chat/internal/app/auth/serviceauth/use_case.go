package serviceauth

import (
	"chatvis-chat/internal/app"
	"chatvis-chat/internal/app/usuario"
	"chatvis-chat/internal/pkg"
	"errors"
	"fmt"
)

type ServiceAuth struct {
	repo usuario.UsuarioRepository
}

func NewServiceAuth(repo usuario.UsuarioRepository) *ServiceAuth {
	return &ServiceAuth{
		repo: repo,
	}
}

func (s *ServiceAuth) Authenticate(email, password string) (*app.Usuarios, error) {

	user, code, err := s.repo.GetByEmail(email)
	if err != nil && code != 404 {
		// Solo retornar error si no es el error de "registro no encontrado"
		return nil, fmt.Errorf("error al obtener el usuario: %w", err)
	}

	checkPassword := pkg.CheckPasswordHash(password, user.Password)
	if !checkPassword {
		return nil, errors.New("Las credenciales son incorrectas")
	}

	// Generate a new token for the user
	token, err := pkg.GenerateJWT(fmt.Sprint(user.Id), user.Nombre, user.Email)
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

func (s *ServiceAuth) Logout(email string) error {

	existingUser, code, err := s.repo.GetByEmail(email)
	if err != nil && code != 404 {
		return err
	}

	existingUser.Token = ""
	err = s.repo.UpdateToken(existingUser.Id, "")
	if err != nil {
		return err
	}

	return nil
}

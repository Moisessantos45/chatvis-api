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
		return nil, fmt.Errorf("error al obtener el usuario: %w", err)
	}

	if user == nil {
		return nil, errors.New("Las credenciales son incorrectas")
	}

	checkPassword := pkg.CheckPasswordHash(password, user.Password)
	if !checkPassword {
		return nil, errors.New("Las credenciales son incorrectas")
	}

	token, err := pkg.GenerateJWT(fmt.Sprint(user.Id), user.Nombre, user.Email, user.IsAdmin)
	if err != nil {
		return nil, err
	}

	user.Token = token
	err = s.repo.UpdateToken(user.Id, token)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return user, nil
}

func (s *authUseCase) Logout(email string) error {
	existingUser, code, err := s.repo.GetByEmail(email)
	if err != nil && code != 404 {
		return err
	}

	if existingUser == nil {
		return nil
	}

	existingUser.Token = ""
	err = s.repo.UpdateToken(existingUser.Id, "")
	if err != nil {
		return err
	}

	return nil
}

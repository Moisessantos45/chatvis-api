package domain

// AuthUseCase define las reglas de negocio para la autenticación
type AuthUseCase interface {
	Authenticate(email, password string) (*Usuario, error)
	Logout(email string) error
}

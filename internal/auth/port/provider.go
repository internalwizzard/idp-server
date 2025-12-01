package auth

type AuthProvider interface {
	Validate(token string) error
}

package security_interface

type PasswordHasher interface {
	Hash(password string) (hash string, err error)
	Verify(user Passwordness, password string) (err error)
}

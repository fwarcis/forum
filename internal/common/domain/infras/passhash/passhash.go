package passhash

import "golang.org/x/crypto/bcrypt"

const MaxPassBytes = 18 // 4 bytes for each UTF-8 char; 72 bytes

type PasswordHasher struct {
	cost int
}

func (h PasswordHasher) Generate(pass string) (string, error) {
	res, err := bcrypt.GenerateFromPassword(
		[]byte(pass), h.cost)
	return string(res), err
}

func (h PasswordHasher) Compare(pass, hash string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(pass))
}

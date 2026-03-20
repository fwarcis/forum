package passhash

import "golang.org/x/crypto/bcrypt"

type passHasher struct {
	cost int
}

func NewPasswordHasher(cost int) passHasher {
	return passHasher{cost}
}

func (h passHasher) Generate(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, h.cost)
}

func (h passHasher) Compare(password, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, password)
}

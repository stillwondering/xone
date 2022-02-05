package bcrypt

import "golang.org/x/crypto/bcrypt"

const cost = 8

func HashFromPassword(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, cost)
}

func PasswordMatchesHash(hashedPassword []byte, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}

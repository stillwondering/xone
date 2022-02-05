package bcrypt_test

import (
	"testing"

	"github.com/stillwondering/xone/bcrypt"
)

func TestPasswortMatchesHash(t *testing.T) {
	password := []byte("SuperSecretPassword")

	hash, err := bcrypt.HashFromPassword(password)
	if err != nil {
		t.Fatalf("HashFromPassword() wantErr = %v, got %v", nil, err)
	}

	if !bcrypt.PasswordMatchesHash(hash, password) {
		t.Errorf("PasswordMatchesHash() want = %v, got %v", true, false)
	}
}

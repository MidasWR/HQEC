package hash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	var passwords []string
	passwords = append(passwords, "123456982312")
	passwords = append(passwords, "123456783121")
	passwords = append(passwords, "876543211321")
	for _, password := range passwords {
		hash, err := HashPassword(password)
		assert.NotEmpty(t, hash)
		assert.Nil(t, err)
		err = VerifyPassword(hash, password)
		assert.Nil(t, err)
	}
}

package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateJWT(t *testing.T) {
	var logins []string
	logins = append(logins, "user1")
	logins = append(logins, "user2")
	logins = append(logins, "user3")
	for _, login := range logins {
		token := GenerateJWT(login)
		n_login, err := ParseJWT(token)
		assert.Nil(t, err)
		assert.Equal(t, login, n_login)
	}

}

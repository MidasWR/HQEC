package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(login string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
	})
	tokenString, _ := token.SignedString([]byte("secret"))
	return tokenString
}
func ParseJWT(tokenString string) (string, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}
	login, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("login not found in token claims")
	}
	return login, nil
}

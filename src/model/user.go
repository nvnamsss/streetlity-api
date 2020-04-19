package model

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const RoleAdmin = 10

func Auth(tokenString string) error {
	fmt.Println("[Auth]", tokenString)
	if tokenString == "" {
		return errors.New("Token is empty")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret-key-0985399536aA"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return nil
	} else {
		fmt.Println(err)
		return err
	}
}

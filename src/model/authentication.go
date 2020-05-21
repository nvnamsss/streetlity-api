package model

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type User struct {
	Id int64
}

const RoleAdmin = 10

func CreateToken(id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"exp": time.Now().Add(time.Minute*10 + time.Second*30).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret-key-0985399536aA"))

	return tokenString, err
}

func Authenticate(tokenString string) error {
	fmt.Println("[Authenticate]", tokenString)
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
		_, ok := claims["id"]
		if !ok {
			return errors.New("Invalid token")
		}
		return nil
	} else {
		log.Println("[Authenticate]", err.Error())
		return err
	}
}

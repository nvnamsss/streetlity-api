package model

import (
	"fmt"
	"log"

	"github.com/dgrijalva/jwt-go"
)

func Auth(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret-key"), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err = claims.Valid()

		if err != nil {
			/** refresh logic **/
			log.Println("expired")
		}

		fmt.Println(claims)
		return true
	} else {
		fmt.Println(err)
		return false
	}
}

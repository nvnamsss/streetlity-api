package router

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func confirm(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	tokenString := query["token"]

	status := model.Auth(tokenString[0])

	w.Write([]byte(strconv.FormatBool(status)))
}

func auth(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Token   string
		Message string
	}

	result.Status = true
	result.Message = "Success"

	query := req.URL.Query()
	id, idReady := query["id"]

	if !idReady {
		result.Status = false
		result.Message = "Id is missing"
	}

	if result.Status {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  id,
			"exp": 60,
		})

		tokenString, err := token.SignedString([]byte("secret-key"))

		result.Token = tokenString
		if err != nil {
			log.Println(err.Error())
		}
	}

	jsonData, jsonErr := json.Marshal(result)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	w.Write(jsonData)
}

func HandleAuth(router *mux.Router) {
	log.Println("[Router]", "Handling auth")
	s := router.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/", auth)
	s.HandleFunc("/confirm", confirm)
}

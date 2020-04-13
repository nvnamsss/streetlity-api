package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"
	"streelity/v1/pipeline"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func confirm(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	tokenString := query["token"]

	status, _ := model.Auth(tokenString[0])

	w.Write([]byte(strconv.FormatBool(status)))
}

func auth(w http.ResponseWriter, req *http.Request) {
	var result struct {
		Status  bool
		Token   string
		Message string
	}

	result.Status = true
	query := req.URL.Query()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() error {
		_, idOk := query["id"]

		if !idOk {
			return errors.New("id param is missing")
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() error {
		_, idErr := strconv.ParseInt(query["id"][0], 10, 64)

		if idErr != nil {
			return errors.New("cannot parse id to int")
		}

		return nil
	})

	validateParamsStage.Next(parseValueStage)
	pipe.First = validateParamsStage

	err := pipe.Run()

	if err != nil {
		result.Status = false
		result.Message = err.Error()
	}

	if result.Status {
		id, _ := strconv.ParseInt(query["id"][0], 10, 64)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  id,
			"exp": time.Now().Add(time.Minute*10 + time.Second*30).Unix(),
		})

		tokenString, err := token.SignedString([]byte("secret-key-0985399536aA"))

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

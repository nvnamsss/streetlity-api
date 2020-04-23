package router

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model"
	"streelity/v1/pipeline"

	"github.com/gorilla/mux"
)

func confirm(w http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	tokenString := query["token"]

	err := model.Auth(tokenString[0])

	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte("success"))
	}
}

func auth(w http.ResponseWriter, req *http.Request) {
	var res struct {
		Response
		Token string
	}

	res.Status = true
	query := req.URL.Query()

	var pipe *pipeline.Pipeline = pipeline.NewPipeline()
	validateParamsStage := pipeline.NewStage(func() error {
		_, idOk := query["id"]

		if !idOk {
			return errors.New("id param is missing")
		}

		return nil
	})

	parseValueStage := pipeline.NewStage(func() (struct{ Id int64 }, error) {
		id, idErr := strconv.ParseInt(query["id"][0], 10, 64)

		if idErr != nil {
			return struct{ Id int64 }{}, errors.New("cannot parse id to int")
		}

		return struct{ Id int64 }{Id: id}, nil
	})

	validateParamsStage.NextStage(parseValueStage)
	pipe.First = validateParamsStage
	res.Error(pipe.Run())

	if res.Status {
		id := pipe.GetInt("Id")[0]

		token, err := model.CreateToken(id)
		res.Token = token
		res.Error(err)
	}

	WriteJson(w, res)
}

func HandleAuth(router *mux.Router) {
	log.Println("[Router]", "Handling auth")
	s := router.PathPrefix("/auth").Subrouter()
	s.HandleFunc("/", auth)
	s.HandleFunc("/confirm", confirm)
}

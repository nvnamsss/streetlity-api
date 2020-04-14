package router

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

type Response struct {
	Status  bool
	Message string
}

//Error validate the data of response by err
func (res *Response) Error(err error) {
	if err != nil {
		res.Status = false
		res.Message = err.Error()
	}
}

func (res *Response) Write(w http.ResponseWriter) {
	jsonData, jsonErr := json.Marshal(res)

	if jsonErr != nil {
		log.Println(jsonErr)
	}

	w.Write(jsonData)
}

func ValidateParams(data url.Values, fields ...string) error {
	for _, field := range fields {
		_, ok := data[field]

		if !ok {
			return errors.New(" param is missing")
		}
	}
	return nil
}

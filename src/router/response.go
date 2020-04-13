package router

import (
	"errors"
	"net/url"
)

type Response struct {
	Status  bool
	Message string
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

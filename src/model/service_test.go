package model_test

import (
	"streelity/v1/model"
	"testing"
)

func TestQueryService(t *testing.T) {
	s := model.Service{Lat: 2, Lon: 2}
	model.Db.Find(&s)

}

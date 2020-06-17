package fuel_test

import (
	"streelity/v1/model"
	"streelity/v1/model/fuel"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestCreateService(t *testing.T) {
	model.ConnectSync()
	gofakeit.Seed(0)
	for loop := 0; loop < 100; loop++ {
		var s fuel.Fuel
		addr := gofakeit.Address()
		s.Address = addr.Address
		s.Lat = float32(addr.Latitude)
		s.Lon = float32(addr.Longitude)
		s.Note = gofakeit.Sentence(30)
		if e := fuel.CreateServices(s); e != nil {
			t.Error(e)
		}
	}

	t.Logf("Completed")
}

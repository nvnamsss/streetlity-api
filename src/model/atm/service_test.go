package atm_test

import (
	"streelity/v1/model"
	"streelity/v1/model/atm"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestCreateService(t *testing.T) {
	model.ConnectSync()
	gofakeit.Seed(0)
	minLat := float32(10.8231 - 0.12)
	maxLat := float32(10.8231 + 0.4)
	minLon := float32(106.6297 - 0.1)
	maxLon := float32(106.6297 + 0.22)
	for loop := 0; loop < 100; loop++ {
		var s atm.Atm
		addr := gofakeit.Address()
		s.Address = addr.Address
		s.Lat = gofakeit.Float32Range(minLat, maxLat)
		s.Lon = gofakeit.Float32Range(minLon, maxLon)
		s.Note = gofakeit.Sentence(30)
		s.BankId = int64(gofakeit.RandomInt([]int{1, 2, 3, 4}))
		if _, e := atm.CreateService(s); e != nil {
			t.Error(e)
		}
	}

	t.Log("Completed")
}

package fuel_test

import (
	"streelity/v1/model/fuel"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestCreateReview(t *testing.T) {
	gofakeit.Seed(0)
	for loop := 0; loop < 100; loop++ {
		service_id := int64(gofakeit.Number(0, 10))
		commenter := int64(gofakeit.Number(0, 100))
		score := gofakeit.Float32Range(0, 5)
		body := gofakeit.Sentence(100)

		fuel.CreateReview(service_id, commenter, score, body)
	}

	t.Logf("Completed")
}

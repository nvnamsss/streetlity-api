package fuel_test

import (
	"log"
	"streelity/v1/model"
	"streelity/v1/model/fuel"
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestCreateReview(t *testing.T) {
	model.ConnectSync()
	gofakeit.Seed(0)
	for loop := 0; loop < 100; loop++ {
		service_id := int64(gofakeit.Number(0, 10))
		commenter := int64(gofakeit.Number(0, 100))
		score := gofakeit.Float32Range(0, 5)
		body := gofakeit.Sentence(100)

		if e := fuel.CreateReview(service_id, commenter, score, body); e != nil {
			t.Error(e)
		}
	}

	t.Logf("Completed")
}

func TestReviewAverageScore(t *testing.T) {
	model.ConnectSync()

	average := fuel.ReviewAverageScore(6)

	log.Println(average)

	t.Error("Completed")
}

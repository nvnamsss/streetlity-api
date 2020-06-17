package ratm

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model/atm"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func ReviewById(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Review atm.Review
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := pipeline.NewStage(func() (str struct {
		ReviewId int64
	}, e error) {
		form := req.PostForm
		review_ids, ok := form["review_id"]
		if !ok {
			return str, errors.New("review_id param is missing")
		}

		review_id, e := strconv.ParseInt(review_ids[0], 10, 64)
		if e != nil {
			return str, errors.New("review_id param cannot parse to int64")
		}

		str.ReviewId = review_id
		return
	})

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		review_id := p.GetIntFirstOrDefault("ReviewId")
		if review, e := atm.ReviewById(review_id); e != nil {
			res.Error(e)
		} else {
			res.Review = review
		}
	}

	sres.WriteJson(w, res)
}
func UpdateReview(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Review atm.Review
	}
	res.Status = true

	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.UpdateReviewValidateStage(req)

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		review_id := p.GetIntFirstOrDefault("ReviewId")
		new_body := p.GetStringFirstOrDefault("NewBody")

		if review, e := atm.ReviewById(review_id); e != nil {
			res.Error(e)
		} else {
			review.Body = new_body
			res.Error(review.Save())
			res.Review = review
		}
	}

	sres.WriteJson(w, res)
}

func ReviewByServiceId(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Reviews []atm.Review
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.QueryReviewByOrderValidate(req)

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		service_id := p.GetIntFirstOrDefault("ServiceId")
		order := p.GetIntFirstOrDefault("Order")
		limit := p.GetIntFirstOrDefault("Limit")
		if reviews, e := atm.ReviewByService(service_id, order, limit); e != nil {
			res.Error(e)
		} else {
			res.Reviews = reviews
		}
	}

	sres.WriteJson(w, res)
}

func CreateReview(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Review atm.Review
	}
	res.Status = true
	req.ParseForm()
	p := pipeline.NewPipeline()
	stage := stages.ReviewValidateStage(req)

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		service_id := p.GetIntFirstOrDefault("ServiceId")
		reviewer := p.GetStringFirstOrDefault("Reviewer")
		score := p.GetFloatFirstOrDefault("Score")
		body := p.GetStringFirstOrDefault("Body")
		res.Error(atm.CreateReview(service_id, reviewer, float32(score), body))
	}

	sres.WriteJson(w, res)
}

func ReviewAverageScore(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Value float64
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := stages.ServiceIdValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		service_id := p.GetIntFirstOrDefault("ServiceId")
		res.Value = atm.ReviewAverageScore(service_id)
	}

	sres.WriteJson(w, res)
}

func DeleteReview(w http.ResponseWriter, req *http.Request) {
	var res sres.Response = sres.Response{Status: true}
	p := pipeline.NewPipeline()
	stage := stages.ReviewIdValidate(req)
	p.First = stage
	res.Error(p.Run())

	if res.Status {
		review_id := p.GetIntFirstOrDefault("ReviewId")
		if e := atm.DeleteReview(review_id); e != nil {
			res.Error(e)
		}
	}
	sres.WriteJson(w, res)
}

func HandleReview(router *mux.Router) {
	log.Println("[Router]", "Handling review atm")
	s := router.PathPrefix("/review").Subrouter()

	s.HandleFunc("/", ReviewById).Methods("GET")
	s.HandleFunc("/", UpdateReview).Methods("POST")
	s.HandleFunc("/", DeleteReview).Methods("DELETE")
	s.HandleFunc("/create", CreateReview).Methods("POST")
	s.HandleFunc("/query", ReviewByServiceId).Methods("GET")
	s.HandleFunc("/score", ReviewAverageScore).Methods("GET")
}

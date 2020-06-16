package rtoilet

import (
	"net/http"
	"streelity/v1/model/toilet"
	"streelity/v1/sres"
	"streelity/v1/stages"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func ReviewById(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Review toilet.Review
	}
	res.Status = true
	p := pipeline.NewPipeline()
	stage := stages.ReviewIdValidate(req)

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		review_id := p.GetIntFirstOrDefault("ReviewId")
		if review, e := toilet.ReviewById(review_id); e != nil {
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
		Review toilet.Review
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

		if review, e := toilet.ReviewById(review_id); e != nil {
			res.Error(e)
		} else {
			review.Body = new_body
			review.Save()
		}
	}

	sres.WriteJson(w, res)
}

func ReviewByServiceId(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Reviews []toilet.Review
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.QueryReviewByOrderValidate(req)

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		service_id := p.GetIntFirstOrDefault("ServiceId")
		order := p.GetIntFirstOrDefault("Order")
		if reviews, e := toilet.ReviewByService(service_id, order, 5); e != nil {
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
		Review toilet.Review
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
		res.Error(toilet.CreateReview(service_id, reviewer, float32(score), body))
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
		res.Value = toilet.ReviewAverageScore(service_id)
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
		if e := toilet.DeleteReview(review_id); e != nil {
			res.Error(e)
		}
	}
	sres.WriteJson(w, res)
}

func HandleReview(router *mux.Router) {
	s := router.PathPrefix("/review").Subrouter()

	s.HandleFunc("/", ReviewById).Methods("GET")
	s.HandleFunc("/", UpdateReview).Methods("POST")
	s.HandleFunc("/create", CreateReview).Methods("POST")
	s.HandleFunc("/query", ReviewByServiceId).Methods("GET")
	s.HandleFunc("/score", ReviewAverageScore).Methods("GET")
}

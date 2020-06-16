package rmaintenance

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"streelity/v1/model/atm"
	"streelity/v1/model/maintenance"
	"streelity/v1/router/sres"
	"streelity/v1/router/stages"

	"github.com/gorilla/mux"
	"github.com/nvnamsss/goinf/pipeline"
)

func ReviewById(w http.ResponseWriter, req *http.Request) {
	var res struct {
		sres.Response
		Review maintenance.Review
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
		if review, e := maintenance.ReviewById(review_id); e != nil {
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
		Review maintenance.Review
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

		if review, e := maintenance.ReviewById(review_id); e != nil {
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
		Reviews []maintenance.Review
	}
	res.Status = true

	p := pipeline.NewPipeline()
	stage := stages.QueryReviewByOrderValidate(req)

	p.First = stage
	res.Error(p.Run())

	if res.Status {
		service_id := p.GetIntFirstOrDefault("ServiceId")
		order := p.GetIntFirstOrDefault("Order")
		if reviews, e := maintenance.ReviewByService(service_id, order, 5); e != nil {
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
		Review maintenance.Review
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
		res.Value = maintenance.ReviewAverageScore(service_id)
	}

	sres.WriteJson(w, res)
}

func Handle(router *mux.Router) {
	log.Println("[Router]", "Handling review maintenance")
	s := router.PathPrefix("/review").Subrouter()

	s.HandleFunc("/", ReviewById).Methods("GET")
	s.HandleFunc("/", UpdateReview).Methods("POST")
	s.HandleFunc("/create", CreateReview).Methods("POST")
	s.HandleFunc("/query", ReviewByServiceId).Methods("GET")
	s.HandleFunc("/score", ReviewAverageScore).Methods("GET")
}

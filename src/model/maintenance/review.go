package maintenance

import (
	"log"
	"math"
	"streelity/v1/model"

	"github.com/jinzhu/gorm"
)

type Review struct {
	model.Review
	db *gorm.DB
}

const ReviewTableName = "maintenace_review"

func (Review) TableName() string {
	return ReviewTableName
}

func CreateReview(service_id int64, reviewer string, score float32, body string) (e error) {
	var review Review = Review{}
	review.ServiceId = service_id
	review.Reviewer = reviewer
	review.Score = score
	review.Body = body

	if e := model.Db.Create(&review).Error; e != nil {
		log.Println("[Database]", "create maintenance review", e.Error())
	}

	return
}

func DeleteReview(review_id int64) (e error) {
	var review Review
	review.Id = review_id
	if e := model.Db.Delete(&review).Error; e != nil {
		log.Println("[Database]", "delete maintenance review", e.Error())
	}

	return
}

func ReviewByService(service_id, order int64, limit int) (reviews []Review, e error) {
	if limit < 0 {
		limit = math.MaxInt64
	}

	if e := model.Db.Where("service_id=?", service_id).Offset(order).Limit(limit).Find(&reviews).Error; e != nil {
		log.Println("[Database]", "get maintenance reviews", e.Error())
	}

	for _, review := range reviews {
		review.db = model.Db
	}
	return
}

func ReviewById(review_id int64) (review Review, e error) {
	e = model.GetById(ReviewTableName, review_id, &review)
	review.db = model.Db
	// if e := model.Db.Where("id=?", review_id).Find(&review).Error; e != nil {
	// 	log.Println("[Database]", "maintenance review by id", e.Error())
	// }

	return
}

func ReviewAverageScore(service_id int64) (average float64) {
	if e := model.Db.Table(ReviewTableName).Select("avg(score)").Where("service_id=?", service_id).Row().Scan(&average); e != nil {
		log.Println("[Database]", "maintenance review average score", e.Error())
	}

	return
}

func (r Review) Save() (e error) {
	if e := r.db.Save(&r).Error; e != nil {
		log.Println("[Database]", "save maintenance review", e.Error())
	}

	return
}

func (r Review) Delete() (e error) {
	return DeleteReview(r.Id)
}

package toilet

import (
	"log"
	"math"
	"streelity/v1/model"

	"github.com/jinzhu/gorm"
)

type Review struct {
	Id        int64
	ServiceId int64   `gorm:"column:service_id"`
	Reviewer  int64   `gorm:"column:commenter"`
	Score     float32 `gorm:"column:score"`
	Body      string  `gorm:"column:body"`
	db        *gorm.DB
}

const ReviewTableName = "toilet_review"

func (Review) TableName() string {
	return ReviewTableName
}

func CreateReview(service_id int64, commenter int64, score float32, body string) (e error) {
	var review Review = Review{}
	review.ServiceId = service_id
	review.Reviewer = commenter
	review.Score = score
	review.Body = body

	if e := model.Db.Create(&review).Error; e != nil {
		log.Println("[Database]", "create toilet review", e.Error())
	}

	return
}

func DeleteReview(review_id int64) (e error) {
	var review Review = Review{Id: review_id}
	if e := model.Db.Delete(&review).Error; e != nil {
		log.Println("[Database]", "delete toilet review", e.Error())
	}

	return
}

func ReviewByService(service_id, order int64, limit int) (reviews []Review, e error) {
	if limit < 0 {
		limit = math.MaxInt64
	}

	if e := model.Db.Where("service_id=?", service_id).Offset(order).Limit(limit).Find(&reviews).Error; e != nil {
		log.Println("[Database]", "get toilet reviews", e.Error())
	}

	for _, review := range reviews {
		review.db = model.Db
	}
	return
}

func ReviewById(review_id int64) (review Review, e error) {
	if e := model.Db.Where("id=?", review_id).Find(&review).Error; e != nil {
		log.Println("[Database]", "toilet review by id", e.Error())
	}

	return
}

func ReviewAverageScore(service_id int64) (average float64) {
	if e := model.Db.Table(ReviewTableName).Select("avg(score)").Where("service_id=?", service_id).Row().Scan(&average); e != nil {
		log.Println("[Database]", "toilet review average score", e.Error())
	}

	return
}

func (r Review) Save() (e error) {
	if e := r.db.Save(&r).Error; e != nil {
		log.Println("[Database]", "save toilet review", e.Error())
	}

	return
}

func (r Review) Delete() (e error) {
	return DeleteReview(r.Id)
}

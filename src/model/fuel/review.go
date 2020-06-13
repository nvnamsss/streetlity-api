package fuel

import (
	"log"
	"streelity/v1/model"

	"github.com/jinzhu/gorm"
)

type Review struct {
	Id        int64
	ServiceId int64   `gorm:"column:service_id"`
	Commenter int64   `gorm:"column:commenter"`
	Score     float32 `gorm:"column:score"`
	Body      string  `gorm:"column:body"`
	db        *gorm.DB
}

func (Review) TableName() string {
	return "fuel_review"
}

func CreateReview(service_id int64, commenter int64, score float32, body string) (e error) {
	var review Review = Review{}
	review.ServiceId = service_id
	review.Commenter = commenter
	review.Score = score
	review.Body = body

	if e := model.Db.Create(&review).Error; e != nil {
		log.Println("[Database]", "create fuel review", e.Error())
	}

	return
}

func DeleteReview(review_id int64) (e error) {
	var review Review = Review{Id: review_id}
	if e := model.Db.Delete(&review).Error; e != nil {
		log.Println("[Database]", "delete fuel review", e.Error())
	}

	return
}

func ReviewByService(service_id, order int64, limit int) (reviews []Review, e error) {
	if e := model.Db.Where("service_id=?", service_id).Offset(order).Limit(limit).Find(&reviews).Error; e != nil {
		log.Println("[Database]", "get fuel reviews", e.Error())
	}

	for _, review := range reviews {
		review.db = model.Db
	}
	return
}

func ReviewById(review_id int64) (review Review, e error) {
	if e := model.Db.Where("id=?", review_id).Find(&review).Error; e != nil {
		log.Println("[Database]", "fuel review by id", e.Error())
	}

	return
}

func (r Review) Save() (e error) {
	if e := r.db.Save(&r).Error; e != nil {
		log.Println("[Database]", "save fuel review", e.Error())
	}

	return
}

func (r Review) Delete() (e error) {
	return DeleteReview(r.Id)
}

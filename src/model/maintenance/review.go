package maintenance

import (
	"log"
	"streelity/v1/model"
)

type Review struct {
	Id        int64
	ServiceId int64   `gorm:"column:service_id"`
	Commenter int64   `gorm:"column:commenter"`
	Score     float32 `gorm:"column:score"`
	Body      string  `gorm:"column:body"`
}

func (Review) TableName() string {
	return "maintenance_review"
}

func Create(service_id int64, commenter int64, score float32, body string) (e error) {
	var review Review = Review{}
	review.ServiceId = service_id
	review.Commenter = commenter
	review.Score = score
	review.Body = body

	if e := model.Db.Create(&review).Error; e != nil {
		log.Println("[Database]", "create maintenance review", e.Error())
	}

	return
}

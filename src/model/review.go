package model

type Review struct {
	Id        int64
	ServiceId int64   `gorm:"column:service_id"`
	Reviewer  string   `gorm:"column:reviewer"`
	Score     float32 `gorm:"column:score"`
	Body      string  `gorm:"column:body"`
}
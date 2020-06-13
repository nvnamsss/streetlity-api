package toilet

type Review struct {
	Id        int64
	ServiceId int64   `gorm:"column:service_id"`
	Commenter int64   `gorm:"column:commenter"`
	Score     float32 `gorm:"column:score"`
	Body      string  `gorm:"column:body"`
}

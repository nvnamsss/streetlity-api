package model

import (
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type ToiletUcf struct {
	Id        int64
	Lat       float32 `gorm:"column:lat"`
	Lon       float32 `gorm:"column:lon"`
	Confident int     `gorm:"column:confident"`
}

//TableName determine the table name in database which is using for gorm
func (ToiletUcf) TableName() string {
	return "ToiletUcf"
}

func (s ToiletUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtms query all the atm serivces
func AllToiletUcfs() []ToiletUcf {
	var services []ToiletUcf
	Db.Find(&services)

	return services
}

//ToiletUcfById query the unconfirmed toilet service by specific id
func ToiletUcfById(id int64) (service ToiletUcf, e error) {
	if e := Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//AddToiletUcf add new ToiletUcf service to the database
//
//return error if there is something wrong when doing transaction
func AddToiletUcf(s ToiletUcf) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//AfterSave automatically run everytime the update transaction is done
func (s *ToiletUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident == 5 {
		var f Toilet = Toilet{Lat: s.Lat, Lon: s.Lon}
		AddToilet(f)
		scope.DB().Delete(s)
		log.Println("Confident is enough")
	}

	return
}

package model

import (
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type MaintainerUcf struct {
	Id        int64
	Lat       float32 `gorm:"column:lat"`
	Lon       float32 `gorm:"column:lon"`
	Confident int     `gorm:"column:confident"`
}

func (MaintainerUcf) TableName() string {
	return "maintainer_ucf"
}

func (s MaintainerUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllMaintainerUcfs query all maintainer services
func AllMaintainerUcfs() []MaintainerUcf {
	var services []MaintainerUcf
	Db.Find(&services)

	return services
}

//AddMaintainerUcf add new unconfirmed maintainer service to the database
//
//return error if there is something wrong when doing transaction
func AddMaintainerUcf(s MaintainerUcf) (e error) {
	if e = Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//MaintainerUcfById query the unconfirmed maintainer service by specific id
func MaintainerUcfById(id int64) (service MaintainerUcf, e error) {
	if e := Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

func (s *MaintainerUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident == 5 {
		var m Maintainer = Maintainer{Lat: s.Lat, Lon: s.Lon}
		AddMaintainer(m)
		scope.DB().Delete(s)
		log.Println("Confident is enough")
	}

	return
}

package model

import (
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type AtmUcf struct {
	ServiceUcf
	BankId int64 `gorm:"column:bank_id"`
}

//TableName determine the table name in database which is using for gorm
func (AtmUcf) TableName() string {
	return "atm_ucf"
}

//Location determine the location of service as r2.Point
func (s AtmUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtmUcfs query all the AtmUcf serivces
func AllAtmUcfs() []AtmUcf {
	var services []AtmUcf
	Db.Find(&services)

	return services
}

//AtmUcfById query the AtmUcf service by specific id
func AtmUcfById(id int64) (service AtmUcf, e error) {
	if e = Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e)
	}

	return
}

//UpvoteAtmUcf upvote the unconfirmed atm by specific id
func UpvoteAtmUcf(id int64) error {
	s, e := AtmUcfById(id)

	if e != nil {
		return e
	}

	s.Confident += 1
	Db.Save(&s)

	return nil
}

//AddAtmUcf add new AtmUcf service to the database
//
//return error if there is something wrong when doing transaction
func AddAtmUcf(s AtmUcf) (e error) {
	if e = Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

func (s *AtmUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident == 5 {
		var a Atm = Atm{Service: Service{Lat: s.Lat, Lon: s.Lon, Address: s.Address}}
		AddAtm(a)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Atm]", "Confident is enough. Added", a)
	}

	return
}

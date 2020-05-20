package model

import (
	"errors"
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type AtmUcf struct {
	ServiceUcf
	BankId int64 `gorm:"column:bank_id"`
}

var confident int = 1

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
	return upvoteAtmUcf(id, 1)
}

func UpvoteAtmUcfImmediately(id int64) error {
	return upvoteAtmUcf(id, confident)
}

func upvoteAtmUcf(id int64, value int) (e error) {
	s, e := AtmUcfById(id)

	if e != nil {
		return e
	}

	s.Confident += value
	if e := Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "upvote unconfirmed atm", id, ":", e.Error())
	}

	return
}

//AddAtmUcf add new AtmUcf service to the database
//
//return error if there is something wrong when doing transaction
func AddAtmUcf(s AtmUcf) (e error) {
	var existed AtmUcf
	if e = Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&existed).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e = Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	//Temporal
	UpvoteAtmUcf(s.Id)
	return
}

func (s *AtmUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident >= confident {
		var a Atm = Atm{Service: s.GetService(), BankId: s.BankId}
		AddAtm(a)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Atm]", "Confident is enough. Added", a)
	}

	return
}

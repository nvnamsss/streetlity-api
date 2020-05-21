package model

import (
	"errors"
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type ToiletUcf struct {
	ServiceUcf
}

//TableName determine the table name in database which is using for gorm
func (ToiletUcf) TableName() string {
	return "toilet_ucf"
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

//UpvoteToiletUcf upvote the unconfirmed toilet by specific id
func UpvoteToiletUcf(id int64) error {
	return upvoteToiletUcf(id, 1)
}

func UpvoteToiletUcfImmediately(id int64) error {
	return upvoteToiletUcf(id, confident)
}

func upvoteToiletUcf(id int64, value int) (e error) {
	s, e := ToiletUcfById(id)

	if e != nil {
		return
	}

	s.Confident += value
	if e := Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "upvote unconfirmed toilet", id, ":", e.Error())
	}

	return
}

func queryToiletUcf(s ToiletUcf) (service ToiletUcf, e error) {
	service = s

	if e := Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query unconfirmed atm", e.Error())
	}

	return
}

func ToiletUcfByService(s ServiceUcf) (service ToiletUcf, e error) {
	service.ServiceUcf = s
	return queryToiletUcf(service)
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
func AddToiletUcf(s ToiletUcf) (e error) {
	var existed ToiletUcf
	if e = Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&existed).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e = Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	//Temporal
	UpvoteToiletUcf(s.Id)
	return
}

//AfterSave automatically run everytime the update transaction is done
func (s *ToiletUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident >= confident {
		var t Toilet = Toilet{Service: s.GetService()}
		AddToilet(t)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Toilet]", "Confident is enough. Added", t)
	}

	return
}

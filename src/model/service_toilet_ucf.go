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

//UpvoteToiletUcf upvote the unconfirmed toilet by specific id
func UpvoteToiletUcf(id int64) error {
	s, e := ToiletUcfById(id)

	if e != nil {
		return e
	}

	s.Confident += 1
	Db.Save(&s)

	return nil
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

	return
}

//AfterSave automatically run everytime the update transaction is done
func (s *ToiletUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident == confident {
		var t Toilet = Toilet{Service: Service{Lat: s.Lat, Lon: s.Lon, Address: s.Address}}
		AddToilet(t)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Toilet]", "Confident is enough. Added", t)
	}

	return
}

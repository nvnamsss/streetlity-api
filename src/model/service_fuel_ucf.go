package model

import (
	"errors"
	"log"

	"github.com/jinzhu/gorm"
)

//FuelUcf representation the Fuel service which is not confirmed
type FuelUcf struct {
	ServiceUcf
}

func (FuelUcf) TableName() string {
	return "fuel_ucf"
}

//AllFuelsUcf query all unconfirmed fuel services
func AllFuelsUcf() []FuelUcf {
	var services []FuelUcf
	Db.Find(&services)

	return services
}

//AddFuelUcf add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddFuelUcf(s FuelUcf) (e error) {
	var existed FuelUcf
	if e = Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&existed).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e = Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	//Temporal
	UpvoteFuelUcf(s.Id)
	return
}

func queryFuelUcf(s FuelUcf) (service FuelUcf, e error) {
	service = s

	if e := Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query unconfirmed fuel", e.Error())
	}

	return
}

func FuelUcfByService(s ServiceUcf) (service FuelUcf, e error) {
	service.ServiceUcf = s
	return queryFuelUcf(service)
}

//FuelUcfById query the fuel service by specific id
func FuelUcfById(id int64) (service FuelUcf, e error) {
	if e = Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//UpvoteFuelUcf upvote the unconfirmed fuel by specific id
func UpvoteFuelUcf(id int64) error {
	return upvoteFuelUcf(id, 1)
}

func UpvoteFuelUcfImmediately(id int64) error {
	return upvoteFuelUcf(id, confident)
}

func upvoteFuelUcf(id int64, value int) (e error) {
	s, e := FuelUcfById(id)

	if e != nil {
		return
	}

	s.Confident += value
	if e := Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "upvote unconfirmed fuel", id, ":", e.Error())
	}

	return
}

func (s *FuelUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident >= confident {
		var f Fuel = Fuel{Service: s.GetService()}
		AddFuel(f)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Fuel]", "Confident is enough. Added", f)
	}

	return
}

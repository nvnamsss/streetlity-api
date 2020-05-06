package model

import (
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
func AddFuelUcf(s FuelUcf) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
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
	s, e := FuelUcfById(id)

	if e != nil {
		return e
	}

	s.Confident += 1
	Db.Save(&s)

	return nil
}

func (s *FuelUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident == 5 {
		var f Fuel = Fuel{Lat: s.Lat, Lon: s.Lon}
		AddFuel(f)
		scope.DB().Delete(s)
		log.Println("Confident is enough")
	}

	return
}

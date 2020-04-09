package model

import (
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type Atm struct {
	gorm.Model
	Location r2.Point
}

func (Atm) TableName() string {
	return "atm"
}

func AllAtms() []Fuel {
	var services []Fuel
	Db.Find(&services)

	return services
}

func AtmById(id int64) Atm {
	var service Atm
	Db.Find(&service, id)

	return service
}

func AllAtmsInRange(circle spatial.Circle) []Fuel {
	var services []Fuel
	Db.Find(&services)

	return services
}

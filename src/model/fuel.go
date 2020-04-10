package model

import (
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
)

type Fuel struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
	// Location r2.Point
}

//Determine table name
func (Fuel) TableName() string {
	return "fuel"
}

func (s Fuel) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

func AllFuels() []Fuel {
	var services []Fuel
	Db.Find(&services)

	return services
}

func FuelById(id int64) Fuel {
	var service Fuel
	Db.Find(&service, id)

	return service
}

func FuelsInRange(circle spatial.Circle) []Fuel {
	return nil
}

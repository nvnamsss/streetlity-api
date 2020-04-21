package model

import (
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

//AllFuels query all fuel services
func AllFuels() []Fuel {
	var services []Fuel
	Db.Find(&services)

	return services
}

//AddFuel add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddFuel(s Fuel) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//FuelById query the fuel service which specific id
func FuelById(id int64) Fuel {
	var service Fuel
	Db.Find(&service, id)

	return service
}

//FuelsInRange query the fuel services which is in the radius of a location
func FuelsInRange(p r2.Point, max_range float64) []Fuel {
	var result []Fuel = []Fuel{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isFuel := item.(Fuel)
			if isFuel && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

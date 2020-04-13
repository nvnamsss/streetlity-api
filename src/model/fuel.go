package model

import (
	"log"

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

func AddFuel(f Fuel) error {
	Db.Create(&f)

	var last Fuel
	Db.Last(&last)

	if last.Lat == f.Lat && last.Lon == f.Lon {
		log.Println("Create new fuel is succeed")
	}

	return nil
}

func FuelById(id int64) Fuel {
	var service Fuel
	Db.Find(&service, id)

	return service
}

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

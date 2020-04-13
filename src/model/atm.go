package model

import (
	"log"

	"github.com/golang/geo/r2"
)

type Atm struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
}

func (Atm) TableName() string {
	return "atm"
}

func (s Atm) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

func AllAtms() []Atm {
	var services []Atm
	Db.Find(&services)

	return services
}

func AtmById(id int64) Atm {
	var service Atm
	Db.Find(&service, id)

	return service
}

func AddAtm(s Atm) error {
	Db.Create(&s)

	var last Atm
	Db.Last(&last)

	if last.Lat == s.Lat && last.Lon == s.Lon {
		log.Println("Create new atm is succeed")
	}

	return nil
}

func AllAtmsInRange(p r2.Point, max_range float64) []Atm {
	var result []Atm = []Atm{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isFuel := item.(Atm)
			if isFuel && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

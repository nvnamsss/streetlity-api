package model

import (
	"github.com/golang/geo/r2"
)

type Toilet struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
}

//TableName determine the table name in database which is using for gorm
func (Toilet) TableName() string {
	return "toilet"
}

func (s Toilet) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtms query all the atm serivces
func AllToilets() []Toilet {
	var services []Toilet
	Db.Find(&services)

	return services
}

//AddToilet add new toilet service to the database
//
//return error if there is something wrong when doing transaction
func AddToilet(s Toilet) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//ToiletsInRange query the toilet services which is in the radius of a location
func ToiletsInRange(p r2.Point, max_range float64) []Toilet {
	var result []Toilet = []Toilet{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			s, isToilet := item.(Toilet)

			if isToilet {
				location := item.Location()
				d := distance(location, p)
				if d < max_range {
					result = append(result, s)
				}
			}
		}
	}
	return result
}

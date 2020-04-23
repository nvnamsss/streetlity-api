package model

import (
	"github.com/golang/geo/r2"
)

type ToiletUcf struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
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

//AddToiletUcf add new ToiletUcf service to the database
//
//return error if there is something wrong when doing transaction
func AddToiletUcf(s ToiletUcf) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

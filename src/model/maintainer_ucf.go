package model

import "github.com/golang/geo/r2"

type MaintainerUcf struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
}

func (MaintainerUcf) TableName() string {
	return "maintainer_ucf"
}

func (s MaintainerUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllFuels query all fuel services
func AllMaintainerUcfs() []MaintainerUcf {
	var services []MaintainerUcf
	Db.Find(&services)

	return services
}

//AddFuel add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddMaintainerUcf(s MaintainerUcf) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//FuelById query the fuel service by specific id
func MaintainerUcfById(id int64) MaintainerUcf {
	var service MaintainerUcf
	Db.Find(&service, id)

	return service
}

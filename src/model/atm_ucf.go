package model

import "github.com/golang/geo/r2"

type AtmUcf struct {
	Id     int64
	Lat    float32 `gorm:"column:lat"`
	Lon    float32 `gorm:"column:lon"`
	BankId int64   `gorm:"column:bank_id"`
}

//TableName determine the table name in database which is using for gorm
func (AtmUcf) TableName() string {
	return "atm_ucf"
}

//Location determine the location of service as r2.Point
func (s AtmUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtmUcfs query all the AtmUcf serivces
func AllAtmUcfs() []AtmUcf {
	var services []AtmUcf
	Db.Find(&services)

	return services
}

//AtmUcfById query the AtmUcf service by specific id
func AtmUcfById(id int64) AtmUcf {
	var service AtmUcf
	Db.Find(&service, id)

	return service
}

//AddAtmUcf add new AtmUcf service to the database
//
//return error if there is something wrong when doing transaction
func AddAtmUcf(s AtmUcf) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

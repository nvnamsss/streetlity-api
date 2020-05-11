package model

import (
	"log"

	"github.com/golang/geo/r2"
)

type Atm struct {
	Service
	BankId int64 `gorm:"column:bank_id"`
}

type Bank struct {
	Id   int64
	Name string `gorm:"column:name"`
}

//TableName determine the table name in database which is using for gorm
func (Atm) TableName() string {
	return "atm"
}

func (Bank) TableName() string {
	return "bank"
}

//Location determine the location of service as r2.Point
func (s Atm) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtms query all the atm serivces
func AllAtms() []Atm {
	var services []Atm
	Db.Find(&services)

	return services
}

//AtmById query the atm service by specific id
func AtmById(id int64) (service Atm, e error) {
	if e = Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//AtmByIds query the atm services by specific ids
func AtmByIds(ids ...int64) (services []Atm) {
	for _, id := range ids {
		s, e := AtmById(id)
		if e != nil {
			continue
		}

		services = append(services, s)
	}

	return
}

//AddAtm add new atm service to the database
//
//return error if there is something wrong when doing transaction
func AddAtm(s Atm) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//AtmsInRange query the atm services which is in the radius of a location
func AtmsInRange(p r2.Point, max_range float64) []Atm {
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

func AllBanks() []Bank {
	var banks []Bank
	Db.Find(&banks)

	return banks
}

func AddBank(s Bank) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

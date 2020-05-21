package model

import (
	"errors"
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

//FuelUcf representation the Fuel service which is confirmed
type Fuel struct {
	Service
	// Id  int64
	// Lat float32 `gorm:"column:lat"`
	// Lon float32 `gorm:"column:lon"`
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
func AddFuel(s Fuel) (e error) {
	if e = Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Fuel{}).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e := Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "add fuel", e.Error())
	}

	return
}

func queryFuel(s Fuel) (service Fuel, e error) {
	service = s

	if e := Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query fuel", e.Error())
	}

	return
}

//FuelByService get fuel by provide Service
func FuelByService(s Service) (services Fuel, e error) {
	services.Service = s
	return queryFuel(services)
}

//FuelById query the fuel service by specific id
func FuelById(id int64) (service Fuel, e error) {
	if e = Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
		return service, errors.New("Problem occured when query")
	}

	return
}

//ToiletByIds query the toilets service by specific id
func FuelByIds(ids ...int64) (services []Fuel) {
	for _, id := range ids {
		s, e := FuelById(id)
		if e != nil {
			continue
		}

		services = append(services, s)
	}

	return
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

func (s Fuel) AfterCreate(scope *gorm.Scope) (e error) {
	if e = services.AddItem(s); e != nil {
		log.Println("[Database]", "After create fuel", e.Error())
	}

	return
}

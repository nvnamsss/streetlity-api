package atm

import (
	"errors"
	"log"
	"math"
	"streelity/v1/model"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
)

type Atm struct {
	model.Service
	BankId int64 `gorm:"column:bank_id"`
}

const ServiceTableName = "atm"

var tag string = "[ATM]"

//TableName determine the table name in database which is using for gorm
func (Atm) TableName() string {
	return ServiceTableName
}

var services spatial.RTree

//Location determine the location of service as r2.Point
func (s Atm) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllServices query all the atm serivces
func AllServices() []Atm {
	var services []Atm
	model.Db.Find(&services)

	return services
}

func queryAtm(s Atm) (service Atm, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query atm", e.Error())
	}

	return
}

//ServiceByService get atm by provide Service
func ServiceByService(s model.Service) (services Atm, e error) {
	services.Service = s
	return queryAtm(services)
}

//ServiceById query the atm service by specific id
func ServiceById(id int64) (service Atm, e error) {
	db := model.Db.Find(&service, id)
	if e := db.Error; e != nil {
		log.Println("[Database]", "Atm service", id, ":", e.Error())
	}

	if db.RowsAffected == 0 {
		e = errors.New("Atm service was not found")
		log.Println("[Database]", "atm", e.Error())
	}

	return
}

//ServicesByIds query the atm services by specific ids
func ServicesByIds(ids ...int64) (services []Atm) {
	for _, id := range ids {
		s, e := ServiceById(id)
		if e != nil {
			continue
		}

		services = append(services, s)
	}

	return
}

//CreateService add new atm service to the database
//
//return error if there is something wrong when doing transaction
func CreateService(s Atm) (e error) {
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Atm{}).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "Add atm", e.Error())
		return
	}

	return nil
}

func distance(p1 r2.Point, p2 r2.Point) float64 {
	x := math.Pow(p1.X-p2.X, 2)
	y := math.Pow(p1.Y-p2.Y, 2)
	return math.Sqrt(x + y)
}

//ServicesInRange query the atm services which is in the radius of a location
func ServicesInRange(p r2.Point, max_range float64) []Atm {
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

func LoadService() {
	log.Println("[ATM]", "Loading service")

	atms := AllServices()
	for _, atm := range atms {
		services.AddItem(atm)
	}
}

func init() {
	model.OnConnected.Subscribe(LoadService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadService)
	})
}

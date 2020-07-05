package fuel

import (
	"errors"
	"log"
	"math"
	"strconv"
	"streelity/v1/model"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
	"github.com/nvnamsss/goinf/spatial"
)

//FuelUcf representation the Fuel service which is confirmed
type Fuel struct {
	model.Service
	// Id  int64
	// Lat float32 `gorm:"column:lat"`
	// Lon float32 `gorm:"column:lon"`
}

var services spatial.RTree
var map_services map[int64]Fuel

const ServiceTableName = "fuel"

//Determine table name
func (Fuel) TableName() string {
	return ServiceTableName
}

func (s Fuel) GetId() string {
	id := strconv.FormatInt(s.Id, 10)
	return id
}

func (s Fuel) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllServices query all fuel services
func AllServices() (services []Fuel, e error) {
	if e = model.Db.Find(&services).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//CreateService add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func CreateService(s Fuel) (service Fuel, e error) {
	service = s
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Fuel{}).Error; e == nil {
		return s, errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&service).Error; e != nil {
		log.Println("[Database]", "add fuel", e.Error())
	}

	return
}

//UpvoteService upvote the unconfirmed atm by specific id
func UpvoteService(id int64) error {
	return upvoteService(id, 1)
}

func UpvoteServiceImmediately(id int64) error {
	return upvoteService(id, confident)
}

func upvoteService(id int64, value int) (e error) {
	s, e := ServiceById(id)

	if e != nil {
		return e
	}

	s.Confident += value
	if e := model.Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "upvote unconfirmed fuel", id, ":", e.Error())
	}

	return
}

func queryFuel(s Fuel) (service Fuel, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query fuel", e.Error())
	}

	return
}

//ServiceByService get fuel by provide model.Service
func ServiceByService(s model.Service) (services Fuel, e error) {
	services.Service = s
	return queryFuel(services)
}

//ServiceById query the fuel service by specific id
func ServiceById(id int64) (service Fuel, e error) {
	db := model.Db.Where("id=?", id).First(&service)
	if e := db.Error; e != nil {
		log.Println("[Database]", "Fuel service", id, ":", e.Error())
	}

	if db.RowsAffected == 0 {
		e = errors.New("Fuel service was not found")
		log.Println("[Database]", "fuel", e.Error())
	}

	return
}

func ServiceByLocation(lat, lon float64) (service Fuel, e error) {
	e = model.GetServiceByLocation(ServiceTableName, lat, lon, &service)
	return
}

func ServiceByAddress(address string) (service Fuel, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &service)
	return
}

func ServicesByAddress(address string) (services []Fuel, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &services)
	return
}

//ToiletByIds query the toilets service by specific id
func ServicesByIds(ids ...int64) (services []Fuel) {
	for _, id := range ids {
		s, e := ServiceById(id)
		if e != nil {
			continue
		}

		services = append(services, s)
	}

	return
}

func distance(p1 r2.Point, p2 r2.Point) float64 {
	x := math.Pow(p1.X-p2.X, 2)
	y := math.Pow(p1.Y-p2.Y, 2)
	return math.Sqrt(x + y)
}

//ServicesInRange query the fuel services which is in the radius of a location
func ServicesInRange(p r2.Point, max_range float64) []Fuel {
	var result []Fuel = []Fuel{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isFuel := item.(Fuel)
			if isFuel && d < max_range {
				result = append(result, map_services[s.Id])
			}
		}
	}
	return result
}

func (s *Fuel) AfterSave(scope *gorm.Scope) (e error) {
	if s.Confident > confident {
		ucf_services.RemoveItem(s)
		if e = services.AddItem(*s); e != nil {
			log.Println("[Database]", "fuel offical", e.Error())
		}
	} else {
		ucf_services.AddItem(s)
	}

	map_services[s.Id] = *s
	return
}

func (s Fuel) AfterCreate(scope *gorm.Scope) (e error) {
	if e = services.AddItem(s); e != nil {
		log.Println("[Database]", "After create fuel", e.Error())
	}

	return
}

func LoadService() {
	log.Println("[Fuel]", "Loading service")
	map_services = make(map[int64]Fuel)
	ss, _ := AllServices()
	for _, s := range ss {
		if s.Confident > confident {
			services.AddItem(s)
		} else {
			ucf_services.AddItem(s)
		}
		map_services[s.Id] = s
	}
}

func init() {
	model.OnConnected.Subscribe(LoadService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadService)
	})
}

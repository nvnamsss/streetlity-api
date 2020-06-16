package toilet

import (
	"errors"
	"log"
	"math"
	"streelity/v1/model"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type Toilet struct {
	model.Service
}

var services spatial.RTree

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
	model.Db.Find(&services)

	return services
}

//AddToilet add new toilet service to the database
//
//return error if there is something wrong when doing transaction
func AddToilet(s Toilet) (e error) {
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Toilet{}).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e := model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "add toilet", e.Error())
	}

	return
}

func queryToilet(s Toilet) (service Toilet, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query toilet", e.Error())
	}

	return
}

//ToiletByService get toilet by provide a Service
func ToiletByService(s model.Service) (services Toilet, e error) {
	services.Service = s
	return queryToilet(services)
}

func ToiletById(id int64) (service Toilet, e error) {
	if e = model.Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//ToiletByIds query the toilets service by specific id
func ToiletByIds(ids ...int64) (services []Toilet) {
	for _, id := range ids {
		s, e := ToiletById(id)
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

func (s Toilet) AfterCreate(scope *gorm.Scope) (e error) {
	if e = services.AddItem(s); e != nil {
		log.Println("[Database]", "After create toilet", e.Error())
	}

	return
}

func LoadService() {
	log.Println("[Toilet]", "Loading service")

	toilets := AllToilets()
	for _, service := range toilets {
		services.AddItem(service)
	}
}

func init() {
	model.OnConnected.Subscribe(LoadService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadService)
	})
}
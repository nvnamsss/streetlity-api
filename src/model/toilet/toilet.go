package toilet

import (
	"errors"
	"log"
	"math"
	"strconv"
	"streelity/v1/model"
	"strings"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
	"github.com/nvnamsss/goinf/spatial"
)

type Toilet struct {
	model.Service
	Name string `gorm:"column:name"`
}

var services spatial.RTree
var map_services map[int64]Toilet

const ServiceTableName = "toilet"

//TableName determine the table name in database which is using for gorm
func (Toilet) TableName() string {
	return ServiceTableName
}

func (s Toilet) GetId() string {
	id := strconv.FormatInt(s.Id, 10)
	return id
}

func (s Toilet) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtms query all the atm serivces
func AllServices() (services []Toilet, e error) {
	if e = model.Db.Find(&services).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//CreateService add new toilet service to the database
//
//return error if there is something wrong when doing transaction
func CreateService(s Toilet) (service Toilet, e error) {
	service = s
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Toilet{}).Error; e == nil {
		return s, errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&service).Error; e != nil {
		log.Println("[Database]", "add toilet", e.Error())
	}

	return
}

//UpvoteService upvote the unconfirmed atm by specific id
func UpvoteService(id int64) error {
	return upvoteService(id, 1)
}

func DownvoteService(id int64) error {
	return upvoteService(id, -1)
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
		log.Println("[Database]", "upvote unconfirmed toilet", id, ":", e.Error())
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

//ServiceByService get toilet by provide a Service
func ServiceByService(s model.Service) (services Toilet, e error) {
	services.Service = s
	return queryToilet(services)
}

func ServiceById(id int64) (service Toilet, e error) {
	e = model.GetById(ServiceTableName, id, &service)
	return
}

func ServiceByLocation(lat, lon float64) (service Toilet, e error) {
	e = model.GetServiceByLocation(ServiceTableName, lat, lon, &service)
	return
}

func ServiceByAddress(address string) (service Toilet, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &service)
	return
}

func ServicesByAddress(address string) (services []Toilet, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &services)
	return
}

//ServicesByIds query the toilets service by specific id
func ServicesByIds(ids ...int64) (services []Toilet) {
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

//ServicesInRange query the toilet services which is in the radius of a location
func ServicesInRange(p r2.Point, max_range float64) []Toilet {
	var result []Toilet = []Toilet{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			s, isToilet := item.(Toilet)

			if isToilet {
				location := item.Location()
				d := distance(location, p)
				if d < max_range {
					result = append(result, map_services[s.Id])
				}
			}
		}
	}
	return result
}

func Import(bytes []byte, t string) (e error) {
	switch t {
	case "RawText":
		ImportByRawText(string(bytes))
		break
	}

	return
}

func ImportByRawText(data string) (e error) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		fields := strings.Split(line, ";")
		m := make(map[string]string)
		s := Toilet{}
		for _, field := range fields {
			log.Println(field)
			att := strings.Split(field, ":")
			if len(att) <= 1 {
				continue
			}

			m[att[0]] = att[1]
		}

		if lat, e := strconv.ParseFloat(m["lat"], 64); e != nil {
			log.Println("[Fuel]", "import", "cannot parse lat to float")
			continue
		} else {
			s.Lat = float32(lat)
		}

		if lon, e := strconv.ParseFloat(m["lon"], 64); e != nil {
			log.Println("[Fuel]", "import", "cannot parse lon to float")
			continue
		} else {
			s.Lon = float32(lon)
		}

		if address, ok := m["address"]; ok {
			s.Address = address
		}

		if note, ok := m["note"]; ok {
			s.Note = note
		}

		s.Name = m["name"]
		s.Contributor = "Streetlity"
		s.Confident = confident + 1
		CreateService(s)
	}

	return
}

func (s *Toilet) AfterSave(scope *gorm.Scope) (e error) {
	if s.Confident > confident {
		ucf_services.RemoveItem(s)
		if e = services.AddItem(*s); e != nil {
			log.Println("[Database]", "toilet offical", e.Error())
		}
	} else {
		ucf_services.AddItem(s)
	}

	map_services[s.Id] = *s
	return
}

func (s Toilet) AfterCreate(scope *gorm.Scope) (e error) {
	if e = services.AddItem(s); e != nil {
		log.Println("[Database]", "After create toilet", e.Error())
	}

	return
}

func LoadService() {
	log.Println("[Toilet]", "Loading service")
	map_services = make(map[int64]Toilet)
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

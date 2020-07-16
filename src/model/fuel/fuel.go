package fuel

import (
	"errors"
	"log"
	"math"
	"net/url"
	"strconv"
	"streelity/v1/model"
	"strings"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
	"github.com/nvnamsss/goinf/spatial"
)

//FuelUcf representation the Fuel service which is confirmed
type Fuel struct {
	model.Service
	Name string `gorm:"column:name"`
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

func UpdateService(id int64, values url.Values) (service Fuel, e error) {
	service, e = ServiceById(id)
	if e != nil {
		return
	}

	_, ok := values["lat"]
	if ok {
		if lat, e := strconv.ParseFloat(values["lat"][0], 64); e != nil {
			service.Lat = float32(lat)
		}
	}
	_, ok = values["lon"]
	if ok {
		if lon, e := strconv.ParseFloat(values["lon"][0], 64); e != nil {
			service.Lon = float32(lon)
		}
	}

	_, ok = values["note"]
	if ok {
		service.Note = values["note"][0]
	}

	if _, ok = values["address"]; ok {
		service.Address = values["address"][0]
	}

	if _, ok = values["images"]; ok {
		service.SetImages(values["images"]...)
	}

	if _, ok = values["name"]; ok {
		service.Name = values["name"][0]
	}

	if e := model.Db.Save(&service).Error; e != nil {
		log.Println("[Database]", "update ", ServiceTableName, e.Error())
	}
	return
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
		s := Fuel{}
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

func (s *Fuel) AfterSave(scope *gorm.Scope) (e error) {
	if s.Confident > confident {
		if _, ok := map_services[s.Id]; !ok {
			ucf_services.RemoveItem(s)
			if e = services.AddItem(*s); e != nil {
				log.Println("[Database]", "fuel offical", e.Error())
			}
			map_services[s.Id] = *s
			delete(map_ucfservices, s.Id)
		}
	} else {
		if _, ok := map_ucfservices[s.Id]; !ok {
			services.RemoveItem(s)
			ucf_services.AddItem(*s)
			map_ucfservices[s.Id] = *s
			delete(map_services, s.Id)
		}
	}

	return
}

// func (s Fuel) AfterCreate(scope *gorm.Scope) (e error) {
// 	if e = services.AddItem(s); e != nil {
// 		log.Println("[Database]", "After create fuel", e.Error())
// 	}

// 	return
// }

func LoadService() {
	log.Println("[Fuel]", "Loading service")
	map_services = make(map[int64]Fuel)
	map_ucfservices = make(map[int64]Fuel)
	ss, _ := AllServices()
	for _, s := range ss {
		if s.Confident > confident {
			services.AddItem(s)
			map_services[s.Id] = s
		} else {
			ucf_services.AddItem(s)
			map_ucfservices[s.Id] = s
		}
	}
}

func init() {
	model.OnConnected.Subscribe(LoadService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadService)
	})
}

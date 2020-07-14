package atm

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

type Atm struct {
	model.Service
	BankId int64 `gorm:"column:bank_id"`
}

var map_services map[int64]Atm
var services spatial.RTree
var tag string = "[ATM]"

const ServiceTableName = "atm"

//TableName determine the table name in database which is using for gorm
func (Atm) TableName() string {
	return ServiceTableName
}

func (s Atm) GetId() string {
	id := strconv.FormatInt(s.Id, 10)
	return id
}

//Location determine the location of service as r2.Point
func (s Atm) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllServices query all the atm serivces
func AllServices() (services []Atm, e error) {
	if e = model.Db.Find(&services).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
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
	e = model.GetById(ServiceTableName, id, &service)
	return
}

func ServiceByLocation(lat, lon float64) (service Atm, e error) {
	e = model.GetServiceByLocation(ServiceTableName, lat, lon, &service)
	return
}

func ServiceByAddress(address string) (service Atm, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &service)
	return
}

func ServicesByAddress(address string) (services []Atm, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &services)
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
func CreateService(s Atm) (service Atm, e error) {
	service = s
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Atm{}).Error; e == nil {
		return s, errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&service).Error; e != nil {
		log.Println("[Database]", "Add atm", e.Error())
		return
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
		log.Println("[Database]", "upvote unconfirmed atm", id, ":", e.Error())
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
		s := Atm{}
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
		if bank, e := BankByName(m["name"]); e != nil {
			continue
		} else {
			s.BankId = bank.Id
		}

		s.Contributor = "Streetlity"
		s.Confident = confident + 1
		CreateService(s)
	}

	return
}

func (s *Atm) AfterSave(scope *gorm.Scope) (e error) {
	if s.Confident > confident {
		if _, ok := map_services[s.Id]; !ok {
			ucf_services.RemoveItem(s)
			if e = services.AddItem(*s); e != nil {
				log.Println("[Database]", "atm offical", e.Error())
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

// func (s Atm) AfterCreate(scope *gorm.Scope) (e error) {
// 	if e = services.AddItem(s); e != nil {
// 		log.Println("[Database]", "After create atm", e.Error())
// 	}

// 	log.Println("[Database]", "New atm added", s)
// 	return
// }

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
				result = append(result, map_services[s.Id])
			}
		}
	}
	return result
}

func LoadService() {
	log.Println("[ATM]", "Loading service")
	map_services = make(map[int64]Atm)
	map_ucfservices = make(map[int64]Atm)

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

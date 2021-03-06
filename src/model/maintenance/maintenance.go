package maintenance

import (
	"encoding/json"
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

type Maintenance struct {
	model.Service
	Maintainer string `gorm:"column:maintainer"`
	Name       string `gorm:"column:name"`
	// Id  int64
	// Lat float32 `gorm:"column:lat"`
	// Lon float32 `gorm:"column:lon"`
}

const ServiceTableName = "maintenance"

var services spatial.RTree
var map_services map[int64]Maintenance

func (Maintenance) TableName() string {
	return ServiceTableName
}

func (s Maintenance) GetId() string {
	id := strconv.FormatInt(s.Id, 10)
	return id
}

func (s Maintenance) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllServices query all maintenance services
func AllServices() (services []Maintenance, e error) {
	if e = model.Db.Find(&services).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//CreateService add new maintenance service to the database
//
//return error if there is something wrong when doing transaction
func CreateService(s Maintenance) (service Maintenance, e error) {
	service = s
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Maintenance{}).Error; e == nil {
		return s, errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&service).Error; e != nil {
		log.Println("[Database]", "add maintennace", e.Error())
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
		log.Println("[Database]", "upvote unconfirmed maintenance", id, ":", e.Error())
	}

	return
}

func (m *Maintenance) AddMaintainer(maintainer string) (e error) {
	ms := m.GetMaintainers()
	if _, ok := ms[maintainer]; ok {
		log.Println("[Maintenance]", "Add", maintainer, "is already work for", m.Name)
		return errors.New("maintainer is already exist")
	} else {
		ms[maintainer] = "Employee"
	}

	e = m.SetMaintainer(ms)
	return
}

func (m *Maintenance) RemoveMaintainer(maintainer string) (e error) {
	ms := m.GetMaintainers()
	_, ok := ms[maintainer]
	if !ok {
		log.Println("[Maintenance]", "Remove", maintainer, "is not work for", m.Name)
		return errors.New("maintainer is not exist")
	}

	delete(ms, maintainer)
	e = m.SetMaintainer(ms)
	return
}

func (m *Maintenance) SetMaintainer(maintainers map[string]string) (e error) {
	data, e := json.Marshal(maintainers)
	if e != nil {
		log.Println("[Maintenance]", "Set maintainer", e.Error())
	}

	m.Maintainer = string(data)
	return
}

func (m Maintenance) GetMaintainers() (maintainer map[string]string) {
	maintainer = make(map[string]string)
	if e := json.Unmarshal([]byte(m.Maintainer), &maintainer); e != nil {
		log.Println("[Maintenance]", "Get maintainer", e.Error())
	}

	return maintainer
}

func queryMaintenance(s Maintenance) (service Maintenance, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query maintenance", e.Error())
	}

	return
}

//ServiceById query the maintenance service by specific id
func ServiceById(id int64) (service Maintenance, e error) {
	e = model.GetById(ServiceTableName, id, &service)
	// db := model.Db.Find(&service, id)
	// if e := db.Error; e != nil {
	// 	log.Println("[Database]", "Maintenance service", id, ":", e.Error())
	// }

	// if db.RowsAffected == 0 {
	// 	e = errors.New("Ucf Maintenance service was not found")
	// 	log.Println("[Database]", "Maintenance ucf", e.Error())
	// }

	return
}

func ServiceByLocation(lat, lon float64) (service Maintenance, e error) {
	e = model.GetServiceByLocation(ServiceTableName, lat, lon, &service)
	return
}

func ServiceByAddress(address string) (service Maintenance, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &service)
	return
}

func ServicesByAddress(address string) (services []Maintenance, e error) {
	e = model.GetServiceByAddress(ServiceTableName, address, &services)
	return
}

//ServiceByService get maintenance by provide Service
func ServiceByService(s model.Service) (services Maintenance, e error) {
	services.Service = s
	return queryMaintenance(services)
}

//ServicesByIds query the maintenances service by specific id
func ServicesByIds(ids ...int64) (services []Maintenance) {
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

//ServicesInRange query the maintenance services which is in the radius of a location
func ServicesInRange(p r2.Point, max_range float64) []Maintenance {
	var result []Maintenance = []Maintenance{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isMaintenance := item.(Maintenance)
			if isMaintenance && d < max_range {
				result = append(result, map_services[s.Id])
			}
		}
	}
	return result
}

func UpdateService(id int64, values url.Values) (service Maintenance, e error) {
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

func AddMaintainer(id int64, maintainer string) (service Maintenance, e error) {
	service, e = ServiceById(id)
	if e != nil {
		return
	}
	if e = service.AddMaintainer(maintainer); e != nil {
		return
	}

	if e = model.Db.Save(&service).Error; e != nil {
		log.Println("[Database]", "add maintainer", e.Error())
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
		s := Maintenance{}
		for _, field := range fields {
			log.Println(field)
			att := strings.Split(field, ":")
			if len(att) <= 1 {
				continue
			}

			m[att[0]] = att[1]
		}

		if lat, e := strconv.ParseFloat(m["lat"], 64); e != nil {
			log.Println("[Maintenance]", "import", "cannot parse lat to float")
			continue
		} else {
			s.Lat = float32(lat)
		}

		if lon, e := strconv.ParseFloat(m["lon"], 64); e != nil {
			log.Println("[Maintenance]", "import", "cannot parse lon to float")
			continue
		} else {
			s.Lon = float32(lon)
		}

		s.Name = m["name"]
		if address, ok := m["address"]; ok {
			s.Address = address
		}

		if note, ok := m["note"]; ok {
			s.Note = note
		}
		s.Contributor = "Streetlity"
		s.Confident = confident + 1
		CreateService(s)
	}

	return
}

func RemoveMaintainer(id int64, maintainer string) (service Maintenance, e error) {
	service, e = ServiceById(id)
	if e != nil {
		return
	}
	if e = service.RemoveMaintainer(maintainer); e != nil {
		return
	}

	if e = model.Db.Save(&service).Error; e != nil {
		log.Println("[Database]", "add maintainer", e.Error())
	}
	return
}

func (s *Maintenance) AfterSave(scope *gorm.Scope) (e error) {
	if s.Confident > confident {
		if _, ok := map_services[s.Id]; !ok {
			ucf_services.RemoveItem(s)
			if e = services.AddItem(*s); e != nil {
				log.Println("[Database]", "maintenance offical", e.Error())
			}
			delete(map_ucfservices, s.Id)
		}
		map_services[s.Id] = *s
	} else {
		if _, ok := map_ucfservices[s.Id]; !ok {
			services.RemoveItem(s)
			map_ucfservices[s.Id] = *s
			delete(map_services, s.Id)
		}
		ucf_services.AddItem(*s)
	}

	return
}

// func (s Maintenance) AfterCreate(scope *gorm.Scope) (e error) {
// 	if e = services.AddItem(s); e != nil {
// 		log.Println("[Database]", "After create maintenance", e.Error())
// 	} else {
// 		map_services[s.Id] = s
// 	}

// 	log.Println("[Database]", "New maintennace added")
// 	return
// }

func LoadService() {
	log.Println("[Maintenance]", "Loading service")
	map_services = make(map[int64]Maintenance)
	map_ucfservices = make(map[int64]Maintenance)
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

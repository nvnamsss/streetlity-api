package toilet

import (
	"errors"
	"log"
	"streelity/v1/model"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
	"github.com/nvnamsss/goinf/spatial"
)

type ToiletUcf struct {
	model.ServiceUcf
}

var confident int = 5
var map_ucfservices map[int64]Toilet
var ucf_services spatial.RTree

const UcfServiceTableName = "toilet_ucf"

//TableName determine the table name in database which is using for gorm
func (ToiletUcf) TableName() string {
	return "toilet_ucf"
}

func (s ToiletUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtms query all the atm serivces
func AllToiletUcfs() []ToiletUcf {
	var services []ToiletUcf
	model.Db.Find(&services)

	return services
}

//UpvoteUcf upvote the unconfirmed toilet by specific id
func UpvoteUcf(id int64) error {
	return upvoteToiletUcf(id, 1)
}

func UpvoteUcfImmediately(id int64) error {
	return upvoteToiletUcf(id, confident)
}

func upvoteToiletUcf(id int64, value int) (e error) {
	s, e := UcfById(id)

	if e != nil {
		return
	}

	s.Confident += value
	if e := model.Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "upvote unconfirmed toilet", id, ":", e.Error())
	}

	return
}

func queryToiletUcf(s ToiletUcf) (service ToiletUcf, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query unconfirmed atm", e.Error())
	}

	return
}

func UcfByService(s model.ServiceUcf) (service ToiletUcf, e error) {
	service.ServiceUcf = s
	return queryToiletUcf(service)
}

//UcfById query the unconfirmed toilet service by specific id
func UcfById(id int64) (service ToiletUcf, e error) {
	e = model.GetById(UcfServiceTableName, id, &service)
	return
}

func UcfByLocation(lat, lon float64) (service ToiletUcf, e error) {
	e = model.GetServiceByLocation(UcfServiceTableName, lat, lon, &service)
	return
}

func UcfByAddress(address string) (service ToiletUcf, e error) {
	e = model.GetServiceByAddress(UcfServiceTableName, address, &service)
	return
}

func UcfsByAddress(address string) (services []ToiletUcf, e error) {
	e = model.GetServiceByAddress(UcfServiceTableName, address, &services)
	return
}

//CreateUcf add new ToiletUcf service to the database
//
//return error if there is something wrong when doing transaction
func CreateUcf(s ToiletUcf) (ucf ToiletUcf, e error) {
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Toilet{}).Error; e == nil {
		return ucf, errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&ToiletUcf{}).Error; e == nil {
		return ucf, errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	} else {
		ucf = s
	}

	return
}

//UcfInRange query the unconfirmed fuel services that are in the radius of a location
func UcfInRange(p r2.Point, max_range float64) []Toilet {
	var result []Toilet = []Toilet{}
	trees := ucf_services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isService := item.(Toilet)
			if isService && d < max_range {
				result = append(result, map_ucfservices[s.Id])
			}
		}
	}
	return result
}

func DeleteUcf(id int64) (e error) {
	var ucf ToiletUcf
	ucf.Id = id
	if e := model.Db.Delete(&ucf).Error; e != nil {
		log.Println("[Database]", "delete ucf fuel", e.Error())
	}

	return
}

//AfterSave automatically run everytime the update transaction is done
func (s *ToiletUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident >= confident {
		var t Toilet = Toilet{Service: s.GetService()}
		CreateService(t)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Toilet]", "Confident is enough. Added", t)
	} else {
		ucf_services.AddItem(s)
	}

	return
}

func LoadUnconfirmedService() {
	log.Println("[Toilet]", "Loading unconfirmed service")

	toilets := AllToiletUcfs()
	for _, service := range toilets {
		ucf_services.AddItem(service)
	}
}

// func init() {
// 	model.OnConnected.Subscribe(LoadUnconfirmedService)
// 	model.OnDisconnect.Subscribe(func() {
// 		model.OnConnected.Unsubscribe(LoadUnconfirmedService)
// 	})
// }

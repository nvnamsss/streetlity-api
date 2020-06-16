package atm

import (
	"errors"
	"log"
	"streelity/v1/model"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type AtmUcf struct {
	model.ServiceUcf
	BankId int64 `gorm:"column:bank_id"`
}

var confident int = 1
var ucf_services spatial.RTree

//TableName determine the table name in database which is using for gorm
func (AtmUcf) TableName() string {
	return "atm_ucf"
}

//Location determine the location of service as r2.Point
func (s AtmUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllAtmUcfs query all the AtmUcf serivces
func AllAtmUcfs() []AtmUcf {
	var services []AtmUcf
	model.Db.Find(&services)

	return services
}

func queryAtmUcf(s AtmUcf) (service AtmUcf, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query unconfirmed atm", e.Error())
	}

	return
}

func AtmUcfByService(s model.ServiceUcf) (service AtmUcf, e error) {
	service.ServiceUcf = s
	return queryAtmUcf(service)
}

//AtmUcfById query the AtmUcf service by specific id
func AtmUcfById(id int64) (service AtmUcf, e error) {
	db := model.Db.Find(&service, id)
	if e := db.Error; e != nil {
		log.Println("[Database]", "Atm service", id, ":", e.Error())
	}

	if db.RowsAffected == 0 {
		e = errors.New("Ucf Atm service was not found")
		log.Println("[Database]", "Atm ucf", e.Error())
	}

	return
}

//UpvoteAtmUcf upvote the unconfirmed atm by specific id
func UpvoteAtmUcf(id int64) error {
	return upvoteAtmUcf(id, 1)
}

func UpvoteAtmUcfImmediately(id int64) error {
	return upvoteAtmUcf(id, confident)
}

func upvoteAtmUcf(id int64, value int) (e error) {
	s, e := AtmUcfById(id)

	if e != nil {
		return e
	}

	s.Confident += value
	if e := model.Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "upvote unconfirmed atm", id, ":", e.Error())
	}

	return
}

//AddAtmUcf add new AtmUcf service to the database
//
//return error if there is something wrong when doing transaction
func AddAtmUcf(s AtmUcf) (e error) {
	var existed AtmUcf
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&existed).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	//Temporal
	UpvoteAtmUcf(s.Id)
	return
}

//UcfInRange query the unconfirmed atm services that are in the radius of a location
func UcfInRange(p r2.Point, max_range float64) []AtmUcf {
	var result []AtmUcf = []AtmUcf{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isFuel := item.(AtmUcf)
			if isFuel && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

func DeleteUcf(id int64) (e error) {
	var ucf AtmUcf
	ucf.Id = id
	if e := model.Db.Delete(&ucf).Error; e != nil {
		log.Println("[Database]", "delete ucf fuel", e.Error())
	}

	return
}

func (s *AtmUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident >= confident {
		var a Atm = Atm{Service: s.GetService(), BankId: s.BankId}
		AddAtm(a)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Atm]", "Confident is enough. Added", a)
	} else {
		ucf_services.AddItem(s)
	}

	return
}

func LoadUnconfirmedService() {
	log.Println("[ATM]", "Loading unconfirmed service")

	maintenances := AllAtmUcfs()
	for _, service := range maintenances {
		ucf_services.AddItem(service)
	}
}

func init() {
	model.OnConnected.Subscribe(LoadUnconfirmedService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadUnconfirmedService)
	})
}

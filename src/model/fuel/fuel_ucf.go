package fuel

import (
	"errors"
	"log"
	"streelity/v1/model"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

var confident int = 5
var ucf_services spatial.RTree

//FuelUcf representation the Fuel service which is not confirmed
type FuelUcf struct {
	model.ServiceUcf
}

func (FuelUcf) TableName() string {
	return "fuel_ucf"
}

func (s FuelUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllFuelsUcf query all unconfirmed fuel services
func AllFuelsUcf() []FuelUcf {
	var services []FuelUcf
	model.Db.Find(&services)

	return services
}

//AddFuelUcf add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddFuelUcf(s FuelUcf) (e error) {
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&FuelUcf{}).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	//Temporal
	UpvoteFuelUcf(s.Id)
	return
}

func queryFuelUcf(s FuelUcf) (service FuelUcf, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query unconfirmed fuel", e.Error())
	}

	return
}

func FuelUcfByService(s model.ServiceUcf) (service FuelUcf, e error) {
	service.ServiceUcf = s
	return queryFuelUcf(service)
}

func DeleteUcf(id int64) (e error) {
	ucf := FuelUcf{model.ServiceUcf{Id: id}}
	if e := model.Db.Delete(&ucf).Error; e != nil {
		log.Println("[Database]", "delete ucf fuel", e.Error())
	}

	return
}

//UcfInRange query the unconfirmed fuel services that are in the radius of a location
func UcfInRange(p r2.Point, max_range float64) []FuelUcf {
	var result []FuelUcf = []FuelUcf{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isFuel := item.(FuelUcf)
			if isFuel && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

//FuelUcfById query the fuel service by specific id
func FuelUcfById(id int64) (service FuelUcf, e error) {
	if e = model.Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//UpvoteFuelUcf upvote the unconfirmed fuel by specific id
func UpvoteFuelUcf(id int64) error {
	return upvoteFuelUcf(id, 1)
}

func UpvoteFuelUcfImmediately(id int64) error {
	return upvoteFuelUcf(id, confident)
}

func upvoteFuelUcf(id int64, value int) (e error) {
	s, e := FuelUcfById(id)

	if e != nil {
		return
	}

	s.Confident += value
	if e := model.Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "upvote unconfirmed fuel", id, ":", e.Error())
	}

	return
}

func (s *FuelUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident >= confident {
		var f Fuel = Fuel{Service: s.GetService()}
		AddFuel(f)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Fuel]", "Confident is enough. Added", f)
	} else {
		ucf_services.AddItem(s)
	}

	return
}

func LoadUnconfirmedService() {
	log.Println("[Fuel]", "Loading unconfirmed service")

	fuels := AllFuelsUcf()
	for _, fuel := range fuels {
		ucf_services.AddItem(fuel)
	}
}

func init() {
	model.OnConnected.Subscribe(LoadUnconfirmedService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadService)
	})
}

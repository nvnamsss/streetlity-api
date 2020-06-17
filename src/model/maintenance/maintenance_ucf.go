package maintenance

import (
	"errors"
	"log"
	"streelity/v1/model"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

var confident int = 5

type MaintenanceUcf struct {
	model.ServiceUcf
	Name string `gorm:"column:name"`
}

func (MaintenanceUcf) TableName() string {
	return "maintenance_ucf"
}

func (s MaintenanceUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllUcfs query all maintenance services
func AllUcfs() []MaintenanceUcf {
	var services []MaintenanceUcf
	if e := model.Db.Find(&services).Error; e != nil {
		log.Println("[Database]", "All maintenance service", e.Error())
	}

	return services
}

//UpvoteUcf upvote the unconfirmed maintainer by specific id
func UpvoteUcf(id int64) (e error) {
	return upvoteMaintenanceUcf(id, 1)
}

//UpvoteMaintenanceUcfById upvote the unconfirmed maintainer by specific id
//with out caring about the remaining confident
func UpvoteUcfImmediately(id int64) (e error) {
	return upvoteMaintenanceUcf(id, confident)
}

func upvoteMaintenanceUcf(id int64, value int) (e error) {
	s, e := UcfById(id)

	if e != nil {
		return e
	}

	s.Confident += value
	if e = model.Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "Upvote maintenance service", id, ":", e.Error())
	}

	return
}

//CreateUcf add new unconfirmed maintainer service to the database
//
//return error if there is something wrong when doing transaction
func CreateUcf(s MaintenanceUcf) (ucf MaintenanceUcf, e error) {
	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Maintenance{}).Error; e == nil {
		return ucf, errors.New("The service location is existed")
	}

	if e = model.Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&MaintenanceUcf{}).Error; e == nil {
		return ucf, errors.New("The service location is existed or some problems is occured")
	}

	if e = model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "Add maintenance service:", e.Error())
	} else {
		ucf = s
	}

	//Temporal
	UpvoteUcfImmediately(s.Id)
	return
}

func queryMaintenanceUcf(s MaintenanceUcf) (service MaintenanceUcf, e error) {
	service = s

	if e := model.Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query unconfirmed maintenance", e.Error())
	}

	return
}

func UcfServiceByService(s model.ServiceUcf) (service MaintenanceUcf, e error) {
	service.ServiceUcf = s
	return queryMaintenanceUcf(service)
}

func MaintenanceUcfByAddress() {
}

//UcfById query the unconfirmed maintainer service by specific id
func UcfById(id int64) (service MaintenanceUcf, e error) {
	db := model.Db.Find(&service, id)
	if e := db.Error; e != nil {
		log.Println("[Database]", "Fuel service", id, ":", e.Error())
	}

	if db.RowsAffected == 0 {
		e = errors.New("Ucf Fuel service was not found")
		log.Println("[Database]", "fuel ucf", e.Error())
	}

	return
}

func DeleteUcf(id int64) (e error) {
	var ucf MaintenanceUcf
	ucf.Id = id
	if e := model.Db.Delete(&ucf).Error; e != nil {
		log.Println("[Database]", "delete ucf fuel", e.Error())
	}

	return
}

//UcfInRange query the unconfirmed fuel services that are in the radius of a location
func UcfInRange(p r2.Point, max_range float64) []MaintenanceUcf {
	var result []MaintenanceUcf = []MaintenanceUcf{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isService := item.(MaintenanceUcf)
			if isService && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

func (s *MaintenanceUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident >= confident {
		var m Maintenance = Maintenance{Service: s.GetService(), Name: s.Name}
		CreateService(m)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Maintenance]", "Confident is enough. Added", m)
	}

	return
}

func LoadUcfService() {
	log.Println("[Maintenance]", "Loading unconfirmed service")

	maintenances := AllUcfs()
	for _, service := range maintenances {
		services.AddItem(service)
	}
}

func init() {
	model.OnConnected.Subscribe(LoadUcfService)
	model.OnDisconnect.Subscribe(func() {
		model.OnConnected.Unsubscribe(LoadUcfService)
	})
}

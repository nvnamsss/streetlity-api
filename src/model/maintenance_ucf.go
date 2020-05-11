package model

import (
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type MaintenanceUcf struct {
	ServiceUcf
}

func (MaintenanceUcf) TableName() string {
	return "maintenance_ucf"
}

func (s MaintenanceUcf) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllMaintenanceUcfs query all maintainer services
func AllMaintenanceUcfs() []MaintenanceUcf {
	var services []MaintenanceUcf
	Db.Find(&services)

	return services
}

//UpvoteMaintenanceUcf upvote the unconfirmed maintainer by specific id
func UpvoteMaintenanceUcf(id int64) error {
	s, e := MaintenanceUcfById(id)

	if e != nil {
		return e
	}

	s.Confident += 1
	Db.Save(&s)

	return nil
}

//AddMaintenanceUcf add new unconfirmed maintainer service to the database
//
//return error if there is something wrong when doing transaction
func AddMaintenanceUcf(s MaintenanceUcf) (e error) {
	if e = Db.Create(&s).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

//MaintenanceUcfById query the unconfirmed maintainer service by specific id
func MaintenanceUcfById(id int64) (service MaintenanceUcf, e error) {
	if e := Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

func (s *MaintenanceUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident == 5 {
		var m Maintenance = Maintenance{Service: Service{Lat: s.Lat, Lon: s.Lon, Address: s.Address}}
		AddMaintenance(m)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Maintenance]", "Confident is enough. Added", m)
	}

	return
}

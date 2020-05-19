package model

import (
	"errors"
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type MaintenanceUcf struct {
	ServiceUcf
	Name string `gorm:"column:name"`
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
	if e := Db.Find(&services).Error; e != nil {
		log.Println("[Database]", "All maintenance service", e.Error())
	}

	return services
}

//UpvoteMaintenanceUcf upvote the unconfirmed maintainer by specific id
func UpvoteMaintenanceUcf(id int64) (e error) {
	return upvoteMaintenanceUcf(id, 1)
}

func UpvoteMaintenanceUcfImmediately(id int64) (e error) {
	return upvoteMaintenanceUcf(id, confident)
}

func upvoteMaintenanceUcf(id int64, value int) (e error) {
	s, e := MaintenanceUcfById(id)

	if e != nil {
		return e
	}

	s.Confident += value
	if e = Db.Save(&s).Error; e != nil {
		log.Println("[Database]", "Upvote maintenance service", id, ":", e.Error())
	}

	return
}

//AddMaintenanceUcf add new unconfirmed maintainer service to the database
//
//return error if there is something wrong when doing transaction
func AddMaintenanceUcf(s MaintenanceUcf) (e error) {
	var existed MaintenanceUcf
	if e = Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&existed).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e = Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "Add maintenance service:", e.Error())
	}

	//Temporal
	UpvoteMaintenanceUcf(s.Id)
	return
}

//MaintenanceUcfById query the unconfirmed maintainer service by specific id
func MaintenanceUcfById(id int64) (service MaintenanceUcf, e error) {
	if e := Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", "Maintenance service", id, ":", e.Error())
	}

	return
}

func (s *MaintenanceUcf) AfterSave(scope *gorm.Scope) (err error) {
	if s.Confident == confident {
		var m Maintenance = Maintenance{Service: s.GetService(), Name: s.Name}
		AddMaintenance(m)
		scope.DB().Delete(s)
		log.Println("[Unconfirmed Maintenance]", "Confident is enough. Added", m)
	}

	return
}

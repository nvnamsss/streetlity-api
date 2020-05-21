package model

import (
	"errors"
	"log"

	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type Maintenance struct {
	Service
	Owner string `gorm:"column:owner"`
	Name  string `gorm:"column:name"`
	// Id  int64
	// Lat float32 `gorm:"column:lat"`
	// Lon float32 `gorm:"column:lon"`
}

func (Maintenance) TableName() string {
	return "maintenance"
}

func (s Maintenance) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllMaintenances query all maintenance services
func AllMaintenances() []Maintenance {
	var services []Maintenance
	Db.Find(&services)

	return services
}

//AddMaintenance add new maintenance service to the database
//
//return error if there is something wrong when doing transaction
func AddMaintenance(s Maintenance) (e error) {
	if e = Db.Where("lat=? AND lon=?", s.Lat, s.Lon).Find(&Maintenance{}).Error; e == nil {
		return errors.New("The service location is existed or some problems is occured")
	}

	if e := Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "add maintennace", e.Error())
	}

	return
}

func queryMaintenance(s Maintenance) (service Maintenance, e error) {
	service = s

	if e := Db.Find(&service).Error; e != nil {
		log.Println("[Database]", "query maintenance", e.Error())
	}

	return
}

//MaintenanceById query the maintenance service by specific id
func MaintenanceById(id int64) (service Maintenance, e error) {
	if e = Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e)
	}

	return
}

//MaintenanceByService get maintenance by provide Service
func MaintenanceByService(s Service) (services Maintenance, e error) {
	services.Service = s
	return queryMaintenance(services)
}

//MaintenanceByIds query the maintenances service by specific id
func MaintenanceByIds(ids ...int64) (services []Maintenance) {
	for _, id := range ids {
		s, e := MaintenanceById(id)
		if e != nil {
			continue
		}

		services = append(services, s)
	}

	return
}

//MaintenancesInRange query the maintenance services which is in the radius of a location
func MaintenancesInRange(p r2.Point, max_range float64) []Maintenance {
	var result []Maintenance = []Maintenance{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isMaintenance := item.(Maintenance)
			if isMaintenance && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

func UpdateMaintenance(id int64, values map[string]string) {
	service, e := MaintenanceById(id)
	if e != nil {
		return
	}

	_, ok := values["owner"]
	if ok {
		service.Owner = values["owner"]
	}

	if e = Db.Save(&service).Error; e != nil {
		log.Println("[Database]", "Update maintenance", e.Error())
	}
}

func (s Maintenance) AfterCreate(scope *gorm.Scope) (e error) {
	if e = services.AddItem(s); e != nil {
		log.Println("[Database]", "After create maintenance", e.Error())
	}

	log.Println("[Database]", "New maintennace added")
	return
}

package model

import (
	"log"

	"github.com/golang/geo/r2"
)

type Maintenance struct {
	Service
	Owner string `gorm:"column:owner"`
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
func AddMaintenance(s Maintenance) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//MaintenanceById query the maintenance service by specific id
func MaintenanceById(id int64) (service Maintenance, e error) {
	if e = Db.Find(&service, id).Error; e != nil {
		log.Println("[Database]", e)
	}

	return
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
		log.Println("[Database]", e.Error())
	}
}

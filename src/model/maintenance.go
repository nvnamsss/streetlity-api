package model

import "github.com/golang/geo/r2"

type Maintenance struct {
	Service
	Owner int64 `gorm:"column:owner"`
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

//AllFuels query all fuel services
func AllMaintenances() []Maintenance {
	var services []Maintenance
	Db.Find(&services)

	return services
}

//AddFuel add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddMaintenance(s Maintenance) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//FuelById query the fuel service by specific id
func MaintenanceById(id int64) Maintenance {
	var service Maintenance
	Db.Find(&service, id)

	return service
}

func MaintenanceByIds(ids ...int64) (services []Maintenance) {
	for _, id := range ids {
		services = append(services, MaintenanceById(id))
	}

	return
}

//FuelsInRange query the fuel services which is in the radius of a location
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

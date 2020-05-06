package model

import "github.com/golang/geo/r2"

type Maintainer struct {
	Service
	// Id  int64
	// Lat float32 `gorm:"column:lat"`
	// Lon float32 `gorm:"column:lon"`
}

func (Maintainer) TableName() string {
	return "maintainer"
}

func (s Maintainer) Location() r2.Point {
	var p r2.Point = r2.Point{X: float64(s.Lat), Y: float64(s.Lon)}
	return p
}

//AllFuels query all fuel services
func AllMaintainers() []Maintainer {
	var services []Maintainer
	Db.Find(&services)

	return services
}

//AddFuel add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddMaintainer(s Maintainer) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//FuelById query the fuel service by specific id
func MaintainerById(id int64) Maintainer {
	var service Maintainer
	Db.Find(&service, id)

	return service
}

//FuelsInRange query the fuel services which is in the radius of a location
func MaintainersInRange(p r2.Point, max_range float64) []Maintainer {
	var result []Maintainer = []Maintainer{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			s, isMaintainer := item.(Maintainer)
			if isMaintainer && d < max_range {
				result = append(result, s)
			}
		}
	}
	return result
}

package model

//FuelUcf representation the Fuel service which is not confirmed
type FuelUcf struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
}

func (FuelUcf) TableName() string {
	return "fuel_ucf"
}

//AllFuelsUcf query all unconfirmed fuel services
func AllFuelsUcf() []FuelUcf {
	var services []FuelUcf
	Db.Find(&services)

	return services
}

//AddFuelUcf add new fuel service to the database
//
//return error if there is something wrong when doing transaction
func AddFuelUcf(s FuelUcf) error {
	if dbc := Db.Create(&s); dbc.Error != nil {
		return dbc.Error
	}

	return nil
}

//FuelUcfById query the fuel service by specific id
func FuelUcfById(id int64) FuelUcf {
	var service FuelUcf
	Db.Find(&service, id)

	return service
}

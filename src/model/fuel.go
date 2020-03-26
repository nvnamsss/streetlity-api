package model

type Fuel struct {
	Id  int64
	Lat float32 `gorm:"column:lat"`
	Lon float32 `gorm:"column:lon"`
	// Location r2.Point
}

//Determine table name
func (Fuel) TableName() string {
	return "fuel"
}

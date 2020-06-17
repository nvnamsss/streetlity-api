package model

import (
	"log"
	"math"
	"regexp"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
)

// type Services struct {
// 	Atms        []Atm
// 	Fuels       []Fuel
// 	Toilets     []Toilet
// 	Maintenance []Maintenance
// }

type ServiceUcf struct {
	Id        int64   `gorm:"column:id"`
	Lat       float32 `gorm:"column:lat"`
	Lon       float32 `gorm:"column:lon"`
	Note      string  `gorm:"column:note"`
	Address   string  `gorm:"column:address"`
	Confident int     `gorm:"column:confident"`
	Images    string  `gorm:"column:images"`
}

type Service struct {
	Id      int64   `gorm:"column:id"`
	Lat     float32 `gorm:"column:lat"`
	Lon     float32 `gorm:"column:lon"`
	Note    string  `gorm:"column:note"`
	Address string  `gorm:"column:address"`
	Images  string  `gorm:"column:images"`
}

func (s Service) GetImagesArray() (images []string) {
	reg, e := regexp.Compile(";")
	if e != nil {
		log.Println("[Database]", "wrong images data", s.Images)
		return
	}

	images = reg.Split(s.Images, -1)
	return
}

func (s Service) SetImages(images ...string) {
	len := len(images)
	if len == 0 {
		return
	}

	imgString := images[0]
	for loop := 1; loop < len; loop++ {
		imgString += ";" + images[loop]
	}

	s.Images = imgString
}

func (s ServiceUcf) GetImagesArray() (images []string) {
	reg, e := regexp.Compile(";")
	if e != nil {
		log.Println("[Database]", "wrong images data", s.Images)
		return
	}

	images = reg.Split(s.Images, -1)
	return
}

func (s *ServiceUcf) SetImages(images ...string) {
	len := len(images)
	if len == 0 {
		return
	}

	imgString := images[0]
	for loop := 1; loop < len; loop++ {
		imgString += ";" + images[loop]
	}

	s.Images = imgString
}

var services spatial.RTree

func distance(p1 r2.Point, p2 r2.Point) float64 {
	x := math.Pow(p1.X-p2.X, 2)
	y := math.Pow(p1.Y-p2.Y, 2)
	return math.Sqrt(x + y)
}

// func QueryService(s Service) {
// 	if e := Db.Find(&s).Error; e != nil {
// 		log.Println(e)
// 	}
// }

//LoadService loading all kind of service in Database and storage it into spatial tree.
//
//The functions which are using spatial tree need LoadService ran before to work as expectation.
// func LoadService() {
// 	fuels := AllFuels()
// 	atms := AllAtms()
// 	toilets := AllToilets()
// 	maintenances := AllMaintenances()

// 	for _, fuel := range fuels {
// 		services.AddItem(fuel)
// 	}

// 	for _, atm := range atms {
// 		services.AddItem(atm)
// 	}

// 	for _, toilet := range toilets {
// 		services.AddItem(toilet)
// 	}

// 	for _, maintenance := range maintenances {
// 		services.AddItem(maintenance)
// 	}

// }

func (s ServiceUcf) GetService() (service Service) {
	service.Lat = s.Lat
	service.Lon = s.Lon
	service.Note = s.Note
	service.Address = s.Address
	service.Images = s.Images

	return
}

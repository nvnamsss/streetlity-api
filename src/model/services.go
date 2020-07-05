package model

import (
	"errors"
	"log"
	"math"
	"net/url"
	"regexp"
	"strconv"

	"github.com/golang/geo/r2"
	"github.com/nvnamsss/goinf/spatial"
)

// type Services struct {
// 	Atms        []Atm
// 	Fuels       []Fuel
// 	Toilets     []Toilet
// 	Maintenance []Maintenance
// }

type ServiceUcf struct {
	Id          int64   `gorm:"column:id"`
	Lat         float32 `gorm:"column:lat"`
	Lon         float32 `gorm:"column:lon"`
	Note        string  `gorm:"column:note"`
	Address     string  `gorm:"column:address"`
	Confident   int     `gorm:"column:confident"`
	Images      string  `gorm:"column:images"`
	Contributor string  `gorm:"column:contributor"`
}

type Service struct {
	Id          int64   `gorm:"column:id"`
	Lat         float32 `gorm:"column:lat"`
	Lon         float32 `gorm:"column:lon"`
	Note        string  `gorm:"column:note"`
	Address     string  `gorm:"column:address"`
	Images      string  `gorm:"column:images"`
	Contributor string  `gorm:"column:contributor"`
	Confident   int     `gorm:"column:confident"`
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

func (s *Service) SetImages(images ...string) {
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

func (s ServiceUcf) GetId() string {
	id := strconv.FormatInt(s.Id, 10)
	return id
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
	service.Contributor = s.Contributor
	return
}

func GetServiceByLocation(tablename string, lat, lon float64, ref interface{}) (e error) {
	db := Db.Table(tablename).Where("ABS(lat-?) < 0.00001 AND ABS(lon-?) < 0.00001", lat, lon).Find(ref)
	e = db.Error
	if db.RowsAffected == 0 {
		e := errors.New("record was not found")
		log.Println("[Database]", "get by location", tablename, e.Error())
	}
	return
}

func GetServiceByAddress(tablename string, address string, ref interface{}) (e error) {
	db := Db.Table(tablename).Where("address LIKE ?", "%"+address+"%").First(ref)
	e = db.Error

	if db.RowsAffected == 0 {
		e := errors.New("record was not found")
		log.Println("[Database]", "get by address", tablename, e.Error())
	}
	return
}

func UpdateService(tablename string, id int64, values url.Values, ref interface{}) (e error) {
	if e = GetById(tablename, id, ref); e != nil {
		return
	}

	service := ref.(Service)
	_, ok := values["lat"]
	if ok {
		if lat, e := strconv.ParseFloat(values["lat"][0], 64); e != nil {
			service.Lat = float32(lat)
		}
	}
	_, ok = values["lon"]
	if ok {
		if lon, e := strconv.ParseFloat(values["lon"][0], 64); e != nil {
			service.Lon = float32(lon)
		}
	}

	_, ok = values["note"]
	if ok {
		service.Note = values["lon"][0]
	}

	if _, ok = values["address"]; ok {
		service.Address = values["address"][0]
	}

	if _, ok = values["images"]; ok {
		service.SetImages(values["images"]...)
	}

	if e := Db.Save(service).Error; e != nil {
		log.Println("[Database]", "update ", tablename, e.Error())
	}

	return
}

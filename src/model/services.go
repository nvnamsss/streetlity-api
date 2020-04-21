package model

import (
	"fmt"
	"math"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
)

type Services struct {
	Atms    []Atm
	Fuels   []Fuel
	Toilets []Toilet
}

type Service interface {
}

var services spatial.RTree

func distance(p1 r2.Point, p2 r2.Point) float64 {
	x := math.Pow(p1.X-p2.X, 2)
	y := math.Pow(p1.Y-p2.Y, 2)
	return math.Sqrt(x + y)
}

func ServicesInRange(p r2.Point, max_range float64) []Service {
	var result []Service = []Service{}
	trees := services.InRange(p, max_range)

	for _, tree := range trees {
		for _, item := range tree.Items {
			location := item.Location()

			d := distance(location, p)
			if d < max_range {
				result = append(result, item)
			}
		}
	}

	fmt.Println(result)

	return result
}

//LoadService loading all kind of service in Database and storage it into spatial tree.
//
//The functions which are using spatial tree need LoadService ran before to work as expectation.
func LoadService() {
	fuels := AllFuels()
	atms := AllAtms()
	toilets := AllToilets()

	for _, fuel := range fuels {
		services.AddItem(fuel)
	}

	for _, atm := range atms {
		services.AddItem(atm)
	}

	for _, toilet := range toilets {
		services.AddItem(toilet)
	}
}

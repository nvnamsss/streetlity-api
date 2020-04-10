package model

import (
	"math"
	"streelity/v1/spatial"

	"github.com/golang/geo/r2"
)

type Service interface {
}

var services spatial.RTree

func ServicesInRange() {
	p := r2.Point{X: 5, Y: 5}
	services.Find(func(rtree *spatial.RTree) bool {
		p2 := rtree.Rect.Center()
		d := math.Sqrt(math.Pow(p2.X-p.X, 2) + math.Pow(p2.Y-p.Y, 2))

		return d > 5
	})
}

func LoadService() {
	fuels := AllFuels()

	for _, fuel := range fuels {
		services.AddItem(fuel)
	}

	ServicesInRange()
}

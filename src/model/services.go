package model

import (
	"streelity/v1/spatial"
)

var services spatial.RTree

func meomeo() {
	var f *Fuel = new(Fuel)
	item := spatial.Item(f)
	services.Items = append(services.Items, &item)

}

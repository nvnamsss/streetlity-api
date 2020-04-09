package spatial

import (
	"errors"
	"fmt"
	"math"

	"github.com/golang/geo/r2"
)

type RTree struct {
	Ancestor   *RTree
	Descendant []*RTree
	Items      []*Item
	Rect       r2.Rect
	MaxItem    int
}

type Item interface {
	GetLocation() r2.Point
}

//Add new cell to current cell
func (r *RTree) Add(l *RTree) error {
	if l == nil {
		return errors.New("Additional RTree is nil")
	}

	if l.Ancestor == r {
		return errors.New("Additional RTree is already existed in current RTree")
	}

	// if len is greater than MaxItems, add new cell instead
	if len(r.Descendant) > r.MaxItem {

	}

	r.Descendant = append(r.Descendant, l)
	l.Ancestor = r

	return nil
}

//Add new item to tree, new node will be added if required node is not exist
func (r *RTree) AddItem(item Item) error {
	location := item.GetLocation()
	fmt.Println(location)

	return nil
}
func (r *RTree) Find(location r2.Point) *RTree {

	return nil
}

//Get the distance between two RTree
//Location is measure by center of two RTree
func (r RTree) Distance(l RTree) float64 {
	var p1 r2.Point = r2.Point{X: r.Rect.X.Lo, Y: r.Rect.Y.Lo}
	var p2 r2.Point = r2.Point{X: r.Rect.X.Hi, Y: r.Rect.Y.Hi}

	return math.Sqrt(math.Pow(p1.X-p2.X, 2) + math.Pow(p1.Y-p2.Y, 2))
}

//Create new pointer RTree
func NewRTree() *RTree {
	qt := new(RTree)
	qt.Ancestor = nil

	return qt
}

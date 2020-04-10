package spatial

import (
	"errors"
	"math"

	"github.com/golang/geo/r1"
	"github.com/golang/geo/r2"
)

//RTree repesentation a tree for spatial indexing which improve the time for searching a location in space
type RTree struct {
	Ancestor   *RTree
	Descendant []*RTree
	Items      []Item
	Rect       r2.Rect
	MaxItem    int
	level      int
}

type Item interface {
	Location() r2.Point
}

//AddTree new RTree to current RTree
func (r *RTree) AddTree(l *RTree) error {
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
	l.level = r.level + 1
	return nil
}

//The AddItem add new item to the tree, new node will be added if required node is not exist
func (rt *RTree) AddItem(item Item) error {
	tree := rt.Nearest(item.Location(), 20)

	if tree == nil {
		tree = NewRTree()
		rt.AddTree(tree)
		location := item.Location()
		tree.Rect.X.Lo = location.X - 0.5
		tree.Rect.Y.Lo = location.Y - 0.5
		tree.Rect.X.Hi = location.X + 0.5
		tree.Rect.Y.Hi = location.X + 0.5
	}

	tree.Items = append(tree.Items, item)
	tree.UpdateRect()

	return nil
}

//UpdateRect modify the Rect size and affect to the Rect of Ancestor
//UpdateRect should be used for automatically updating size of Rect
func (rt *RTree) UpdateRect() {
	if len(rt.Items) == 0 {
		return
	}

	location := rt.Items[0].Location()
	var x r1.Interval = r1.Interval{Lo: location.X, Hi: location.X}
	var y r1.Interval = r1.Interval{Lo: location.Y, Hi: location.Y}

	for _, item := range rt.Items {
		p := item.Location()

		x.Lo = math.Min(x.Lo, p.X)
		x.Hi = math.Max(x.Hi, p.X)

		y.Lo = math.Min(y.Lo, p.Y)
		y.Hi = math.Max(y.Hi, p.Y)
	}

	rt.Rect.X = x
	rt.Rect.Y = y
}

//Nearest find the nearest RTree in Descendant which is in max_range from the location
//nil will be return if there are no RTree
func (rt *RTree) Nearest(location r2.Point, max_range float64) *RTree {
	var min float64 = math.MaxFloat64
	var result *RTree = nil
	for _, tree := range rt.Descendant {
		d := tree.Distance(location)

		if d < min {
			d = min
			result = tree
		}
	}

	return result
}

//Find the Tree which is holding finding item
//The smallest tree will be returned if found
//otherwise the return value will be nil
func (rt *RTree) Contains(item Item) *RTree {
	location := item.Location()
	var result *RTree = nil

	for _, element := range rt.Descendant {
		if element.Rect.ContainsPoint(location) {
			small := element.Contains(item)
			if small != nil {
				result = small
			} else {
				result = element
			}
		}
	}

	return result
}

func (rt RTree) Level() int {
	level := 0

	for rt.Ancestor != nil {
		level += 1
		rt = *rt.Ancestor
	}

	return level
}

//Find all RTree which are matched with the function
func (rt *RTree) Find(match func(tree *RTree) bool) []RTree {
	var result []RTree
	if match(rt) {
		result = append(result, *rt)

		for _, item := range rt.Descendant {
			result = append(result, item.Find(match)...)
		}
	}

	return result
}

//Get the distance between two RTree
//Location is measure by center of two RTree
func (r RTree) Distance(p r2.Point) float64 {
	center := r.Rect.Center()
	x := math.Pow(center.X-p.X, 2)
	y := math.Pow(center.Y-p.Y, 2)

	return math.Sqrt(x + y)
}

//Create new pointer RTree
func NewRTree() *RTree {
	qt := new(RTree)
	qt.Ancestor = nil
	qt.level = 0

	return qt
}

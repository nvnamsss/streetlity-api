package spatial

import (
	"errors"
	"math"

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
	GetLocation() r2.Point
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
func (r *RTree) AddItem(item Item) error {
	tree := r.Contains(item)

	if tree == nil {
		tree = NewRTree()
		tree.Items = append(tree.Items, item)
		r.AddTree(tree)
		location := item.GetLocation()
		tree.Rect.X.Lo = location.X - 0.5
		tree.Rect.Y.Lo = location.Y - 0.5
		tree.Rect.X.Hi = location.X + 0.5
		tree.Rect.Y.Hi = location.X + 0.5
	} else {
		tree.Items = append(tree.Items, item)
	}

	return nil
}

//Find the Tree which is holding finding item
//The smallest tree will be returned if found
//otherwise the return value will be nil
func (r *RTree) Contains(item Item) *RTree {
	location := item.GetLocation()
	var result *RTree = nil

	for _, element := range r.Descendant {
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

func (r RTree) Level() int {
	level := 0

	for r.Ancestor != nil {
		level += 1
		r = *r.Ancestor
	}

	return level
}

//Find all RTree which are matched with the function
func (r *RTree) Find(match func(tree *RTree) bool) []RTree {
	var result []RTree
	if match(r) {
		result = append(result, *r)

		for _, item := range r.Descendant {
			result = append(result, item.Find(match)...)
		}
	}

	return result
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
	qt.level = 0

	return qt
}

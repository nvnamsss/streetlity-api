package spatial

import (
	"errors"

	"github.com/golang/geo/r2"
)

type Quadtree struct {
	Ancestor   *Quadtree
	Descendant []*Quadtree
	Items      []*Item
	Rect       r2.Rect
}

type Item interface {
}

//Add new cell to current cell
func (r *Quadtree) Add(l *Quadtree) error {
	if l == nil {
		return errors.New("Additional Quadtree is nil")
	}

	if l.Ancestor == r {
		return errors.New("Additional Quadtree is already existed in current Quadtree")
	}

	r.Descendant = append(r.Descendant, l)
	l.Ancestor = r

	return nil
}

func (*Quadtree) InsertItem(item Item) {

}

func NewQuadtree() *Quadtree {
	qt := new(Quadtree)
	qt.Ancestor = nil

	return qt
}

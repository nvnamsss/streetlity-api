package Astar

import r2 "github.com/golang/geo/r2"

type Node struct {
	Id        int64
	Location  r2.Point
	Neighbors []Node
	Data      NodeData /*Data perspective with neightbors*/
	StreetId  []int64
}

type NodeData struct {
	G float64
	H float64
	F float64
}

type SortNode []Node

func (a SortNode) Len() int           { return len(a) }
func (a SortNode) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortNode) Less(i, j int) bool { return a[i].Data.F < a[j].Data.F }

func NewNode(id int64, location r2.Point, neighbors []Node, data NodeData) *Node {
	node := new(Node)
	node.Id = id
	node.Location = location
	node.Neighbors = neighbors
	node.Data = data

	return node
}

type Path struct {
	Parent *Path
	Node   Node
}

func NewPath(parent *Path, node Node) *Path {
	p := new(Path)
	p.Parent = parent
	p.Node = node
	return p
}

type Street struct {
	Id      int64
	NodeIds []int64
	Cost    int64
}

func NewStreet(id int64, nodeIds []int64) *Street {
	s := new(Street)
	s.Id = id
	s.NodeIds = nodeIds

	return s
}

func (s *Street) Nullable() (nullable, ok bool) {
	return nullable, true
}

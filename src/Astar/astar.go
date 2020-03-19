package main

import (
	"fmt"
	"math"
	"sort"

	r2 "github.com/golang/geo/r2"
)

var Nodes []Node
var openList []Node
var closeList map[int]bool

/*
openlist: The open list is a collection of all generated Nodes.
This means that those are Nodes that were Neighbors of expanded Nodes.
As mentioned above, the open list is often implemented as a priority queue
so the search can simply dequeue the nest best Node.

closedlist: The closed list is a collection of all expanded Nodes.
This means that those are Nodes that were already "searched".
in big domains, the closed list can't fit all Nodes,
so the closed list has to be implemented smartly.
For example, it is possible to reduce the memory required using a Bloom Filter.
This prevents the search from visiting Nodes again and again
*/
type Node struct {
	Id        int
	Location  r2.Point
	Neighbors []Node
	Data      NodeData /*Data perspective with neightbors*/
	Street    *Street
}

type Street struct {
	Id int
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

func d(from r2.Point, to r2.Point) (value float64) {
	return math.Sqrt(math.Pow(from.X-to.X, 2) + math.Pow(from.Y-to.Y, 2))
}

/*estimate distance from current Node to end Node*/
func h_euclid(from r2.Point, to r2.Point) (value float64) {
	return math.Sqrt(math.Pow(from.X-to.X, 2) + math.Pow(from.Y-to.Y, 2))
}

func IndexOf(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func Remove(s []Node, index int) []Node {
	return append(s[:index], s[index+1:]...)
}

func RemoveItem(s []Node, item Node) []Node {
	var index int = IndexOf(len(s), func(i int) bool { return s[i].Id == item.Id })
	return Remove(s, index)
}

func GetCurrent(s []Node) Node {
	return s[len(s)-1]
}

func prepend(s []Node, node Node) []Node {
	s = append(s, Node{})
	copy(s[1:], s)
	s[0] = node
	return s
}

func reconstruct_path(path *Path) []Node {
	var nodes []Node
	for path.Parent != nil {
		nodes = prepend(nodes, path.Node)
		path = path.Parent
	}

	return nodes
}

/*find the shortest path*/
func route(from Node, to Node) ([]Node, bool) {
	var path *Path = NewPath(nil, from)
	var current Node = from
	current.Data.G = 0
	current.Data.H = h_euclid(current.Location, to.Location)
	current.Data.F = current.Data.H

	openList = append(openList, current)
	closeList = make(map[int]bool)

	fmt.Println("[Astar]", "First node F", current.Data.F)

	for len(openList) > 0 {
		current = GetCurrent(openList)
		closeList[current.Id] = true

		if current.Id == to.Id {
			return reconstruct_path(path), true
		}

		path = NewPath(path, current)
		openList = RemoveItem(openList, current)

		for index, element := range current.Neighbors {
			var eG float64 = current.Data.G + d(current.Location, element.Location)
			var eH float64 = h_euclid(element.Location, to.Location)
			var eF float64 = eG + eH
			closeList[element.Id] = false

			fmt.Println("[Astar]", current.Data.F, eF)
			if eF <= current.Data.F {
				// var pos = IndexOf(len(openList), func(i int) bool { return openList[i].Id == element.Id })
				if !closeList[element.Id] {
					element.Data.G = eG
					element.Data.H = eH
					element.Data.F = eF
					openList = append(openList, element)

					var s SortNode = openList
					sort.Stable(s)
				}

			}
			fmt.Println("[Foreach]", index, element)
		}
	}

	return nil, false
}

func Test() {
	var p1 r2.Point = r2.Point{X: 1, Y: 2}
	var p2 r2.Point = r2.Point{X: 3, Y: 4}
	var p3 r2.Point = r2.Point{X: 4, Y: 5}
	var p4 r2.Point = r2.Point{X: 10, Y: 5}
	var n1 Node = Node{Id: 1, Location: p1}
	var n2 Node = Node{Id: 2, Location: p2}
	var n3 Node = Node{Id: 3, Location: p3}
	var n4 Node = Node{Id: 4, Location: p4}

	n1.Neighbors = []Node{n2}
	n2.Neighbors = []Node{n3}
	n2.Neighbors = []Node{n4}
	var path, result = route(n1, n4)

	if result {
		fmt.Println("[AStar]", "Success")
		fmt.Println("[AStar]", path)
	} else {
		fmt.Println("[AStar", "Cannot found any path")
	}

	fmt.Println("[AStar]", h_euclid(p1, p2))
}

func init() {

}

func main() {
	Test()
}

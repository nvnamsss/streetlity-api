package Astar

import (
	"fmt"
	"math"
	"sort"

	r2 "github.com/golang/geo/r2"
)

var Streets map[int64]Street = make(map[int64]Street)
var Nodes map[int64]Node = make(map[int64]Node)

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

func getNextNode(s []Node) Node {
	return s[0]
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
func Route(from Node, to Node) ([]Node, bool) {
	fmt.Println(from, to)
	var openList []Node = []Node{}
	var closeList map[int64]bool = make(map[int64]bool)
	var path *Path = NewPath(nil, from)
	var current Node = from
	var paths map[int64]*Path = make(map[int64]*Path)

	current.Data.G = 0
	current.Data.H = h_euclid(current.Location, to.Location)
	current.Data.F = math.MaxFloat64

	openList = append(openList, current)
	closeList = make(map[int64]bool)

	for len(openList) > 0 {
		current = getNextNode(openList)
		closeList[current.Id] = true
		fmt.Println("[Astar]", "Current", current)
		path = NewPath(path, current)
		paths[current.Id] = path

		if current.Id == to.Id {
			fmt.Println("[Astar]", "Path is found")
			return reconstruct_path(path), true
		}

		openList = RemoveItem(openList, current)

		for _, item := range current.StreetId {
			var street Street = Streets[item]

			for _, item2 := range street.NodeIds {
				if closeList[item2] {
					continue
				}

				var eG float64 = d(current.Location, Nodes[item2].Location)
				var eH float64 = h_euclid(Nodes[item2].Location, to.Location)
				var eF float64 = eG + eH
				// fmt.Println("[Astar]", "Compare F", current.Data.F, eF)
				if _, ok := closeList[item2]; !ok {
					closeList[item2] = false
				}

				if eF <= current.Data.F {
					if !closeList[item2] {
						data := Nodes[item2]
						data.Data.G = eG
						data.Data.H = eH
						data.Data.F = eF
						Nodes[item2] = data

						openList = append(openList, Nodes[item2])

						var s SortNode = openList
						sort.Stable(s)
					}
				}
			}

		}

	}

	fmt.Println("[Astar]", "Cannot find any path")
	return nil, false
}

func Test() {
	var p1 r2.Point = r2.Point{X: 1, Y: 2}
	var p2 r2.Point = r2.Point{X: 3, Y: 4}
	var p3 r2.Point = r2.Point{X: 4, Y: 5}
	var p4 r2.Point = r2.Point{X: 10, Y: 5}
	var n1 *Node = new(Node)
	n1.Id = 1
	n1.Location = p1
	var n2 *Node = new(Node)
	n2.Id = 2
	n2.Location = p2
	var n3 *Node = new(Node)
	n3.Id = 3
	n3.Location = p3
	var n4 *Node = new(Node)
	n4.Id = 4
	n4.Location = p4

	n2.Neighbors = []Node{*n3, *n4}
	n1.Neighbors = []Node{*n2}

	var path, result = Route(*n1, *n4)

	if result {
		fmt.Println("[AStar]", "Success")
		fmt.Println("[AStar]", path)
	} else {
		fmt.Println("[AStar", "Cannot found any path")
	}

}

func main() {
	Test()
}

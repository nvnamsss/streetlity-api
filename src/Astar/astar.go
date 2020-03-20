package Astar

import (
	"fmt"
	"math"

	r2 "github.com/golang/geo/r2"
)

var start Node
var target Node
var openList []Node

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
	Data      NodeData
}

type NodeData struct {
	G float64
	H float64
	F float64
}

func d(from r2.Point, to r2.Point) (value float64) {
	return math.Sqrt(math.Pow(from.X-to.X, 2) + math.Pow(from.Y-to.Y, 2))
}

/*distance from current Node to start Node*/
func g(from r2.Point, to r2.Point) (value float64) {
	return 0
}

/*estimate distance from current Node to end Node*/
func h(from r2.Point, to r2.Point) (value float64) {
	var p r2.Point
	p.X = 1
	p.Y = 2
	fmt.Println(p)
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

/*find the shortest path*/
func route(from Node, to Node, newH float32) []Node {
	var path []Node
	var current Node = from
	openList = append(openList, from)

	for len(openList) > 0 {
		//var nearestNode Node //nearest Node
		if current.Id == to.Id {
			return path
		}

		//remove current out openList
		RemoveItem(openList, current)
		//iterate neighbor
		for index, element := range current.Neighbors {
			var eG float64 = g(element.Location, from.Location) + d(current.Location, element.Location)
			var eH float64 = h(current.Location, to.Location)
			var eF float64 = eG + eH

			if eF < current.Data.F {
				var pos = IndexOf(len(openList), func(i int) bool { return openList[i].Id == element.Id })
				//do task
				if pos == -1 {
					openList = append(openList, element)
					element.Data.G = eG
					element.Data.H = eH
					element.Data.F = eF
				}

			}
			// index is the index where we are
			// element is the element from someSlice for where we are
			fmt.Println(eG)
			fmt.Println("[Foreach]", index, element)
		}
	}
	//initialize the open list
	//initialize the closed list

	// while the open list is not empty
	// a) find the Node with the least f on
	//    the open list, call it "q"

	// b) pop q off the open list

	// c) generate q's 8 successors and set their
	//    parents to q

	// d) for each successor
	//     i) if successor is the goal, stop search
	//       successor.g = q.g + distance between
	//                           successor and q
	//       successor.h = distance from goal to
	//       successor (This can be done using many
	//       ways, we will discuss three heuristics-
	//       Manhattan, Diagonal and Euclidean
	//       Heuristics)

	//       successor.f = successor.g + successor.h

	//     ii) if a Node with the same position as
	//         successor is in the OPEN list which has a
	//        lower f than successor, skip this successor

	//     iii) if a Node with the same position as
	//         successor  is in the CLOSED list which has
	//         a lower f than successor, skip this successor
	//         otherwise, add  the Node to the open list
	//  end (for loop)

	// e) push q on the closed list
	// end (while loop)

	return path
}

func init() {
	var p1 r2.Point = r2.Point{X: 1, Y: 2}
	var p2 r2.Point = r2.Point{X: 3, Y: 4}

	fmt.Println("[AStar]", h(p1, p2))
}

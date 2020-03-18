package Astar

import (
	"fmt"
	"math"

	r2 "github.com/golang/geo/r2"
)

var start node
var target node
var openList []node

/*
openlist: The open list is a collection of all generated nodes.
This means that those are nodes that were neighbors of expanded nodes.
As mentioned above, the open list is often implemented as a priority queue
so the search can simply dequeue the nest best node.

closedlist: The closed list is a collection of all expanded nodes.
This means that those are nodes that were already "searched".
in big domains, the closed list can't fit all nodes,
so the closed list has to be implemented smartly.
For example, it is possible to reduce the memory required using a Bloom Filter.
This prevents the search from visiting nodes again and again
*/
type node struct {
	id        int
	location  r2.Point
	neighbors []node
}

/*distance from current node to start node*/
func g(from r2.Point, to r2.Point) (h float64) {
	return 0
}

/*estimate distance from current node to end node*/
func h_euclid(from r2.Point, to r2.Point) (h float64) {
	var p r2.Point
	p.X = 1
	p.Y = 2
	fmt.Println(p)
	return math.Sqrt(math.Pow(from.X-to.X, 2) + math.Pow(from.Y-to.Y, 2))
}

func remove() {

}

/*find the shortest path*/
func route(from node, to node, h float32) []node {
	var path []node
	var current node = from
	openList = append(openList, from)

	for len(openList) > 0 {
		var nearestNode node //nearest node
		if current.id == to.id {
			return path
		}

		//remove current out openList
		remove()

		//iterate neighbor
		for index, element := range current.neighbors {
			var eG float64 = g(element.location, start.location) + h_euclid(current.location, element.location)

			// index is the index where we are
			// element is the element from someSlice for where we are
			fmt.Println(eG)
			fmt.Println("[Foreach]", index, element)
		}
	}
	//initialize the open list
	//initialize the closed list

	// while the open list is not empty
	// a) find the node with the least f on
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

	//     ii) if a node with the same position as
	//         successor is in the OPEN list which has a
	//        lower f than successor, skip this successor

	//     iii) if a node with the same position as
	//         successor  is in the CLOSED list which has
	//         a lower f than successor, skip this successor
	//         otherwise, add  the node to the open list
	//  end (for loop)

	// e) push q on the closed list
	// end (while loop)

	return path
}

func init() {
	var p1 r2.Point = r2.Point{X: 1, Y: 2}
	var p2 r2.Point = r2.Point{X: 3, Y: 4}

	fmt.Println("[AStar]", h_euclid(p1, p2))
}

package Astar

import (
	"fmt"
	"math"

	r2 "github.com/golang/geo/r2"
)

type node struct {
	id int
}

func h_euclid(from r2.Point, to r2.Point) (h float64) {
	var p r2.Point
	p.X = 1
	p.Y = 2
	fmt.Println(p)
	return math.Sqrt(math.Pow(from.X-to.X, 2) + math.Pow(from.Y-to.Y, 2))
}

func route(from node, to node, h float32) []node {
	var path []node
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

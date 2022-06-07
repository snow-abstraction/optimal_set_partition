package main

import (
	"math"

	"github.com/snow-abstraction/optimal_set_partition/tree"
)

func main() {

	var unprocessed []*tree.Node
	root := tree.CreateRoot()
	unprocessed = append(unprocessed, root)
	k3, k4 := tree.Branch(root, math.MaxFloat64, 3, 4)
	unprocessed = append(unprocessed, k3, k4)
	// for len(unprocessed) > 0 {

	// 	node := unprocessed[0]
	// 	unprocessed[0] = nil
	// 	unprocessed = unprocessed[1:]

	// 	// find lb
	// 	// check if integral solution found
	// 	// find branching rows/constraints
	// 	// prune or branch

	// 		printTree([]*Node{node})

	// }

	tree.PrintTree(unprocessed)
}

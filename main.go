package main

import (
	"fmt"
	"math"
)

type NodeKind byte

const (
	root NodeKind = 0
	// In the "both" branch subproblem the two branching
	// constraints should be covered by the same variable
	bothBranch = 1
	// In the "diff" branch subproblem the two branching
	// constrains should be covered by different variables
	diffBranch = 2
)

// constraint branch-and-bound node
// The subproblem the node represents can be calculated by applying
// the branch type of it and its ancestors.
type node struct {
	kind       NodeKind
	parent     *node // nil if root node
	lowerBound float64
	// The following have no meaning for the root node
	branchConstraintOne uint32
	branchConstraintTwo uint32
}

func createRoot() *node {
	return &node{root, nil, math.MaxFloat64, math.MaxUint32, math.MaxUint32}
}

// Branches the parent node on the two constrains to create two new Nodes
func Branch(parent *node, lowerBound float64, branchConstraintOne uint32,
	branchConstraintTwo uint32) (*node, *node) {

	if parent == nil {
		panic("Cannot branch nil node.")
	}

	return &node{bothBranch, parent, lowerBound, branchConstraintOne, branchConstraintTwo},
		&node{diffBranch, parent, lowerBound, branchConstraintOne, branchConstraintTwo}

}

// For printing the implicit tree struct of Nodes
type PrintNode struct {
	referenceNode   *node
	bothBranchChild *PrintNode
	diffBranchChild *PrintNode
}

// For the start node and its ancestors, create corresponding PrintNodes if
// they are not already in printNodeByNode. And set the links for the PrintNodes
// from the parent to its children.
func add(printNodeByNode map[*node]*PrintNode, start *node) *PrintNode {
	curr := start // curr = current
	var prev *node

	var prevPNode *PrintNode
	var currPNode *PrintNode

	// Isn't necessary to actually follow parents to the root node
	// if a node is already in PNodeByNode, but we do so to
	// check for errors.
	for curr != nil {
		var ok bool
		currPNode, ok = printNodeByNode[curr]
		if !ok {
			currPNode = &PrintNode{curr, nil, nil}
			printNodeByNode[curr] = currPNode
		}

		if prev != nil {
			switch prev.kind {
			case root:
				panic(fmt.Sprintf("Node of kind root has a non-nil paren %+v.", *curr))
			case bothBranch:
				if currPNode.bothBranchChild == nil {
					currPNode.bothBranchChild = prevPNode
				} else if currPNode.bothBranchChild != prevPNode {
					panic(fmt.Sprintf(
						"bothBranchChild set before to a different node for node %+v.", *curr))
				}
			case diffBranch:
				if currPNode.diffBranchChild == nil {
					currPNode.diffBranchChild = prevPNode
				} else if currPNode.diffBranchChild != prevPNode {
					panic(fmt.Sprintf(
						"diffBranchChild set before to a different node for node %+v.", *curr))
				}
			default:
				panic(fmt.Sprintf("Unknown kind for for node %+v.", *curr))

			}
		}

		prev = curr
		prevPNode = currPNode
		curr = curr.parent

	}

	return currPNode
}

func printImpl(depth int, node *PrintNode) {
	if node == nil {
		return
	}

	for i := 0; i < depth; i++ {
		fmt.Printf(" ")
	}

	fmt.Printf("%+v\n", *node.referenceNode)
	printImpl(depth+2, node.diffBranchChild)
	printImpl(depth+2, node.diffBranchChild)

}

// For the nodes, find all ancestors and print the tree of nodes
func printTree(nodes []*node) {
	if len(nodes) == 0 {
		return
	}

	var root *PrintNode
	m := make(map[*node]*PrintNode)
	for _, node := range nodes {
		r := add(m, node)
		if root != nil && r != root {
			panic("Two different root nodes found.")
		}
		root = r
	}

	printImpl(0, root)

}

func main() {

	var unprocessed []*node
	root := createRoot()
	unprocessed = append(unprocessed, root)
	k1 := &node{bothBranch, root, math.MaxFloat64, 1, 2}
	k2 := &node{diffBranch, root, math.MaxFloat64, 1, 2}
	unprocessed = append(unprocessed, k1)
	unprocessed = append(unprocessed, k2)
	k3, k4 := Branch(k2, math.MaxFloat64, 3, 4)
	unprocessed = append(unprocessed, k3, k4)
	// for len(unprocessed) > 0 {

	// 	node := unprocessed[0]
	// 	unprocessed[0] = nil
	// 	unprocessed = unprocessed[1:]

	// 	// find lb
	// 	// check if integral soluion found
	// 	// find branching rows/constraints
	// 	// prune or branch

	// 		printTree([]*Node{node})

	// }

	printTree(unprocessed)
}

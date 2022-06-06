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
type Node struct {
	kind       NodeKind
	parent     *Node // nil if root node
	lowerBound float64
	// The following have no meaning for the root node
	branchConstraintOne uint32
	branchConstraintTwo uint32
}

// For printing the imlicit tree struct of Nodes
type PrintNode struct {
	node            *Node
	bothBranchChild *PrintNode
	diffBranchChild *PrintNode
}

// returns root
func add(pNodeByNode map[*Node]*PrintNode, node *Node) *PrintNode {
	var prev *Node
	curr := node

	var prevPNode *PrintNode
	var currPNode *PrintNode

	// Isn't necessary to actually follow parents to the root node
	// if a node is already in PNodeByNode, but we do so to
	// check for errors.
	for curr != nil {
		var ok bool
		currPNode, ok = pNodeByNode[curr]
		if !ok {
			currPNode = &PrintNode{curr, nil, nil}
			pNodeByNode[curr] = currPNode
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

	fmt.Printf("%+v\n", *node.node)
	printImpl(depth+2, node.diffBranchChild)
	printImpl(depth+2, node.diffBranchChild)

}

// For the nodes, find all ancestors and print the tree of nodes
func printTree(nodes []*Node) {
	if len(nodes) == 0 {
		return
	}

	var root *PrintNode
	m := make(map[*Node]*PrintNode)
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

	var unprocessed []*Node
	unprocessed = append(unprocessed, &Node{root, nil, math.MaxFloat64, 0, 0})

	for len(unprocessed) > 0 {

		node := unprocessed[0]
		unprocessed[0] = nil
		unprocessed = unprocessed[1:]

		// find lb
		// check if integral soluion found
		// find branching rows/constraints
		// prune or branch

		printTree([]*Node{node})

	}

	printTree(unprocessed)
}

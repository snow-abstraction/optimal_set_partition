package tree

import (
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

// constraint branch-and-bound Node
// The subproblem the Node represents can be calculated by applying
// the branch type of it and its ancestors.
type Node struct {
	kind       NodeKind
	parent     *Node // nil if root node
	lowerBound float64
	// The following have no meaning for the root node
	branchConstraintOne uint32
	branchConstraintTwo uint32
}

func CreateRoot() *Node {
	return &Node{root, nil, math.MaxFloat64, math.MaxUint32, math.MaxUint32}
}

// Branches the parent node on the two constrains to create two new Nodes
func Branch(parent *Node, lowerBound float64, branchConstraintOne uint32,
	branchConstraintTwo uint32) (*Node, *Node) {

	if parent == nil {
		panic("Cannot branch nil node.")
	}

	return &Node{bothBranch, parent, lowerBound, branchConstraintOne, branchConstraintTwo},
		&Node{diffBranch, parent, lowerBound, branchConstraintOne, branchConstraintTwo}

}

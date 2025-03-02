package production

import "fmt"

// Edge represents a belt of resources between two machines.
type Edge struct {
	resource Resource
	fromNode *Node
	toNode   *Node
}

func NewEdge(fromNode *Node, toNode *Node) *Edge {
	return &Edge{
		fromNode: fromNode,
		toNode:   toNode,
	}
}

func (e *Edge) PrintEdge() {
	fmt.Printf("%s <- %s", e.fromNode.Recipe.Name, e.toNode.Recipe.Name)
	fmt.Println()
}

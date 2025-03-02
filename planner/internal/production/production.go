package production

import (
	"github.com/turbo-pioneer/planner/internal"
)

type ProductionLine struct {
	nodes []*Node
	edges []*Edge
}

func NewProductionLine() *ProductionLine {
	return &ProductionLine{
		nodes: make([]*Node, 0),
		edges: make([]*Edge, 0),
	}
}

func (p *ProductionLine) GetNodes() []*Node {
	return p.nodes
}

func (p *ProductionLine) GetEdges() []*Edge {
	return p.edges
}

type LineBuilder struct {
	dataRegistry *internal.DataRegistry
}

func NewLineBuilder(dr *internal.DataRegistry) *LineBuilder {
	return &LineBuilder{
		dataRegistry: dr,
	}
}

func (l *LineBuilder) GenerateNode(recipeName string) (*Node, error) {
	r, err := l.dataRegistry.GetRecipe(recipeName)
	if err != nil {
		return nil, err
	}
	n := NewNode()
	n.Recipe = r

	var inputs = make([]*Resource, len(r.Ingredients))
	for idx, v := range r.Ingredients {
		i, err := l.dataRegistry.GetItem(v.Item)
		if err != nil {
			return nil, err
		}
		inputs[idx] = NewResource(i, v.Amount, r.Time)
	}

	n.Inputs = inputs

	var outputs = make([]*Resource, len(r.Products))
	for idx, v := range r.Products {
		i, err := l.dataRegistry.GetItem(v.Item)
		if err != nil {
			return nil, err
		}
		outputs[idx] = NewResource(i, v.Amount, r.Time)
	}

	n.Outputs = outputs

	return n, nil
}

func (l *LineBuilder) CreateProductionLineFromRecipe(recipeName string) (*ProductionLine, error) {
	prod := NewProductionLine()
	if err := l.generateLine(recipeName, prod, nil); err == nil {
		return prod, nil
	} else {
		return nil, err
	}
}

func (l *LineBuilder) generateLine(recipeName string, productionLine *ProductionLine, parentNode *Node) error {
	n, err := l.GenerateNode(recipeName)
	if err != nil {
		return err
	}
	if parentNode == nil {
		n.Root = true
	} else {
		e := NewEdge(parentNode, n)
		productionLine.edges = append(productionLine.edges, e)
	}
	productionLine.nodes = append(productionLine.nodes, n)
	for _, v := range n.Inputs {
		name := v.Item.Name
		l.generateLine(name, productionLine, n)
	}
	return nil
}

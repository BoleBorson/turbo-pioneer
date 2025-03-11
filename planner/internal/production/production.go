package production

import (
	"fmt"
	"math"

	"github.com/turbo-pioneer/planner/internal"
	"github.com/turbo-pioneer/planner/internal/utils"
)

type ProductionLine struct {
	nodes  []*Node
	edges  []*Edge
	excess map[string]float64 // excess resources expressed in resource/min
}

func NewProductionLine() *ProductionLine {
	return &ProductionLine{
		nodes:  make([]*Node, 0),
		edges:  make([]*Edge, 0),
		excess: make(map[string]float64),
	}
}

func (p *ProductionLine) GetNodes() []*Node {
	return p.nodes
}

func (p *ProductionLine) GetEdges() []*Edge {
	return p.edges
}

func (p *ProductionLine) CreateNodeMap() map[string][]*Node {
	var nodeMap = map[string][]*Node{}
	for _, n := range p.nodes {
		nodeMap[n.Recipe.Name] = append(nodeMap[n.Recipe.Name], n)
	}
	return nodeMap
}

func (p *ProductionLine) PrintRequiredRates() {
	var rates = map[string]float64{}
	for _, n := range p.nodes {
		if n.Root {
			continue
		}
		for _, out := range n.Outputs {
			rates[out.Item.Name] += out.ExpectedRate
		}
	}
	fmt.Println("Required Resource Rates:")
	for key, value := range rates {
		fmt.Printf("%.2f %s/min\n", value, key)
	}
}

func (p *ProductionLine) Print() {
	var root *Node
	for _, n := range p.nodes {
		if n.Root {
			root = n
		}
	}
	p.printTree(root, p.edges, "", true)
}

func (p *ProductionLine) printTree(root *Node, edges []*Edge, indent string, isLast bool) {
	if isLast {
		fmt.Printf("%s└── %s: %.2f %s/min\n", indent, root.Machine.Name, root.Outputs[0].ExpectedRate, root.Recipe.Name)
	} else {
		fmt.Printf("%s├── %s: %.2f %s/min\n", indent, root.Machine.Name, root.Outputs[0].ExpectedRate, root.Recipe.Name)
	}

	for i, edge := range edges {
		if edge.fromNode == root {
			p.printTree(edge.toNode, edges, indent+"│   ", i == len(edges)-1)
		}
	}
}

func (p *ProductionLine) PrintExcess() {
	for k, v := range p.excess {
		fmt.Printf("Excesss %.2f %s/min\n", v, k)
	}
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

	building, err := l.dataRegistry.GetBuilding(r.ProducedIn[0])
	if err != nil {
		return nil, err
	}
	n.Machine = building

	return n, nil
}

func (l *LineBuilder) CreateProductionLineFromRecipe(recipeName string, rate float64) (*ProductionLine, error) {
	prod := NewProductionLine()

	r, err := l.dataRegistry.GetRecipe(recipeName)
	if err != nil {
		return nil, err
	}

	i, err := l.dataRegistry.GetItem(r.Products[0].Item)
	if err != nil {
		return nil, err
	}

	root := NewRootNode(r, i, rate)
	prod.nodes = append(prod.nodes, root)
	// TODO: calculate what rate machines should be run at. Ideally most machines should run at default rate and than one runs at a funky number to make up the difference.
	// numMachines := int(math.Round(rate / utils.Rate(r.Products[0].Amount, r.Time)))
	// ratePerMachine := rate / float64(numMachines)
	if err := l.generateLine(recipeName, prod, root, rate); err != nil {
		return nil, err
	}
	// for i := 0; i < numMachines; i++ {
	// 	if err := l.generateLine(recipeName, prod, root, ratePerMachine); err != nil {
	// 		return nil, err
	// 	}
	// }
	return prod, nil
}

func (l *LineBuilder) generateLine(recipeName string, productionLine *ProductionLine, parentNode *Node, expectedOutputRate float64) error {
	r, err := l.dataRegistry.GetRecipe(recipeName)
	if err != nil {
		return err
	}
	recipeOutputRate := utils.Rate(r.Products[0].Amount, r.Time)
	numMachines := math.Ceil(expectedOutputRate / recipeOutputRate)
	excess := (numMachines * recipeOutputRate) - expectedOutputRate
	if excess > 0 {
		productionLine.excess[r.Name] = productionLine.excess[r.Name] + excess
	}
	for range int(numMachines) {

		n, err := l.GenerateNode(recipeName)
		if err != nil {
			return err
		}

		e := NewEdge(parentNode, n)
		productionLine.edges = append(productionLine.edges, e)
		productionLine.nodes = append(productionLine.nodes, n)

		for _, v := range n.Inputs {
			l.generateLine(v.Item.Name, productionLine, n, v.ExpectedRate)
		}
	}

	return nil
}

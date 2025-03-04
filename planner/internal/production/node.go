package production

import (
	"github.com/turbo-pioneer/planner/internal/building"
	"github.com/turbo-pioneer/planner/internal/item"
	"github.com/turbo-pioneer/planner/internal/recipe"
)

// Node represents a Machine. Most notably it provides information about what input rates it accepts and what output rate it releases.
type Node struct {
	Recipe        *recipe.Recipe
	Root          bool
	Machine       *building.Building
	ScalingFactor float64
	Inputs        []*Resource
	Outputs       []*Resource
}

func NewNode() *Node {
	return &Node{
		Root: false,
	}
}

// NewRootNode creates a Node the represents the expected output of the production line. It accepts the recipe, a single "product" which is the output of the recipe, and the rate the product should be produced at expressed in product per minute).
func NewRootNode(recipe *recipe.Recipe, product *item.Item, rate float64) *Node {
	input := NewResourceFromRate(product, rate)
	return &Node{
		Root:    true,
		Recipe:  recipe,
		Inputs:  []*Resource{input},
		Outputs: []*Resource{input},
		Machine: &building.Building{Name: "Output"},
	}
}

type Resource struct {
	Item *item.Item
	Rate float64 // rate represents the number of items produced per min
}

func NewResource(item *item.Item, amount float64, time int) *Resource {
	rate := (amount / float64(time)) * 60
	return &Resource{
		Item: item,
		Rate: rate,
	}
}

func NewResourceFromRate(item *item.Item, rate float64) *Resource {
	return &Resource{
		Item: item,
		Rate: rate,
	}
}

type Infrastructure struct {
	machineType      *building.Building
	numberOfMachines int
	ratePerMachine   float64
}

func NewInfrastructure(machineType *building.Building, number int, rate float64) *Infrastructure {
	return &Infrastructure{
		machineType:      machineType,
		numberOfMachines: number,
		ratePerMachine:   rate,
	}
}

package production

import (
	"github.com/turbo-pioneer/planner/internal/item"
	"github.com/turbo-pioneer/planner/internal/recipe"
)

type Node struct {
	Recipe *recipe.Recipe
	Machine any
	Inputs []*Resource
	Outputs []*Resource
}

func NewNode() *Node {
	return &Node{}
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
package production

import (
	"github.com/turbo-pioneer/planner/internal/item"
	"github.com/turbo-pioneer/planner/internal/recipe"
)

type Node struct {
	recipe recipe.Recipe
	machine any
	inputs []Resource
	outputs []Resource
}

type Resource struct {
	item item.Item
	rate float64 // rate represents the number of items produced per min
}
package production

import (
	"sync"
	"time"

	"github.com/turbo-pioneer/planner/internal/building"
	"github.com/turbo-pioneer/planner/internal/item"
	"github.com/turbo-pioneer/planner/internal/recipe"
	"github.com/turbo-pioneer/planner/internal/utils"
)

// Node represents a Machine. Most notably it provides information about what input rates it accepts and what output rate it releases.
type Node struct {
	Recipe        *recipe.Recipe
	Root          bool
	Machine       *building.Building
	ScalingFactor float64
	Inputs        map[string]*Resource
	Outputs       map[string]*Resource
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
		Inputs:  map[string]*Resource{input.Item.ClassName: input},
		Outputs: map[string]*Resource{input.Item.ClassName: input},
		Machine: &building.Building{Name: "Output"},
	}
}

func (n *Node) Produce() {
	var wg sync.WaitGroup
	for _, ingredient := range n.Recipe.Ingredients {
		go n.PopResources(ingredient, &wg)
	}
	wg.Wait()
	// once all resouces have been recieved we will then produce the output
	time.Sleep(time.Duration(n.Recipe.Time)) // simulate producing the item
	for _, p := range n.Recipe.Products {
		r := n.Outputs[p.Item]
		r.Buffer.Push(r.Item)
	}
}

func (n *Node) PopResources(ingredient *recipe.Ingredient, wg *sync.WaitGroup) {
	r := n.Inputs[ingredient.Item]
	for i := 0; i < ingredient.Amount; i++ {
		if item := r.Buffer.Pop(); item != nil {
			continue // we discard the item as it was used to produce the output
		} else {
			time.Sleep(time.Second)
			i--
			continue
		}
	}
	wg.Done()
}

func (n *Node) PushResources(wg *sync.WaitGroup) {

}

type Resource struct {
	Item           *item.Item
	ItemsProcessed int
	StartTime      time.Time
	ExpectedRate   float64 // rate represents the number of items produced per min
	ActualRate     float64
	Buffer         *utils.Buffer
	Port           chan item.Item
}

func NewResource(itemObj *item.Item, amount float64, time int) *Resource {
	rate := (amount / float64(time)) * 60 // set expected rate at start, sim will adjust this value
	return &Resource{
		Item:         itemObj,
		ExpectedRate: rate,
		ActualRate:   0,
		Buffer:       utils.NewBuffer(),
		Port:         make(chan item.Item),
	}
}

func NewResourceFromRate(itemObj *item.Item, rate float64) *Resource {
	return &Resource{
		Item:         itemObj,
		ExpectedRate: rate,
		ActualRate:   0,
		Buffer:       utils.NewBuffer(),
		Port:         make(chan item.Item),
	}
}

func (r *Resource) Access() {
	if r.ItemsProcessed == 0 {
		r.StartTime = time.Now()
	}
	r.ItemsProcessed++

	elapsedTime := time.Since(r.StartTime).Minutes()
	r.ActualRate = float64(r.ItemsProcessed) / elapsedTime

	r.Buffer.Push(<-r.Port)
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

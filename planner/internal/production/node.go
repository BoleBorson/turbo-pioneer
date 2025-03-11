package production

import (
	"fmt"
	"log/slog"
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
	Inputs        map[string]*Port
	Outputs       map[string]*Port
}

func NewNode(recipe *recipe.Recipe, rate float64, building *building.Building) *Node {
	var pullPorts = map[string]*Port{}
	for _, input := range recipe.Ingredients {
		pullPorts[input.Item] = NewPort(&item.Item{}, input.Amount, recipe.Time)
	}
	var pushPorts = map[string]*Port{}
	for _, output := range recipe.Products {
		pushPorts[output.Item] = NewPort(&item.Item{}, output.Amount, recipe.Time)
	}

	return &Node{
		Root:    false,
		Recipe:  recipe,
		Machine: building,
		Inputs:  pullPorts,
		Outputs: pushPorts,
	}
}

// NewRootNode creates a Node the represents the expected output of the production line. It accepts the recipe, a single "product" which is the output of the recipe, and the rate the product should be produced at expressed in product per minute).
func NewRootNode(recipe *recipe.Recipe, product *item.Item, rate float64) *Node {
	input := NewPortFromRate(product, rate)
	return &Node{
		Root:    true,
		Recipe:  recipe,
		Inputs:  map[string]*Port{input.Item.ClassName: input},
		Outputs: map[string]*Port{input.Item.ClassName: input},
		Machine: &building.Building{Name: "Output"},
	}
}

func (n *Node) ConnectInput(belt chan *item.Item, itemName string) error {
	if port, ok := n.Inputs[itemName]; ok {
		port.Connection = belt
	} else {
		return fmt.Errorf("no input port on this node accepts: %s", itemName)
	}
	return nil
}

func (n *Node) ConnectOutput(belt chan *item.Item, itemName string) error {
	if port, ok := n.Outputs[itemName]; ok {
		port.Connection = belt
	} else {
		return fmt.Errorf("no output port on this node accepts: %s", itemName)
	}
	return nil
}

// Start is used to startup all the goroutines required to run the machine.
func (n *Node) Start() (chan struct{}, error) {
	var done = make(chan struct{})
	// startup inputs
	for _, p := range n.Inputs {
		if p.Connection == nil {
			return nil, fmt.Errorf("Can't start node, port %s has no belt connected", p.Item.Name)
		}
		go p.Pull(done)
	}
	for _, p := range n.Outputs {
		// a node can produce even if its outputs has no belts connected
		go p.Push(done)
	}
	go n.Produce(done)
	return done, nil
}

func (n *Node) Produce(done chan struct{}) {
	// loop indefinitily as machines should be constantly producing
	for {
		select {
		case <-done:
			return
		default:
			var wg sync.WaitGroup
			for _, ingredient := range n.Recipe.Ingredients {
				wg.Add(1)
				go n.PopResources(ingredient, &wg)
			}
			wg.Wait()
			// once all resouces have been recieved we will then produce the output
			time.Sleep(time.Duration(n.Recipe.Time) * time.Second) // simulate producing the item
			for _, p := range n.Recipe.Products {
				slog.Info("produced item", "item", p.Item)
				wg.Add(1)
				go n.PushResources(p, &wg)
			}
			wg.Wait()
		}
	}
}

// PopResources provides the Node a way to pull items out of it's Pull Port Buffers. The produce function must wait until all calls to PopResources are done.
// Then the produce function can create the new item.
func (n *Node) PopResources(ingredient *recipe.Ingredient, wg *sync.WaitGroup) {
	r := n.Inputs[ingredient.Item]
	for i := 0; i < int(ingredient.Amount); i++ {
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

// PushResources provides the Node a way to push resources into its Push Ports Buffers.
func (n *Node) PushResources(p *recipe.Product, wg *sync.WaitGroup) {
	r := n.Outputs[p.Item]
	r.Buffer.Push(r.Item)
	wg.Done()
}

type Port struct {
	Item           *item.Item
	ItemsProcessed int
	StartTime      time.Time
	ExpectedRate   float64 // rate represents the number of items produced per min
	ActualRate     float64
	Buffer         *utils.Buffer
	Connection     chan *item.Item // represents the place where a belt connects to a machine.
	Flag           bool
}

func NewPort(itemObj *item.Item, amount float64, time int) *Port {
	rate := (amount / float64(time)) * 60 // set expected rate at start, sim will adjust this value
	return &Port{
		Item:         itemObj,
		ExpectedRate: rate,
		ActualRate:   0,
		Buffer:       utils.NewBuffer(),
		Connection:   make(chan *item.Item, 100),
	}
}

func NewPortFromRate(itemObj *item.Item, rate float64) *Port {
	return &Port{
		Item:         itemObj,
		ExpectedRate: rate,
		ActualRate:   0,
		Buffer:       utils.NewBuffer(),
		Connection:   make(chan *item.Item, 100),
	}
}

func (p *Port) Pull(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case i, ok := <-p.Connection:
			if ok {
				// set start time on first item recieved or it will take forever for the rate to correct itself
				if p.ItemsProcessed == 0 {
					p.StartTime = time.Now()
				}
				// caculate current input rate
				p.ItemsProcessed++
				elapsedTime := time.Since(p.StartTime).Minutes()
				p.ActualRate = float64(p.ItemsProcessed) / elapsedTime
				if p.ActualRate > p.ExpectedRate {
					p.ActualRate = p.ExpectedRate
				}

				p.Buffer.Push(i)
			} else {
				// if the channel is closed, we stop trying to pull items
				return
			}
		default:
			if !p.Flag {
				p.Flag = true
				slog.Info(fmt.Sprintf("No more items on the belt item: %s rate: %.2f", p.Item.Name, p.ActualRate))
			}
			// if there is no resource yet we keep polling until it comes available
			continue
		}

	}
}

func (p *Port) Push(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			i := p.Buffer.Pop()
			if i != nil {
				continue
			}
			if p.ItemsProcessed == 0 {
				p.StartTime = time.Now()
			}
			// caculate current output rate
			p.ItemsProcessed++
			elapsedTime := time.Since(p.StartTime).Minutes()
			p.ActualRate = float64(p.ItemsProcessed) / elapsedTime

			// send the item onto the next belt
			p.Connection <- i
		}
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

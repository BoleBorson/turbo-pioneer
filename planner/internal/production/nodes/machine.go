package nodes

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/turbo-pioneer/planner/internal/models"
	"github.com/turbo-pioneer/planner/internal/production/port"
)

// Machine represents a models. Most notably it provides information about what input rates it accepts and what output rate it releases.
type Machine struct {
	Recipe        *models.Recipe
	Root          bool
	IsProducing   bool
	Building      *models.Building
	ScalingFactor float64
	Inputs        map[string]*port.Port
	Outputs       map[string]*port.Port
}

func NewMachineNode(r *models.Recipe, rate float64, building *models.Building) *Machine {
	var pullPorts = map[string]*port.Port{}
	for _, input := range r.Ingredients {
		pullPorts[input.Item] = port.NewPort(&models.Item{}, input.Amount, r.Time)
	}
	var pushPorts = map[string]*port.Port{}
	for _, output := range r.Products {
		pushPorts[output.Item] = port.NewPort(&models.Item{}, output.Amount, r.Time)
	}

	return &Machine{
		Root:     false,
		Recipe:   r,
		Building: building,
		Inputs:   pullPorts,
		Outputs:  pushPorts,
	}
}

// NewRootMachine creates a Machine the represents the expected output of the production line. It accepts the recipe, a single "product" which is the output of the recipe, and the rate the product should be produced at expressed in product per minute).
func NewRootMachine(recipe *models.Recipe, product *models.Item, rate float64) *Machine {
	input := port.NewPortFromRate(product, rate)
	return &Machine{
		Root:     true,
		Recipe:   recipe,
		Inputs:   map[string]*port.Port{input.Item.ClassName: input},
		Outputs:  map[string]*port.Port{input.Item.ClassName: input},
		Building: &models.Building{Name: "Output"},
	}
}

func (n *Machine) ConnectInput(belt chan *models.Item, itemName string) (*port.Port, error) {
	if p, ok := n.Inputs[itemName]; ok {
		p.Connection = belt
		return p, nil
	} else {
		return nil, fmt.Errorf("no input port.Port on this Machine accepts: %s", itemName)
	}
}

func (n *Machine) ConnectOutput(belt chan *models.Item, itemName string) error {
	if p, ok := n.Outputs[itemName]; ok {
		p.Connection = belt
	} else {
		return fmt.Errorf("no output port.Port on this Machine accepts: %s", itemName)
	}
	return nil
}

// Start is used to startup all the goroutines required to run the models.
func (n *Machine) Start() (chan struct{}, error) {
	var done = make(chan struct{})
	// startup inputs
	for _, p := range n.Inputs {
		if p.Connection == nil {
			return nil, fmt.Errorf("can't start Machine, port.Port %s has no belt connected", p.Item.Name)
		}
		go p.Pull(done)
	}
	for _, p := range n.Outputs {
		// a Machine can produce even if its outputs has no belts connected
		go p.Push(done)
	}
	go n.Produce(done)
	return done, nil
}

func (n *Machine) Produce(done chan struct{}) {
	// loop indefinitily as Buildings should be constantly producing
	for {
		select {
		case <-done:
			return
		default:
			var wg sync.WaitGroup
			n.IsProducing = false
			for _, ingredient := range n.Recipe.Ingredients {
				wg.Add(1)
				go n.popResources(ingredient, &wg)
			}
			wg.Wait()
			n.IsProducing = true
			// once all resouces have been recieved we will then produce the output
			time.Sleep(time.Duration(n.Recipe.Time) * time.Second) // simulate producing the item
			for _, p := range n.Recipe.Products {
				slog.Info("produced item", "item", p.Item)
				wg.Add(1)
				go n.pushResources(p, &wg)
			}
			wg.Wait()
		}
	}
}

// PopResources provides the Machine a way to pull items out of it's Pull port.Port Buffers. The produce function must wait until all calls to PopResources are done.
// Then the produce function can create the new models.
func (n *Machine) popResources(ingredient *models.Ingredient, wg *sync.WaitGroup) {
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

// PushResources provides the Machine a way to push resources into its Push Ports Buffers.
func (n *Machine) pushResources(p *models.Product, wg *sync.WaitGroup) {
	r := n.Outputs[p.Item]
	r.Buffer.Push(r.Item)
	wg.Done()
}

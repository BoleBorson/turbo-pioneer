package nodes

import (
	"fmt"
	"sync"

	"github.com/turbo-pioneer/planner/internal/models"
	"github.com/turbo-pioneer/planner/internal/production/port"
)

type Miner struct {
	Resource      *models.Recipe
	MinerType     *models.Miner
	ScalingFactor float64
	Inputs        map[string]*port.Port
	Outputs       map[string]*port.Port
}

func (m *Miner) ConnectInput(belt chan *models.Item, itemName string) error {
	return fmt.Errorf("error miners do not support inputs")
}

func (m *Miner) ConnectOutput(belt chan *models.Item, itemName string) error {
	if p, ok := m.Outputs[itemName]; ok {
		p.Connection = belt
	} else {
		return fmt.Errorf("no output Port on this Machine accepts: %s", itemName)
	}
	return nil
}

// PushResources provides the Machine a way to push resources into its Push Ports Buffers.
func (m *Miner) pushResources(p *models.Product, wg *sync.WaitGroup) {
	r := m.Outputs[p.Item]
	r.Buffer.Push(r.Item)
	wg.Done()
}

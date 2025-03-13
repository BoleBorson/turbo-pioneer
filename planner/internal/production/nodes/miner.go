package nodes

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/turbo-pioneer/planner/internal/models"
	"github.com/turbo-pioneer/planner/internal/production/port"
)

type Miner struct {
	ResourceNode   *ResourceNode
	MinerType      *models.Miner
	ScalingFactor  float64
	ProductionTime time.Duration
	Inputs         map[string]*port.Port
	Outputs        map[string]*port.Port
}

func NewMiner(rn *ResourceNode, mt *models.Miner, sf float64) *Miner {
	rate := (rn.purity / int(mt.ExtractCycleTime)) * mt.ItemsPerCycle
	pt := time.Duration((1.0 / float64(rate)) * 60 * float64(time.Second))

	var pushPorts = map[string]*port.Port{}
	pushPorts[rn.item.ClassName] = port.NewPortFromRate(rn.item, float64(rate))
	return &Miner{
		ResourceNode:   rn,
		MinerType:      mt,
		ScalingFactor:  sf,
		ProductionTime: pt,
		Outputs:        pushPorts,
	}
}

func (m *Miner) ConnectInput(belt chan *models.Item, itemName string) error {
	return fmt.Errorf("error miners do not support inputs")
}

func (m *Miner) ConnectOutput(belt chan *models.Item, itemName string) (*port.Port, error) {
	if p, ok := m.Outputs[itemName]; ok {
		p.Connection = belt
		return p, nil
	} else {
		return nil, fmt.Errorf("no output Port on this Machine accepts: %s", itemName)
	}
}

// Start is used to startup all the goroutines required to run the models.
func (m *Miner) Start() (chan struct{}, error) {
	var done = make(chan struct{})
	for _, p := range m.Outputs {
		// a Machine can produce even if its outputs has no belts connected
		go p.Push(done)
	}
	go m.Produce(done)
	return done, nil
}

func (m *Miner) Produce(done chan struct{}) {
	// loop indefinitily as Miners should be constantly producing
	for {
		select {
		case <-done:
			return
		default:
			var wg sync.WaitGroup
			slog.Info("production time", "time", m.ProductionTime)
			time.Sleep(m.ProductionTime) // simulate producing the item
			wg.Add(1)
			slog.Info("pushing item", "item", m.ResourceNode.item.ClassName)
			go m.pushResources(m.ResourceNode.item.ClassName, &wg)
			wg.Wait()
		}
	}
}

// PushResources provides the Machine a way to push resources into its Push Ports Buffers.
func (m *Miner) pushResources(itemName string, wg *sync.WaitGroup) {
	r := m.Outputs[itemName]
	r.Buffer.Push(r.Item)
	wg.Done()
}

const (
	Impure = 30
	Normal = 60
	Pure   = 120
)

type ResourceNode struct {
	item   *models.Item
	purity int
}

func NewResourceNode(i *models.Item, p int) *ResourceNode {
	return &ResourceNode{
		item:   i,
		purity: p,
	}
}

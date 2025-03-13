package port

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/turbo-pioneer/planner/internal/models"
	"github.com/turbo-pioneer/planner/internal/utils"
)

type Port struct {
	Item           *models.Item
	ItemsProcessed int
	StartTime      time.Time
	ExpectedRate   float64 // rate represents the number of items produced per min
	ActualRate     float64
	Buffer         *utils.Buffer
	Connection     chan *models.Item // represents the place where a belt connects to a Building.
	Flag           bool
}

func NewPort(itemObj *models.Item, amount float64, time int) *Port {
	rate := (amount / float64(time)) * 60 // set expected rate at start, sim will adjust this value
	return &Port{
		Item:         itemObj,
		ExpectedRate: rate,
		ActualRate:   0,
		Buffer:       utils.NewBuffer(),
		Connection:   make(chan *models.Item, 100),
	}
}

func NewPortFromRate(itemObj *models.Item, rate float64) *Port {
	return &Port{
		Item:         itemObj,
		ExpectedRate: rate,
		ActualRate:   0,
		Buffer:       utils.NewBuffer(),
		Connection:   make(chan *models.Item, 100),
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

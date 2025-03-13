package production

import (
	"time"

	"github.com/turbo-pioneer/planner/internal/models"
	"github.com/turbo-pioneer/planner/internal/production/port"
	"github.com/turbo-pioneer/planner/internal/utils"
)

const (
	MK1 = 60.0
	MK2 = 120.0
	MK3 = 270.0
	MK4 = 480.0
	MK5 = 780.0
	MK6 = 1200.0
)

// Edge represents a belt of resources between two machines.
type Edge struct {
	Buffer         *utils.Buffer
	ItemsProcessed int
	TimeOnBelt     time.Duration
	fromPort       *port.Port
	fromConn       chan *models.Item
	toPort         *port.Port
	toConn         chan *models.Item
}

func NewEdge(fromPort *port.Port, toPort *port.Port, speed float64) *Edge {
	fromConn := make(chan *models.Item, 100)
	toConn := make(chan *models.Item, 100)

	fromPort.Connection = fromConn
	toPort.Connection = toConn

	time := time.Duration((1.0 / speed) * 60 * float64(time.Second))

	return &Edge{
		fromPort:   fromPort,
		fromConn:   fromConn,
		toPort:     toPort,
		toConn:     toConn,
		TimeOnBelt: time,
		Buffer:     utils.NewBuffer(),
	}
}

func (e *Edge) Start() chan struct{} {
	done := make(chan struct{})
	go e.pull(done)
	go e.push(done)
	return done
}

func (e *Edge) pull(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		case i, ok := <-e.fromConn:
			if ok {
				e.Buffer.Push(i)
			} else {
				return
			}
		default:
			continue
		}

	}
}

func (e *Edge) push(done chan struct{}) {
	for {
		select {
		case <-done:
			return
		default:
			time.Sleep(e.TimeOnBelt)
			i := e.Buffer.Pop()
			if i != nil {
				continue
			}
			e.toConn <- i
			e.ItemsProcessed++
		}
	}
}

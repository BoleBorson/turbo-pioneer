package production

import (
	"time"

	"github.com/turbo-pioneer/planner/internal/item"
	"github.com/turbo-pioneer/planner/internal/utils"
)

// Edge represents a belt of resources between two machines.
type Edge struct {
	Buffer         *utils.Buffer
	ItemsProcessed int
	TimeOnBelt     time.Duration
	fromPort       *Port
	fromConn       chan *item.Item
	toPort         *Port
	toConn         chan *item.Item
}

func NewEdge(fromPort *Port, toPort *Port, speed int) *Edge {
	fromConn := make(chan *item.Item, 100)
	toConn := make(chan *item.Item, 100)

	fromPort.Connection = fromConn
	toPort.Connection = toConn

	time := time.Duration(1000/float64(speed)/60.0) * time.Millisecond

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

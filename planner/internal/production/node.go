package production

import (
	"github.com/turbo-pioneer/planner/internal/models"
)

type Node interface {
	ConnectInput(belt chan *models.Item, itemName string) error
	ConnectOutput(belt chan *models.Item, itemName string) error
	// Start initalizes the Nodes Production routines
	Start() (chan struct{}, error)
	// Produce is a life-cycle function which controls the production of a recipe at that recipes defined rate
	Produce(done chan struct{})
}

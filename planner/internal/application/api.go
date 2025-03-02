package application

import (
	"github.com/turbo-pioneer/planner/internal"
)

type Application struct {
	dataRegistry *internal.DataRegistry
	productionLine any
}
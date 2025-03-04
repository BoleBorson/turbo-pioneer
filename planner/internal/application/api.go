package application

import (
	"os"

	"github.com/turbo-pioneer/planner/internal"
	"github.com/turbo-pioneer/planner/internal/production"
)

type Application struct {
	dataRegistry *internal.DataRegistry
	lineBuilder  *production.LineBuilder
}

func NewApplication() (*Application, error) {
	b, err := os.ReadFile("/home/cole/code-projects/turbo-pioneer/data/data1.0.json")
	if err != nil {
		return nil, err
	}
	dr := internal.NewRegistry()
	r, err := dr.LoadRegistryFromFile(b)
	if err != nil {
		return nil, err
	}

	return &Application{
		dataRegistry: r,
		lineBuilder:  production.NewLineBuilder(r),
	}, nil
}

func (a *Application) GenerateLine(recipeName string, rate float64) (*production.ProductionLine, error) {
	prod, err := a.lineBuilder.CreateProductionLineFromRecipe(recipeName, rate)
	if err != nil {
		return nil, err
	}
	return prod, nil
}

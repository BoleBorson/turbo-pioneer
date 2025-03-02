package application

import (
	"os"

	"github.com/turbo-pioneer/planner/internal"
	"github.com/turbo-pioneer/planner/internal/production"
)

type Application struct {
	dataRegistry *internal.DataRegistry
	productionLine any
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
	}, nil
}

func (a *Application) GenerateNode(recipeName string) (*production.Node, error) {
	r, err := a.dataRegistry.GetRecipe(recipeName)
	if err != nil {
		return nil, err
	}
	n := production.NewNode()
	n.Recipe = r
	
	var inputs = make([]*production.Resource, len(r.Ingredients))
	for idx, v := range r.Ingredients {
		i, err := a.dataRegistry.GetItem(v.Item)
		if err != nil {
			return nil, err
		}
		inputs[idx] = production.NewResource(i, v.Amount, r.Time)
	}

	n.Inputs = inputs

	var outputs = make([]*production.Resource, len(r.Products))
	for idx, v := range r.Products {
		i, err := a.dataRegistry.GetItem(v.Item)
		if err != nil {
			return nil, err
		}
		outputs[idx] = production.NewResource(i, v.Amount, r.Time)
	}

	n.Outputs = outputs

	return n, nil
}
package internal

import (
	"github.com/turbo-pioneer/planner/internal/building"
	"github.com/turbo-pioneer/planner/internal/item"
	"github.com/turbo-pioneer/planner/internal/recipe"
	reg "github.com/turbo-pioneer/planner/internal/registry"
)

type DataRegistry struct {
	recipes   reg.Registry
	items     reg.Registry
	buildings reg.Registry
}

func NewRegistry() *DataRegistry {
	return &DataRegistry{}
}

func (dr *DataRegistry) LoadRegistryFromFile(b []byte) (*DataRegistry, error) {
	rreg := recipe.NewRegistry()
	ireg := item.NewRegistry()
	breg := building.NewRegistry()
	if r, err := rreg.LoadRegistryFromFile(b); err == nil {
		dr.recipes = r
	} else {
		return nil, err
	}
	if i, err := ireg.LoadRegistryFromFile(b); err == nil {
		dr.items = i
	} else {
		return nil, err
	}
	if b, err := breg.LoadRegistryFromFile(b); err == nil {
		dr.buildings = b
	} else {
		return nil, err
	}
	return dr, nil
}

func (dr *DataRegistry) GetRecipe(s string) (*recipe.Recipe, error) {
	if r, err := dr.recipes.Get(s); err == nil {
		v := r.(*recipe.Recipe)
		return v, nil
	} else {
		return nil, err
	}
}

func (dr *DataRegistry) GetItem(s string) (*item.Item, error) {
	if r, err := dr.items.Get(s); err == nil {
		v := r.(*item.Item)
		return v, nil
	} else {
		return nil, err
	}
}

func (dr *DataRegistry) GetBuilding(s string) (*building.Building, error) {
	if r, err := dr.buildings.Get(s); err == nil {
		v := r.(*building.Building)
		return v, nil
	} else {
		return nil, err
	}
}

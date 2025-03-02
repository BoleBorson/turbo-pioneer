package internal

import (
	"github.com/turbo-pioneer/planner/internal/item"
	"github.com/turbo-pioneer/planner/internal/recipe"
)

type DataRegistry struct {
	recipes Registry
	items Registry
}

func NewRegistry() *DataRegistry {
	return &DataRegistry{}
}

func LoadRegistryFromFile(b []byte) *Registry  {
	dr := NewRegistry()
	rreg := recipe.NewRegistry()
	ireg := item.NewRegistry()
	if r, err := rreg.LoadRegistryFromFile(b); err != nil {
		dr.recipes = r
	}
	if i, err := ireg.LoadRegistryFromFile(b); err != nil {
		dr.items = i
	}
	
}
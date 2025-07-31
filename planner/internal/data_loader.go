package internal

import (
	"errors"

	"github.com/turbo-pioneer/planner/internal/models"
	reg "github.com/turbo-pioneer/planner/internal/registry"
)

type DataRegistry struct {
	recipes   reg.Registry
	items     reg.Registry
	buildings reg.Registry
	miners    reg.Registry
}

func NewRegistry() *DataRegistry {
	return &DataRegistry{}
}

func (dr *DataRegistry) LoadRegistryFromFile(b []byte) (*DataRegistry, error) {
	rreg := models.NewRecipeRegistry()
	ireg := models.NewItemRegistry()
	breg := models.NewBuildingRegistry()
	mreg := models.NewMinerRegistry()
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
	if m, err := mreg.LoadRegistryFromFile(b); err == nil {
		dr.miners = m
	} else {
		return nil, err
	}
	return dr, nil
}

func (dr *DataRegistry) AllRecipes() (map[string]*models.Recipe, error) {
	recipes := dr.recipes.All()
	if v, ok := recipes.(map[string]*models.Recipe); ok {
		return v, nil
	}
	return nil, errors.New("registry is wrong type")
}

func (dr *DataRegistry) GetRecipe(s string) (*models.Recipe, error) {
	if r, err := dr.recipes.Get(s); err == nil {
		v := r.(*models.Recipe)
		return v, nil
	} else {
		return nil, err
	}
}

func (dr *DataRegistry) GetItem(s string) (*models.Item, error) {
	if r, err := dr.items.Get(s); err == nil {
		v := r.(*models.Item)
		return v, nil
	} else {
		return nil, err
	}
}

func (dr *DataRegistry) GetBuilding(s string) (*models.Building, error) {
	if r, err := dr.buildings.Get(s); err == nil {
		v := r.(*models.Building)
		return v, nil
	} else {
		return nil, err
	}
}

func (dr *DataRegistry) GetMiner(s string) (*models.Miner, error) {
	if r, err := dr.miners.Get(s); err == nil {
		v := r.(*models.Miner)
		return v, nil
	} else {
		return nil, err
	}
}

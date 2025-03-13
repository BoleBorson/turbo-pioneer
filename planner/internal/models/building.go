package models

import (
	"encoding/json"
	"fmt"

	reg "github.com/turbo-pioneer/planner/internal/registry"
)

type BuildingRegistry struct {
	Buildings map[string]*Building `json:"buildings,omitempty"`
}

func NewBuildingRegistry() *BuildingRegistry {
	return &BuildingRegistry{}
}

func (reg *BuildingRegistry) LoadRegistryFromFile(b []byte) (reg.Registry, error) {
	err := json.Unmarshal(b, &reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func (reg *BuildingRegistry) Get(s string) (any, error) {
	if v, ok := reg.Buildings[s]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("building %s, was not found in the registry", s)
	}
}

type Building struct {
	Slug      string    `json:"slug,omitempty"`
	Name      string    `json:"name,omitempty"`
	ClassName string    `json:"className,omitempty"`
	MetaData  *MetaData `json:"metadata,omitempty"`
}

type MetaData struct {
	PowerConsumption         float64 `json:"powerConsumption,omitempty"`
	PowerConsumptionExponent float64 `json:"powerConsumptionExponent,omitempty"`
	ManufacturingSpeed       float64 `json:"manufacturingSpeed,omitempty"`
}

func NewBuilding() *Building {
	return &Building{}
}

func NewMetaData() *MetaData {
	return &MetaData{}
}

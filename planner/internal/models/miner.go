package models

import (
	"encoding/json"
	"fmt"

	reg "github.com/turbo-pioneer/planner/internal/registry"
)

type MinerRegistry struct {
	Items map[string]*Item `json:"items,omitempty"`
}

func NewMinerRegistry() *MinerRegistry {
	return &MinerRegistry{}
}

func (reg *MinerRegistry) LoadRegistryFromFile(b []byte) (reg.Registry, error) {
	err := json.Unmarshal(b, &reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func (reg *MinerRegistry) Get(s string) (any, error) {
	if v, ok := reg.Items[s]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("item %s, was not found in the registry", s)
	}
}

type Miner struct {
	ClassName        string   `json:"class_name,omitempty"`
	AllowedResources []string `json:"allowed_resources,omitempty"`
	AllowLiquids     bool     `json:"allow_liquids,omitempty"`
	AllowSolids      bool     `json:"allow_solids,omitempty"`
	ItemsPerCycle    int      `json:"items_per_cycle,omitempty"`
	ExtractCycleTime float64  `json:"extract_cycle_time,omitempty"`
}

func NewMiner() *Miner {
	return &Miner{}
}

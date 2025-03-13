package models

import (
	"encoding/json"
	"fmt"

	reg "github.com/turbo-pioneer/planner/internal/registry"
)

type MinerRegistry struct {
	Miners map[string]*Miner `json:"miners,omitempty"`
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
	if v, ok := reg.Miners[s]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("miner %s, was not found in the registry", s)
	}
}

type Miner struct {
	ClassName        string   `json:"className,omitempty"`
	AllowedResources []string `json:"allowedResources,omitempty"`
	AllowLiquids     bool     `json:"allowLiquids,omitempty"`
	AllowSolids      bool     `json:"allowSolids,omitempty"`
	ItemsPerCycle    int      `json:"itemsPerCycle,omitempty"`
	ExtractCycleTime float64  `json:"extractCycleTime,omitempty"`
}

func NewMiner() *Miner {
	return &Miner{}
}

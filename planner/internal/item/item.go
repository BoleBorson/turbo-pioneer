package item

import (
	"encoding/json"
	"fmt"

	reg "github.com/turbo-pioneer/planner/internal/registry"
)

type ItemRegistry struct {
	Items map[string]*Item `json:"items,omitempty"`
}

func NewRegistry() *ItemRegistry {
	return &ItemRegistry{}
}

func (reg *ItemRegistry) LoadRegistryFromFile(b []byte) (reg.Registry, error) {
	err := json.Unmarshal(b, &reg)
	if err != nil {
		return nil, err
	}
	return reg, nil
}

func (reg *ItemRegistry) Get(s string) (any, error) {
	if v, ok := reg.Items[s]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("item %s, was not found in the registry", s)
	}
}

type Item struct {
	Slug        string `json:"slug,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ClassName   string `json:"className,omitempty"`
	Liquid      bool   `json:"liquid,omitempty"`
	StackSize   int    `json:"stackSize"`
}

func NewItem() *Item {
	return &Item{}
}

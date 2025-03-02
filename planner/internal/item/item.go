package item

import "encoding/json"

type ItemRegistry struct {
	Items map[string]*Item `json:"items,omitempty"`
}

func NewRegistry() *ItemRegistry {
	return &ItemRegistry{}
}

func (reg *ItemRegistry) LoadRegistryFromFile (b []byte) (*ItemRegistry, error) {
	r := NewRegistry()
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type Item struct {
	Slug        string `json:"slug,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ClassName   string `json:"className,omitempty"`
	Liquid      bool   `json:"liquid,omitempty"`
}

func NewItem() *Item {
	return &Item{}
}
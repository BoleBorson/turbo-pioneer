package recipe

import (
	"encoding/json"
)

type RecipeRegistry struct {
	Recipes map[string]*Recipe `json:"recipes,omitempty"`
}

func NewRegistry() *RecipeRegistry {
	return &RecipeRegistry{}
}

func (reg *RecipeRegistry) LoadRegistryFromFile(b []byte) (*RecipeRegistry, error) {
	r := NewRegistry()
	err := json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

type Recipe struct {
	Slug        string        `json:"slug,omitempty"`
	Name        string        `json:"name,omitempty"`
	ClassName   string        `json:"className,omitempty"`
	Time        int           `json:"time,omitempty"`
	InMachine   bool          `json:"inMachine,omitempty"`
	Ingredients []*Ingredient `json:"ingredients,omitempty"`
	Products    []*Product    `json:"products,omitempty"`
	ProducedIn  []string      `json:"produced_in,omitempty"`
}

type Ingredient struct {
	Item   string  `json:"item,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

type Product struct {
	Item   string  `json:"item,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

func NewRecipe() *Recipe {
	return &Recipe{}
}

func NewIngredient() *Ingredient {
	return &Ingredient{}
}

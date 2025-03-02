package main

import (
	"fmt"
	"os"

	"github.com/turbo-pioneer/planner/internal"
)

func main() {
	b, err := os.ReadFile("/home/cole/code-projects/turbo-pioneer/data/data1.0.json")
	if err != nil {
		panic(err)
	}
	r, err := internal.LoadRegistryFromFile(b)
	if err != nil {
		panic(err)
	}
	
	fmt.Println(r.Recipes["Recipe_Alternate_PolymerResin_C"])
}
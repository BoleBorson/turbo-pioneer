package main

import (
	"fmt"

	"github.com/turbo-pioneer/planner/internal/application"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		panic(err)
	}
	prod, err := app.GenerateLine("Reinforced Iron Plate")
	if err != nil {
		panic(err)
	}

	for _, v := range prod.GetNodes() {
		if v.Root {
			fmt.Println("Root Node")
		}
		fmt.Println(v.Recipe.Name)
	}

	for _, v := range prod.GetEdges() {
		v.PrintEdge()
	}
}

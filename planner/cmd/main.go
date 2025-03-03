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
	// prod, err := app.GenerateLine("Reinforced Iron Plate")
	prod, err := app.GenerateLine("Motor")

	if err != nil {
		panic(err)
	}

	for _, v := range prod.GetNodes() {
		if v.Root {
			fmt.Println("Root Node")
		}
		fmt.Printf("%s, produced in %s \n", v.Recipe.Name, v.Machine.Name)
	}

	for _, v := range prod.GetEdges() {
		v.PrintEdge()
	}
}

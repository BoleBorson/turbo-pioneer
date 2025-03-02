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
	n, err := app.GenerateNode("Recipe_IronRod_C")
	if err != nil {
		panic(err)
	}
	fmt.Println(n.Inputs[0].Rate)
	fmt.Println(n.Outputs[0].Rate)
}
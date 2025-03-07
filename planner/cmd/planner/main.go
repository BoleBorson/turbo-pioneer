package main

import (
	"github.com/turbo-pioneer/planner/internal/application"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		panic(err)
	}
	// prod, err := app.GenerateLine("Reinforced Iron Plate")
	prod, err := app.GenerateLine("Motor", 5)

	if err != nil {
		panic(err)
	}

	prod.Print()
	prod.PrintRequiredRates()
	prod.PrintExcess()

	// nm := prod.CreateNodeMap()
	// for key, value := range nm {
	// 	var totalRate float64
	// 	for _, n := range value {
	// 		totalRate += n.Outputs[0].Rate
	// 	}
	// 	totalRate = math.Round(totalRate*100) / 100
	// 	fmt.Printf("%s, in %d %s's at a total rate of %.2f per min", key, len(value), value[0].Machine.Name, totalRate)
	// 	fmt.Println()
	// }

	// for _, v := range prod.GetNodes() {
	// 	if v.Root {
	// 		fmt.Println("Root Node")
	// 	}
	// 	fmt.Printf("%s, produced in %s \n", v.Recipe.Name, v.Machine.Name)
	// }

	// for _, v := range prod.GetEdges() {
	// 	v.PrintEdge()
	// }
}

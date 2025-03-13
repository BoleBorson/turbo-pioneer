package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/turbo-pioneer/planner/internal/application"
	"github.com/turbo-pioneer/planner/internal/models"
	"github.com/turbo-pioneer/planner/internal/production/nodes"
	"github.com/turbo-pioneer/planner/internal/utils"
)

func main() {
	app, err := application.NewApplication()
	if err != nil {
		panic(err)
	}
	// // prod, err := app.GenerateLine("Reinforced Iron Plate")
	// prod, err := app.GenerateLine("Motor", 5)

	// if err != nil {
	// 	panic(err)
	// }

	// prod.Print()
	// prod.PrintRequiredRates()
	// prod.PrintExcess()
	r, err := app.GetRecipe("Iron Rod")
	if err != nil {
		panic(err)
	}
	rate := utils.Rate(r.Products[0].Amount, r.Time)
	var node = nodes.NewMachineNode(r, rate, &models.Building{})

	in := make(chan *models.Item, 100)
	out := make(chan *models.Item, 100)

	inItem, err := app.GetItem("Desc_IronIngot_C")
	if err != nil {
		panic(err)
	}
	outItem, err := app.GetItem("Desc_IronRod_C")
	if err != nil {
		panic(err)
	}
	for range 5 {
		in <- inItem
	}

	err = node.ConnectInput(in, inItem.ClassName)
	if err != nil {
		panic(err)
	}
	err = node.ConnectOutput(out, outItem.ClassName)
	if err != nil {
		panic(err)
	}
	done, err := node.Start()
	if err != nil {
		panic(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	// Wait for SIGINT (Ctrl+C)
	<-signalChannel
	fmt.Println("Received SIGINT, shutting down...")

	// Close the 'done' channel to signal goroutines to stop
	close(done)
}

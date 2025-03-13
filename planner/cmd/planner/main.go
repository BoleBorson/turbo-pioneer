package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/turbo-pioneer/planner/internal/application"
	"github.com/turbo-pioneer/planner/internal/models"
	"github.com/turbo-pioneer/planner/internal/production"
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

	//  Set up Iron Ore Miner --------------------------------

	ore, err := app.GetItem("Desc_OreIron_C")
	if err != nil {
		panic(err)
	}

	miner, err := app.GetMiner("Desc_MinerMk1_C")
	if err != nil {
		panic(err)
	}

	oreOut := make(chan *models.Item, 100)

	minerNode := nodes.NewMiner(nodes.NewResourceNode(ore, nodes.Normal), miner, 1.0)
	minerOutPort, err := minerNode.ConnectOutput(oreOut, ore.ClassName)
	if err != nil {
		panic(err)
	}

	// Set up Iron Ingot Machine ------------------------------
	r, err := app.GetRecipe("Iron Ingot")
	if err != nil {
		panic(err)
	}
	rate := utils.Rate(r.Products[0].Amount, r.Time)
	var node = nodes.NewMachineNode(r, rate, &models.Building{})

	inOre := make(chan *models.Item, 100)
	outIngot := make(chan *models.Item, 100)

	inItem, err := app.GetItem("Desc_OreIron_C")
	if err != nil {
		panic(err)
	}
	outItem, err := app.GetItem("Desc_IronIngot_C")
	if err != nil {
		panic(err)
	}

	machineInPort, err := node.ConnectInput(inOre, inItem.ClassName)
	if err != nil {
		panic(err)
	}
	err = node.ConnectOutput(outIngot, outItem.ClassName)
	if err != nil {
		panic(err)
	}

	// Set up belt

	belt := production.NewEdge(minerOutPort, machineInPort, production.MK1)

	done1, err := minerNode.Start()
	if err != nil {
		panic(err)
	}

	done2 := belt.Start()
	if err != nil {
		panic(err)
	}

	done3, err := node.Start()
	if err != nil {
		panic(err)
	}

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT)

	// Wait for SIGINT (Ctrl+C)
	<-signalChannel
	fmt.Println("Received SIGINT, shutting down...")

	// Close the 'done' channel to signal goroutines to stop
	close(done1)
	close(done2)
	close(done3)

}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package recipe

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/turbo-pioneer/planner/internal/application"
	"github.com/turbo-pioneer/planner/internal/utils"
)

// recipeCmd represents the recipe command
var RecipeCmd = &cobra.Command{
	Use:   "recipe",
	Short: "recipe provides access to the various parts of Satisfactory recipe",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		recipeName := args[0]
		app, err := application.NewApplication()
		if err != nil {
			fmt.Printf("error starting application, error: %e", err)
			return
		}
		r, err := app.GetRecipe(recipeName)
		if err != nil {
			fmt.Println("recipe not found")
		}
		building, err := app.GetBuilding(r.ProducedIn[0])
		if err != nil {
			fmt.Printf("building %s not found", r.ProducedIn[0])
			return
		}
		fmt.Println("--- " + r.Name + " Recipe ---")
		fmt.Printf("Produced in: %s\n\n", building.Name)
		fmt.Println("Ingredients:")
		for _, i := range r.Ingredients {
			item, err := app.GetItem(i.Item)
			if err != nil {
				fmt.Printf("item %s not found", i.Item)
				return
			}
			rate := utils.Rate(i.Amount, r.Time)
			fmt.Printf("  %.2f %s / min\n", rate, item.Name)
		}
		fmt.Println("\nProducts:")
		for _, p := range r.Products {
			product, err := app.GetItem(p.Item)
			if err != nil {
				fmt.Printf("item %s not found", p.Item)
				return
			}
			rate := utils.Rate(p.Amount, r.Time)
			fmt.Printf("  %.2f %s / min\n", rate, product.Name)
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// recipeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// recipeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

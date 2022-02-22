/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bonusly/utils"
	"fmt"
	"github.com/spf13/cobra"
)

var dryRun bool

// makeitrainCmd represents the makeitrain command
var makeitrainCmd = &cobra.Command{
	Use:   "makeitrain",
	Short: "Award all your remaining bonusly.",
	Long: `Use this command to split your remaining bonusly evenly between the specified recipients.
    Has different modes that decide what happens when you can't split the amount of remaining bonuslys evenly.`,
	Run: func(cmd *cobra.Command, args []string) {
		balance := utils.FetchCurrentGivingBalance()
		if balance < 1 {
			fmt.Println("not enough bonuslys!")
			return
		}
		count := len(recipients)
		canSplitEvenly := balance%count == 0

		amountPerPerson := (balance / count)
		executeTransaction(amountPerPerson, tags, recipients, message)
		if !canSplitEvenly {
			remainingBalance := balance - amountPerPerson*count
			if verbose {
				fmt.Printf("%d bonuslys remaining\n", remainingBalance)
				fmt.Println("Can't split evenly. Continuing with specified mode.")
			}
		}
		outputMessage := "send %d (%d per person) bonusly to %v with message \"%s\"."
		if dryRun {
			fmt.Printf("Would "+outputMessage+"\n", amountPerPerson*count, amountPerPerson, recipients, message)
			return
		}
		if !force {
			fmt.Printf("Will "+outputMessage+" Are you sure? (y/n): ", amountPerPerson*count, amountPerPerson, recipients, message)
			input := ""
			for input != "y" {
				fmt.Scanf("%s", &input)
				if input == "n" {
					fmt.Println("Operation aborted. No bonuslys were transferred.")
					return
				} else if input != "y" {
					fmt.Println("Please enter a valid value. Accepted answers are 'y' or 'n'.")
				}
			}
		}
		executeTransaction(amountPerPerson, tags, recipients, message)
	},
}

func checkMode(mode string) bool {
	// modes := []string{"complete", "oneMessage"}
	return true
}

var mode string

func init() {
	awardCmd.AddCommand(makeitrainCmd)
	makeitrainCmd.Flags().StringVarP(&mode, "mode", "m", "complete", "Choose the algorithm that will be used to distribute the remaining bonuslys")
	makeitrainCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "If this flag is set, the command will not actually execute the transaction, but only show what it would do.")

}

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
				fmt.Println("can't split evenly. Continuing with specified mode.")
			}
		}
		if dryRun {
			fmt.Printf("Would send %d (%d per person) bonusly to %v with message \"%s\".\n",
				amountPerPerson*count, amountPerPerson, recipients, message)
			return
		}

	},
}

func checkMode(mode string) bool {
	// modes := []string{"complete", "oneMessage"}
	return true
}

var mode string

func init() {
	rootCmd.AddCommand(makeitrainCmd)
	makeitrainCmd.Flags().StringSliceVarP(&tags, "hashtags", "g", nil, "Specify optional hashtags that go along with your message.")
	makeitrainCmd.Flags().StringSliceVarP(&recipients, "recipients", "r", nil, "Specify one or more recipients for the bonus.")
	makeitrainCmd.Flags().StringVarP(&mode, "mode", "m", "complete", "Choose the algorithm that will be used to distribute the remaining bonuslys")
	makeitrainCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "If this flag is set, the command will not actually execute the transaction, but only show what it would do.")

}

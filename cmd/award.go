package cmd

import (
	"bonusly/utils"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	recipients  []string
	tags        []string
	message     string
	amount      int
	confirm     bool
	interactive bool
)

// awardCmd represents the award command
var awardCmd = &cobra.Command{
	Use:   "award",
	Short: "Give away bonuses to your coworkers!",
	Long: `This command allows the specification of awards to coworkers.
You may specify the ammount to be given, the destination recipients 
(who will receive) and the message and hashtags sent along.`,
	Run: func(cmd *cobra.Command, args []string) {
		total := len(recipients) * amount

		if validateFlags(amount, tags, recipients, message) {
			return
		}
		// TODO: check if balance is enough
		// TODO: display final balance BEFORE transaction

		fmt.Printf("You will send %d bonusly to %v with message '%v' and hashtags %v.\n", amount, recipients, message, tags)
		if len(recipients) < 1 {
			fmt.Printf("That is a total of %d.\n", total)
		}
		bonusId, err := executeTransaction(amount, tags, recipients, message)
		if err != nil {
			fmt.Println("something went wrong during bonus awarding")
			fmt.Println(err)
		}
		fmt.Printf("Created bonus successfully! Check it out at bonus.ly/bonuses/%s", bonusId)
	},
}

func init() {
	rootCmd.AddCommand(awardCmd)
	awardCmd.PersistentFlags().StringVarP(&message, "message", "m", "You're awesome!", `Your thank you/appreciation message that will be visible
  to everone.`)
	awardCmd.MarkPersistentFlagRequired("message")
	awardCmd.PersistentFlags().StringSliceVarP(&tags, "hashtags", "g", nil, "Specify optional hashtags that go along with your message.")
	awardCmd.MarkPersistentFlagRequired("hashtags")
	awardCmd.PersistentFlags().StringSliceVarP(&recipients, "recipients", "r", nil, "Specify one or more recipients for the bonus.")
	awardCmd.MarkPersistentFlagRequired("recipients")
	awardCmd.Flags().IntVarP(&amount, "amount", "a", 0, "How many bonuslys you want to award.")
	awardCmd.MarkFlagRequired("amount")
	awardCmd.Flags().BoolVarP(&confirm, "confirm", "c", false, "Asks again before sending request.")
	awardCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, `Enables interactive mode. This allows you to assemble
  your message, recipients, etc. step by step.`)
}
func validateFlags(amount int, tags, recipients []string, message string) bool {
	if len(recipients) == 0 {
		fmt.Println("Need to specify at least one recipient! Use '--recipients' or '-r' flag.")
		return false
	}
	if amount < 1 {
		fmt.Println("Please specify at least 1 bonusly to award with '--amount' or '-a' flag.")
		return false
	}
	return true
}

func executeTransaction(amount int, tags, recipients []string, message string) (string, error) {
	for i := 0; i < len(recipients); i++ {
		recipients[i] = "@" + strings.TrimSpace(recipients[i])
	}
	for i := 0; i < len(tags); i++ {
		tags[i] = "#" + strings.TrimSpace(tags[i])
	}
	payload := utils.Bonus{
		Amount: amount,
		Reason: strings.Join(recipients, " ") + " " + message + " " + strings.Join(tags, " "),
	}
	return utils.CreateBonus(payload)
}

package piscator

import (
	"fmt"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
)

var isOrgBool, isPrivateBool bool

var collectCmd = &cobra.Command{
	Use:     "collect",
	Aliases: []string{"c"},
	Short:   "generate URL to collect repos",
	Long:    `generate github url to collect repos from users or an org`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a GitHub name")
			return
		}

		name := args[0]
		isOrgBool, _ = cmd.PersistentFlags().GetBool("org")
		isPrivateBool, _ = cmd.PersistentFlags().GetBool("private")
		res := piscator.GetRepos(name, isOrgBool, isPrivateBool)
		fmt.Println(res)
	},
}

func init() {
	collectCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	collectCmd.PersistentFlags().BoolVarP(&isPrivateBool, "private", "p", false, "Include private repositories")

	rootCmd.AddCommand(collectCmd)
}

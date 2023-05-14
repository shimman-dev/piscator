package piscator

import (
	"fmt"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
)

var isVerbose bool

var netCmd = &cobra.Command{
	Use:     "net",
	Aliases: []string{"c"},
	Short:   "git clone collected repos",
	Long:    `create a directory based on user/org name then create a git repo for each collection`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a GitHub name")
			return
		}

		name := args[0]
		isOrgBool, _ = cmd.PersistentFlags().GetBool("org")
		isPrivateBool, _ = cmd.PersistentFlags().GetBool("private")
		isForkedBool, _ = cmd.PersistentFlags().GetBool("forked")
		makeFileBool, _ = cmd.PersistentFlags().GetBool("makeFile")
		res, err := piscator.GetRepos(name, isOrgBool, isPrivateBool, isForkedBool, makeFileBool)
		if err != nil {
			fmt.Printf("Errors: %s", err)
			return
		}
		isVerbose, _ = cmd.PersistentFlags().GetBool("verbose")
		piscator.CloneReposFromJson(res, name, 5, isVerbose)
		fmt.Println("success friend :)")
	},
}

func init() {
	netCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	netCmd.PersistentFlags().BoolVarP(&isPrivateBool, "private", "p", false, "Include private repositories")
	netCmd.PersistentFlags().BoolVarP(&isForkedBool, "forked", "x", false, "Include forked repositories")
	netCmd.PersistentFlags().BoolVarP(&makeFileBool, "makeFile", "f", false, "Generate a repos.json file")
	netCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "logs detailed messaging to stdout")

	rootCmd.AddCommand(netCmd)
}

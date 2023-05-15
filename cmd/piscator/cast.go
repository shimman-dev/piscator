package piscator

import (
	"fmt"
	"net/http"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
)

var isOrgBool, isPrivateBool, isForkedBool, makeFileBool bool

var castCmd = &cobra.Command{
	Use:     "cast",
	Aliases: []string{"c"},
	Short:   "generate a json struct of github repos",
	Long:    `cast a net into the github sea capturing the URLs of repos from an user/org`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a GitHub username")
			return
		}

		name := args[0]
		isOrgBool, _ = cmd.PersistentFlags().GetBool("org")
		isPrivateBool, _ = cmd.PersistentFlags().GetBool("private")
		isForkedBool, _ = cmd.PersistentFlags().GetBool("forked")
		makeFileBool, _ = cmd.PersistentFlags().GetBool("makeFile")
		res, err := piscator.GetRepos(http.DefaultClient, name, isOrgBool, isPrivateBool, isForkedBool, makeFileBool)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}

func init() {
	castCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	castCmd.PersistentFlags().BoolVarP(&isPrivateBool, "private", "p", false, "Include private repositories")
	castCmd.PersistentFlags().BoolVarP(&isForkedBool, "forked", "x", false, "Include forked repositories")
	castCmd.PersistentFlags().BoolVarP(&makeFileBool, "makeFile", "f", false, "Generate a repos.json file")

	rootCmd.AddCommand(castCmd)
}

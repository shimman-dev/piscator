package piscator

import (
	"fmt"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
)

// var isOrgBool, isPrivateBool bool

var tetherCmd = &cobra.Command{
	Use:     "tether",
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
		res := piscator.GetRepos(name, isOrgBool, isPrivateBool)
		fmt.Println(res)
		piscator.CloneReposFromJson(res, name)
		fmt.Println("success friend :)")
	},
}

func init() {
	tetherCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	tetherCmd.PersistentFlags().BoolVarP(&isPrivateBool, "private", "p", false, "Include private repositories")

	rootCmd.AddCommand(tetherCmd)
}

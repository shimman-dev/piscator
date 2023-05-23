package piscator

import (
	"fmt"
	"net/http"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isSelfBool, isOrgBool, isForkedBool, makeFileBool bool
var githubToken string

var castCmd = &cobra.Command{
	Use:     "cast",
	Aliases: []string{"c"},
	Short:   "generate a json struct of GitHub repos",
	Long: `Ahoy, sailor! Prepare to navigate the GitHub sea and hoist the flag of
exploration with the cast command. Cast your net wide and capture the URLs of
repositories belonging to a user or organization, gathering a bountiful
collection of code treasures. Navigate with ease, discovering new horizons, and
charting your course towards software mastery.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a GitHub username")
			return
		}

		name := args[0]
		tokenFileBool := viper.GetString("github_token") // get token from viper
		isSelfBool, _ = cmd.PersistentFlags().GetBool("self")
		isOrgBool, _ = cmd.PersistentFlags().GetBool("org")
		isForkedBool, _ = cmd.PersistentFlags().GetBool("forked")
		makeFileBool, _ = cmd.PersistentFlags().GetBool("makeFile")

		sleeper := &piscator.RealSleeper{}

		res, err := piscator.GetRepos(http.DefaultClient, sleeper, name, tokenFileBool, isSelfBool, isOrgBool, isForkedBool, makeFileBool)
		if err != nil {
			fmt.Printf("Errors: %s", err)
			return
		}
		fmt.Println(res)
	},
}

func init() {
	viper.AutomaticEnv() // automatically use environment variables

	castCmd.PersistentFlags().BoolVarP(&isSelfBool, "self", "s", false, "Your GitHub user, requires a personal access token")
	castCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	castCmd.PersistentFlags().BoolVarP(&isForkedBool, "forked", "x", false, "Include forked repositories")
	castCmd.PersistentFlags().BoolVarP(&makeFileBool, "makeFile", "f", false, "Generate a repos.json file")
	castCmd.PersistentFlags().StringVarP(&githubToken, "token", "t", "", "GitHub personal access token")

	// bind the token flag to a viper key
	viper.BindPFlag("github_token", castCmd.PersistentFlags().Lookup("token"))

	rootCmd.AddCommand(castCmd)
}

package piscator

import (
	"fmt"
	"net/http"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isSelfBool, isOrgBool, isPrivateBool, isForkedBool, makeFileBool bool
var githubToken string

var castCmd = &cobra.Command{
	Use:     "cast",
	Aliases: []string{"c"},
	Short:   "generate a json struct of GitHub repos",
	Long:    `Avast, ye salty fisherman! Prepare to cast your line with the reel command and embark on a daring fishing expedition in the GitHub waters. As you sail through the digital sea, you'll skillfully create a directory that bears the name of the user or organization, and with each catch, you'll reel in a precious git repository. Like a seasoned fisherman, you'll nurture and cultivate these repositories, transforming them into valuable assets for your coding endeavors. Unleash your fishing prowess, reel in those repositories, and embark on a coding voyage like no other.`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a GitHub username")
			return
		}

		name := args[0]
		tokenFileBool := viper.GetString("github_token") // get token from viper
		isSelfBool, _ = cmd.PersistentFlags().GetBool("self")
		isOrgBool, _ = cmd.PersistentFlags().GetBool("org")
		isPrivateBool, _ = cmd.PersistentFlags().GetBool("private")
		isForkedBool, _ = cmd.PersistentFlags().GetBool("forked")
		makeFileBool, _ = cmd.PersistentFlags().GetBool("makeFile")

		sleeper := &piscator.RealSleeper{}

		res, err := piscator.GetRepos(http.DefaultClient, sleeper, name, tokenFileBool, isSelfBool, isOrgBool, isPrivateBool, isForkedBool, makeFileBool)
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
	castCmd.PersistentFlags().BoolVarP(&isPrivateBool, "private", "p", false, "Include private repositories")
	castCmd.PersistentFlags().BoolVarP(&isForkedBool, "forked", "x", false, "Include forked repositories")
	castCmd.PersistentFlags().BoolVarP(&makeFileBool, "makeFile", "f", false, "Generate a repos.json file")
	castCmd.PersistentFlags().StringVarP(&githubToken, "token", "t", "", "GitHub personal access token")

	// bind the token flag to a viper key
	viper.BindPFlag("github_token", castCmd.PersistentFlags().Lookup("token"))

	rootCmd.AddCommand(castCmd)
}

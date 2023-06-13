package piscator

import (
	"fmt"
	"net/http"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isSelfBool, isOrgBool, isForkedBool, makeFileBool bool
var languageFilter, name, githubToken, username, password, enterprise string

func castRun(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a GitHub username")
		return
	}

	name := args[0]
	tokenFileBool := viper.GetString("github_token") // get token from viper
	isSelfBool, _ := cmd.PersistentFlags().GetBool("self")
	isOrgBool, _ := cmd.PersistentFlags().GetBool("org")
	isForkedBool, _ := cmd.PersistentFlags().GetBool("forked")
	makeFileBool, _ := cmd.PersistentFlags().GetBool("makeFile")

	sleeper := &piscator.RealSleeper{}

	res, err := piscator.GetRepos(http.DefaultClient, sleeper, name, tokenFileBool, username, password, enterprise, isSelfBool, isOrgBool, isForkedBool, makeFileBool)

	if err != nil {
		fmt.Printf("Errors: %s", err)
		return
	}

	if languageFilter != "" {
		res, err = piscator.RepoByLanguage(res, languageFilter)
		if err != nil {
			fmt.Printf("Error filtering repositories by language: %s", err)
			return
		}
	}

	fmt.Println(res)
}

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
	Run:  castRun,
}

func init() {
	viper.AutomaticEnv() // automatically use environment variables

	castCmd.PersistentFlags().BoolVarP(&isSelfBool, "self", "s", false, "Your GitHub user, requires a personal access token")
	castCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	castCmd.PersistentFlags().BoolVarP(&isForkedBool, "forked", "x", false, "Include forked repositories")
	castCmd.PersistentFlags().BoolVarP(&makeFileBool, "makeFile", "f", false, "Generate a repos.json file")

	castCmd.PersistentFlags().StringVarP(&languageFilter, "language", "l", "", "Filter repositories by language(s)")

	castCmd.PersistentFlags().StringVarP(&githubToken, "token", "t", "", "GitHub personal access token")
	castCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "GitHub username")
	castCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "GitHub password")
	castCmd.PersistentFlags().StringVarP(&enterprise, "enterprise", "e", "", "GitHub Enterprise URL")

	// bind the token flags to env keys
	viper.BindPFlag("github_token", castCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("username", castCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", castCmd.PersistentFlags().Lookup("password"))

	viper.BindEnv("github_token", "GITHUB_TOKEN")
	viper.BindEnv("username", "GITHUB_USERNAME")
	viper.BindEnv("password", "GITHUB_PASSWORD")

	rootCmd.AddCommand(castCmd)
	castCmd.AddCommand(generateManCmd)
}

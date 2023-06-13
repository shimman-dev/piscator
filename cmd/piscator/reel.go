package piscator

import (
	"fmt"
	"net/http"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isVerbose bool

func reelRun(cmd *cobra.Command, args []string) {
	if isSelfBool {
		// Use the GitHub username associated with the token
		name = "" // Set name to empty string for self
	} else if len(args) < 1 && githubToken == "" && viper.GetString("github_token") == "" {
		fmt.Println("Please provide a GitHub username or specify a token")
		return
	} else if len(args) >= 1 {
		name = args[0]
	}

	if name == "" {
		fmt.Println("Please provide a GitHub username or org name")
		return
	}

	tokenFileBool := viper.GetString("github_token") // get token from viper
	isForkedBool, _ = cmd.PersistentFlags().GetBool("forked")
	makeFileBool, _ = cmd.PersistentFlags().GetBool("makeFile")

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

	concurrentLimit := int8(10)
	isVerbose, _ = cmd.PersistentFlags().GetBool("verbose")

	err = piscator.CloneReposFromJson(piscator.RealCommandExecutor{}, res, name, concurrentLimit, isVerbose)

	if err != nil {
		fmt.Printf("Errors: %s", err)
		return
	}

	fmt.Println("success friend :)")
}

var reelCmd = &cobra.Command{
	Use:     "reel",
	Aliases: []string{"c"},
	Short:   "git clone collected repos",
	Long: `Avast, ye salty fisherman! Prepare to cast your line with the reel command
and embark on a daring fishing expedition in the GitHub waters. As you sail
through the digital sea, you'll skillfully create a directory that bears the
name of the user or organization, and with each catch, you'll reel in a precious
git repository. Like a seasoned fisherman, you'll nurture and cultivate these
repositories, transforming them into valuable assets for your coding endeavors.
Unleash your fishing prowess, reel in those repositories, and embark on a coding
voyage like no other.`,
	Args: cobra.MinimumNArgs(1),
	Run:  reelRun,
}

func init() {
	viper.AutomaticEnv()

	reelCmd.PersistentFlags().BoolVarP(&isSelfBool, "self", "s", false, "Your GitHub user, requires a personal access token")
	reelCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	reelCmd.PersistentFlags().BoolVarP(&isForkedBool, "forked", "x", false, "Include forked repositories")
	reelCmd.PersistentFlags().BoolVarP(&makeFileBool, "makeFile", "f", false, "Generate a repos.json file")
	reelCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "logs detailed messaging to stdout")

	reelCmd.PersistentFlags().StringVarP(&languageFilter, "language", "l", "", "Filter repositories by language(s)")

	reelCmd.PersistentFlags().StringVarP(&githubToken, "token", "t", "", "GitHub personal access token")
	reelCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "GitHub username")
	reelCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "GitHub password")
	reelCmd.PersistentFlags().StringVarP(&enterprise, "enterprise", "e", "", "GitHub Enterprise URL")

	// bind the token flags to env keys
	viper.BindPFlag("github_token", castCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("username", castCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", castCmd.PersistentFlags().Lookup("password"))

	viper.BindEnv("github_token", "GITHUB_TOKEN")
	viper.BindEnv("username", "GITHUB_USERNAME")
	viper.BindEnv("password", "GITHUB_PASSWORD")

	rootCmd.AddCommand(reelCmd)
	reelCmd.AddCommand(generateManCmd)
}

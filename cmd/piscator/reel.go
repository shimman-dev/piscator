package piscator

import (
	"fmt"
	"net/http"

	"github.com/shimman-dev/piscator/pkg/piscator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var isVerbose bool

var reelCmd = &cobra.Command{
	Use:     "reel",
	Aliases: []string{"c"},
	Short:   "git clone collected repos",
	Long:    `Create a directory based on user/org name then create a git repo for each collection`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please provide a GitHub name")
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

		concurrentLimit := int8(10)
		isVerbose, _ = cmd.PersistentFlags().GetBool("verbose")

		piscator.CloneReposFromJson(piscator.RealCommandExecutor{}, res, name, concurrentLimit, isVerbose)

		fmt.Println("success friend :)")
	},
}

func init() {
	viper.AutomaticEnv() // automatically use environment variables

	reelCmd.PersistentFlags().BoolVarP(&isSelfBool, "self", "s", false, "Your GitHub user, requires a personal access token")
	reelCmd.PersistentFlags().BoolVarP(&isOrgBool, "org", "o", false, "Is an organization")
	reelCmd.PersistentFlags().BoolVarP(&isPrivateBool, "private", "p", false, "Include private repositories")
	reelCmd.PersistentFlags().BoolVarP(&isForkedBool, "forked", "x", false, "Include forked repositories")
	reelCmd.PersistentFlags().BoolVarP(&makeFileBool, "makeFile", "f", false, "Generate a repos.json file")
	reelCmd.PersistentFlags().BoolVarP(&isVerbose, "verbose", "v", false, "logs detailed messaging to stdout")
	reelCmd.PersistentFlags().StringVarP(&githubToken, "token", "t", "", "GitHub personal access token")

	// bind the token flag to a viper key
	viper.BindPFlag("github_token", reelCmd.PersistentFlags().Lookup("token"))

	rootCmd.AddCommand(reelCmd)
}

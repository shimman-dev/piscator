package piscator

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// used for flags
var (
	// is name an organization, not a user
	isOrg string
	// is repo private
	isPrivate string
	// is repo a fork
	isForked string
	// include verbose logging
	verbose string
)

var rootCmd = &cobra.Command{
	Use:   "piscator",
	Short: "piscator is a CLI tool for cloning GitHub repositories",
	Long:  `Grab all the repositories from a GitHub organization or user and clone them locally. Visit https://www.piscator.dev for documentation and usage guides.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Is piscaator working?")
	},
}

func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your piscator command '%s'", err)
		os.Exit(1)
	}

	return nil
}

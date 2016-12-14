package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	contest, language    string
	debug, overwriteMain bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   os.Args[0],
		Short: os.Args[0] + " is a hackerrank local helper",
	}

	rootCmd.PersistentFlags().StringVarP(&contest, "contest", "c", "master", "the contest containing the challenge")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "debug the chatter")
	rootCmd.PersistentFlags().BoolVarP(&overwriteMain, "force", "f", false, "allow overwriting main")
	rootCmd.PersistentFlags().StringVarP(&language, "language", "l", "golang", "perform a challenge in a specific language")

	rootCmd.AddCommand(cmdSubmit, getCommand, cmdNext)

	rootCmd.Execute()

}

var getCommand = &cobra.Command{
	Use:   "get [challenge-slug]",
	Short: "download the challenge",
	Run:   runGetCommand,
}

func runGetCommand(cmd *cobra.Command, args []string) {
	challengeSlug := ""
	if len(args) > 0 {
		challengeSlug = args[0]
	} else {
		var err error
		challengeSlug, err = getChallengeNameFromCache()
		if err != nil {
			challengeSlug = dirName()
		}
	}

	if challengeSlug == "." {
		challengeSlug = dirName()
	}

	err := get(challengeSlug)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func dirName() string {
	p, err := filepath.Abs(".")

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not resolve current directory name")
		os.Exit(1)
	}

	return filepath.Base(p)
}

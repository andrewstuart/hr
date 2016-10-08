package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const submitLong = `Submit your challenge.

"submit" will use, in order, the first argument passed, metadata previously
stored in the local ` + cacheFileName + ` file in the current directory, or attempt
to use the current directory name.`

var cmdSubmit = &cobra.Command{
	Use:   "submit [challenge-slug]",
	Short: "submit a challenge",
	Long:  submitLong,
	Run:   runSubmit,
}

func runSubmit(cmd *cobra.Command, args []string) {
	slug := ""
	if len(args) > 0 {
		slug = args[0]
	} else {
		var err error
		slug, err = getChallengeNameFromCache()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Missing challenge name (as first argument). Trying directory name.")
			slug = dirName()
		}
	}

	f, err := os.OpenFile("./main.go", os.O_RDONLY, 0640)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer f.Close()

	err = submit(slug, f)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}

func getChallengeNameFromCache() (string, error) {
	cache, err := os.OpenFile(cacheFileName, os.O_RDONLY, 0400)
	if err != nil {
		return "", err
	}
	defer cache.Close()

	var chal struct{ Model Challenge }
	err = json.NewDecoder(cache).Decode(&chal)
	if err != nil {
		return "", err
	}

	return chal.Model.Name, nil
}

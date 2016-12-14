package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var cmdNext = &cobra.Command{
	Use:   "next",
	Short: "go to the next challenge",
	Long:  "go to the next challenge",
	Run:   next,
}

func next(cmd *cobra.Command, args []string) {
	var s SubmissionStatus

	f, err := os.OpenFile(statusFileName, os.O_RDONLY, 0640)
	if err != nil {
		if err == os.ErrNotExist {
			log.Printf("It looks like you haven't submitted a challenge from this directory. We can't determine which challenge is next.")
			os.Exit(2)
		}
		log.Printf("Error opening previous submission response: %s", err)
		os.Exit(1)
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(&s)
	if err != nil {
		log.Fatal(err)
	}

	dir := "../" + s.NextChallenge.Slug
	os.Mkdir(dir, 0750)
	os.Chdir(dir)

	err = get(s.NextChallenge.Slug)
	if err != nil {
		log.Fatal(err)
	}
}

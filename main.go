package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
	contest       = flag.String("contest", "master", "the contest containing the challenge")
	debug         = flag.Bool("debug", false, "debug the chatter")
	overwriteMain = flag.Bool("-m", false, "allow overwriting main")
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "Missing challenge name (as first argument)")
		os.Exit(1)
	}

	if flag.Args()[0] == "submit" {
		challengeSlug := ""
		if len(flag.Args()) < 2 {
			cache, err := os.OpenFile(cacheFileName, os.O_RDONLY, 0400)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Missing challenge name (as first argument)")
				os.Exit(1)
			}
			defer cache.Close()

			var chal struct{ Model Challenge }
			err = json.NewDecoder(cache).Decode(&chal)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Missing challenge name and could not decode cache file: %s\n", err)
			}
			challengeSlug = chal.Model.Slug
		} else {
			challengeSlug = flag.Args()[1]
		}

		f, err := os.OpenFile("./main.go", os.O_RDONLY, 0640)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()

		err = submit(challengeSlug, f)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	err := get(*contest, flag.Args()[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type hr struct {
	BodyHTML string `json:"body_html"`
}

var (
	contest       = flag.String("contest", "master", "the contest containing the challenge")
	debug         = flag.Bool("debug", false, "debug the chatter")
	overwriteMain = flag.Bool("-m", false, "allow overwriting main")
)

func main() {
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Missing challenge name (as first argument)")
		os.Exit(1)
	}

	if flag.Args()[0] == "submit" {
		if len(flag.Args()) < 2 {
			fmt.Println("Missing challenge name (as first argument)")
			os.Exit(1)
		}

		f, err := os.OpenFile("./main.go", os.O_RDONLY, 0640)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		err = submit(flag.Args()[1], f)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	err := get(*contest, flag.Args()[0])
	if err != nil {
		log.Fatal(err)
	}

}

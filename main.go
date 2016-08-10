package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/puerkitobio/goquery"
)

type hr struct {
	BodyHTML string `json:"body_html"`
}

var (
	contest = flag.String("contest", "master", "the contest containing the challenge")
	debug   = flag.Bool("debug", false, "debug the chatter")
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

	url := fmt.Sprintf("https://www.hackerrank.com/rest/contests/%s/challenges/%s", *contest, flag.Args()[0])

	log.Println(url)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal("challenges", err)
	}

	defer res.Body.Close()

	var h struct{ Model hr }

	err = json.NewDecoder(res.Body).Decode(&h)
	if err != nil {
		log.Println(res.Status, res.Status)
		log.Fatal("decode error ", err)
	}

	d, err := goquery.NewDocumentFromReader(strings.NewReader(h.Model.BodyHTML))
	if err != nil {
		log.Println("goquery err", err)
	}

	in := d.Find(".challenge_sample_input_body pre").Text()
	out := d.Find(".challenge_sample_output_body pre").Text()

	tmpl, err := template.New("test").Delims("{%", "%}").Parse(testTmpl)
	if err != nil {
		log.Fatal("Template compile error ", err)
	}

	f, err := os.OpenFile("main_test.go", os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]string{"in": in, "out": out})
	if err != nil {
		log.Fatal(err)
	}

	f2, err := os.OpenFile("main.go", os.O_CREATE|os.O_RDWR, 0640)
	if err != nil {
		log.Fatal(err)
	}
	defer f2.Close()

	_, err = f2.Write([]byte(mainTmpl))
	if err != nil {
		log.Fatal(err)
	}
}

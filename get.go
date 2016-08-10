package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type example struct {
	In, Out string
}

// Challenge is a representation of a HackerRank challenge
type Challenge struct {
	BodyHTML    template.HTML `json:"body_html"`
	ContestSlug string        `json:"contest_slug"`
	Name        string        `json:"name"`
	Slug        string        `json:"slug"`
	Link        string        `json:"link"`
	Preview     string        `json:"preview"`
}

func get(contest, challenge string) error {
	url := fmt.Sprintf("https://www.hackerrank.com/rest/contests/%s/challenges/%s", contest, challenge)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var h struct{ Model Challenge }

	cachef, err := os.OpenFile(".response.json", os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer cachef.Close()

	err = json.NewDecoder(io.TeeReader(res.Body, cachef)).Decode(&h)
	if err != nil {
		log.Println(res.Status, res.Status)
		return err
	}

	d, err := goquery.NewDocumentFromReader(strings.NewReader(string(h.Model.BodyHTML)))
	if err != nil {
		log.Println("goquery err", err)
	}

	exPres := d.Find("pre")

	examples := []example{}

	var ex example

	for i := 0; i < exPres.Length(); i++ {
		if i%2 == 0 {
			ex = example{In: exPres.Eq(i).Text()}
		} else {
			ex.Out = exPres.Eq(i).Text()
			examples = append(examples, ex)
		}
	}

	tmpl, err := template.New("test").Delims("{%", "%}").Parse(testTmpl)
	if err != nil {
		return err
	}

	filePerms := os.O_TRUNC | os.O_RDWR | os.O_CREATE

	if !*overwriteMain {
		filePerms |= os.O_EXCL
	}

	fchal, err := os.OpenFile("challenge.html", filePerms&^os.O_EXCL, 0640)
	if err != nil {
		return err
	}
	defer fchal.Close()

	doc := template.Must(template.New("contest").Parse(`<html><head></head><body><link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous"><div class="container"><header class="jumbotron"><h1>{{ .Name }}<a href="{{ .Link }}"><span class="glyphicon glyphicon-link"></span></a></h1><h2>{{ .Preview }}</h2></header>{{ .BodyHTML }}<div class="container"></body></html>`))

	contestFrag := ""
	if h.Model.ContestSlug != "master" && h.Model.ContestSlug != "" {
		contestFrag = fmt.Sprintf("contests/%s/", h.Model.ContestSlug)
	}
	h.Model.Link = fmt.Sprintf("https://www.hackerrank.com/%schallenges/%s", contestFrag, h.Model.Slug)

	err = doc.Execute(fchal, h.Model)
	if err != nil {
		return err
	}
	defer func() {
		exec.Command("xdg-open", "challenge.html").Run()
	}()

	f, err := fileErr(os.OpenFile("main_test.go", filePerms, 0640))
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]interface{}{"examples": examples})
	if err != nil {
		return err
	}

	f2, err := fileErr(os.OpenFile("main.go", filePerms, 0640))
	if err != nil {
		return err
	}
	defer f2.Close()

	_, err = f2.Write([]byte(mainTmpl))
	if err != nil {
		return err
	}
	return nil
}

func fileErr(f *os.File, err error) (*os.File, error) {
	if err != nil && strings.Contains(err.Error(), "file exists") {
		err = fmt.Errorf("main.go already exists - force an overwrite with -m option")
	}
	return f, err
}

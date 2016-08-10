package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type example struct {
	In, Out string
}

func get(contest, challenge string) error {
	url := fmt.Sprintf("https://www.hackerrank.com/rest/contests/%s/challenges/%s", contest, challenge)

	log.Println(url)

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var h struct{ Model hr }

	err = json.NewDecoder(res.Body).Decode(&h)
	if err != nil {
		log.Println(res.Status, res.Status)
		return err
	}

	d, err := goquery.NewDocumentFromReader(strings.NewReader(h.Model.BodyHTML))
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

	f, err := os.OpenFile("main_test.go", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, map[string]interface{}{"examples": examples})
	if err != nil {
		return err
	}

	f2, err := os.OpenFile("main.go", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
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

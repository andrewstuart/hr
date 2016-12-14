package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
)

const (
	loginURL   = "https://www.hackerrank.com/auth/login"
	subRestFmt = "https://www.hackerrank.com/rest/contests/%s/challenges/%s/submissions"
	refFmt     = "https://www.hackerrank.com/challenges/%s"

	statusFileName = ".submission-status.json"
)

type submission struct {
	ContestSlug string `json:"contest_slug"`
	Code        string `json:"code"`
	Language    string `json:"language"`
}

func submit(name string, code io.Reader) error {
	if cl, ok := code.(io.Closer); ok {
		defer cl.Close()
	}

	jar := &session{
		RTer: &csrfTransport{Referrer: fmt.Sprintf(refFmt, name)},
	}
	if len(jar.Cookies(nil)) == 0 {
		err := jar.login()
		if err != nil {
			return err
		}
	}

	cli := &http.Client{Jar: jar, Transport: jar.RTer}

	bs, err := ioutil.ReadAll(code)
	if err != nil {
		return err
	}

	jsbs, err := json.Marshal(submission{ContestSlug: contest, Language: "go", Code: string(bs)})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(subRestFmt, contest, name), bytes.NewReader(jsbs))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var r io.Reader = res.Body

	if debug {
		r = io.TeeReader(r, os.Stdout)
	}

	m := struct{ Model SubmissionStatus }{}
	err = json.NewDecoder(r).Decode(&m)
	if err != nil {
		return err
	}

	stat := &m.Model

	if stat.Status == "" {
		return fmt.Errorf("Emtpy status returned")
	}

	for len(stat.TestcaseStatus) == 0 {
		time.Sleep(3300 * time.Millisecond)

		err = stat.Update(cli)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\rCurrent status: %s\t\t", stat.Status)
		for _, status := range stat.LiveStatus.TestcaseMessage {
			switch status {
			case "Success":
				color.New(color.FgGreen).Print("\u2713")
			case "Terminated due to timeout":
				color.New(color.FgYellow).Print("\u25F7")
			case "Wrong Answer":
				color.New(color.FgRed).Print("\u2717")
			default:
				fmt.Print("?")
			}
		}
	}

	fmt.Printf("\rTest results as follows: %s\n", stat.Status)
	for i, s := range stat.TestcaseMessage {
		fmt.Printf("Test Case %d:\t%s\n", i, s)
	}

	fmt.Printf("\nView results at:\t%s\n", stat.UserURL())
	fmt.Printf("Next contest:\t%s -- %s\n", stat.NextChallengeSlug, stat.NextChallenge.Preview)

	f, err := os.OpenFile(statusFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	json.NewEncoder(f).Encode(stat)

	return f.Close()
}

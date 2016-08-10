package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TestCase struct {
	TestcaseStatus  []int    `json:"testcase_status"`
	TestcaseMessage []string `json:"testcase_message"`
	Progress        int      `json:"progress"`
}

type SubmissionStatus struct {
	TestCase
	ID                int    `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	ChallengeID       int    `json:"challenge_id"`
	ContestSlug       string `json:"contest_slug"`
	TestcaseStatus    []int  `json:"testcase_status"`
	Solved            int    `json:"solved"`
	NextChallengeSlug string `json:"next_challenge_slug"`
	NextChallenge     struct {
		Difficulty  string `json:"difficulty_name"`
		URL         string `json:"url"`
		Slug        string `json:"slug"`
		Name        string `json:"name"`
		Preview     string `json:"preview"`
		SolvedCount int    `json:"solved_count"`
		TotalCount  int    `json:"total_count"`
	} `json:"next_challenge"`
	LiveStatus struct {
		TestCase
	} `json:"live_status"`
	Score  string
	Status string `json:"status"`
}

const viewStatFmt = "https://www.hackerrank.com/%schallenges/%s/submissions"

// UserURL returns the url for the user to view the results in the browser
func (ss *SubmissionStatus) UserURL() string {
	var contestFragment string
	if ss.ContestSlug != "master" {
		contestFragment = fmt.Sprintf("contests/%s/", ss.ContestSlug)
	}
	return fmt.Sprintf(viewStatFmt, contestFragment, ss.Slug)
}

const subStatURL = "https://www.hackerrank.com/rest/contests/%s/submissions/%d?_=%d"

func (ss *SubmissionStatus) Update(cli *http.Client) error {
	url := fmt.Sprintf(subStatURL, ss.ContestSlug, ss.ID, time.Now().UnixNano()/1000000)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	m := struct{ Model *SubmissionStatus }{ss}

	return json.NewDecoder(res.Body).Decode(&m)
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
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
	ID             int    `json:"id"`
	ChallengeID    int    `json:"challenge_id"`
	ContestSlug    string `json:"contest_slug"`
	TestcaseStatus []int  `json:"testcase_status"`
	Solved         int    `json:"solved"`
	NextChallenge  struct {
		Difficulty  string `json:"difficulty_name"`
		URL         string `json:"url"`
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

const subStatURL = "https://www.hackerrank.com/rest/contests/%s/submissions/%d?_=%d"

func (ss *SubmissionStatus) Update(cli *http.Client) error {
	url := fmt.Sprintf(subStatURL, ss.ContestSlug, ss.ID, time.Now().UnixNano()/1000000)
	log.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := cli.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	var m struct{ Model SubmissionStatus }

	err = json.NewDecoder(res.Body).Decode(&m)
	if err != nil {
		return err
	}
	*ss = m.Model
	return nil
}

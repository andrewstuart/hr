package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func getToken(body io.ReadCloser) (string, error) {
	var res struct {
		CSRFToken string `json:"csrf_token"`
		Status    bool
		Messages  []string
	}
	err := json.NewDecoder(body).Decode(&res)
	if err != nil {
		return "", err
	}

	if !res.Status {
		return "", fmt.Errorf("Error authenticating: %s", res.Messages)
	}

	return res.CSRFToken, body.Close()
}

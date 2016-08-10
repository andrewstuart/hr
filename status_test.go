package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type resRT struct {
	resp *http.Response
}

func (r *resRT) RoundTrip(*http.Request) (*http.Response, error) {
	return r.resp, nil
}

func testClient(resBody string, res *http.Response) *http.Client {
	if res == nil {
		res = &http.Response{
			Body:       ioutil.NopCloser(strings.NewReader(resBody)),
			StatusCode: 200,
			Status:     "200 OK",
		}
	}

	return &http.Client{
		Transport: &resRT{resp: res},
	}
}

func TestDecode(t *testing.T) {
	f, err := os.OpenFile("./test/status.json", os.O_RDONLY, 0400)
	require.NoError(t, err, "Error opening file")

	// defer require.NoError(t, f.Close(), "Error closing file")

	var s struct{ Model SubmissionStatus }

	require.NoError(t, json.NewDecoder(f).Decode(&s), "Error decoding")

	require.Len(t, s.Model.LiveStatus.TestcaseStatus, 1)
}

func TestUpdate(t *testing.T) {
	f, err := os.OpenFile("./test/status.json", os.O_RDONLY, 0400)
	require.NoError(t, err, "Error opening file")

	// defer require.NoError(t, f.Close(), "Error closing file")

	var s struct{ Model SubmissionStatus }

	require.NoError(t, json.NewDecoder(f).Decode(&s), "Error decoding")

	require.Len(t, s.Model.LiveStatus.TestcaseStatus, 1)

	status := &s.Model

	fup, err := os.OpenFile("./test/status-2.json", os.O_RDONLY, 0400)
	require.NoError(t, err, "Opening update file")

	res := &http.Response{Body: fup, Status: "200 OK", StatusCode: 200}

	status.Update(testClient("", res))

	require.Len(t, status.LiveStatus.TestcaseStatus, 2)
	require.Equal(t, status.LiveStatus.TestcaseMessage[1], "Success")
}

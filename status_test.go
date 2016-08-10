package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	f, err := os.OpenFile("./test/status.json", os.O_RDONLY, 0400)
	require.NoError(t, err, "Error opening file")

	// defer require.NoError(t, f.Close(), "Error closing file")

	var s struct{ Model SubmissionStatus }

	require.NoError(t, json.NewDecoder(f).Decode(&s), "Error decoding")

	require.Len(t, s.Model.LiveStatus.TestcaseStatus, 1)
}

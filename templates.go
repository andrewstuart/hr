package main

var (
	testTmpl = `package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

` + "var tests = [][]string{{`{% .in %}`, `{% .out %}`}}" + `

func TestMain(t *testing.T) {
	for _, test := range tests {
		b := &bytes.Buffer{}
		cmain(strings.NewReader(test[0]), b)

		require.Equal(t, strings.TrimSpace(test[1]), strings.TrimSpace(string(b.Bytes())))
	}
}`

	mainTmpl = `package main

import (
	"io"
	"os"
)

func cmain(r io.Reader, w io.Writer) {
	// implement here
}

func main() {
	cmain(os.Stdin, os.Stdout)
}
`
)

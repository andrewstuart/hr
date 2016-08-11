package main

const (
	testTmpl = `package main

import (
	"bytes"
	"strings"
	"testing"
)

` + "var tests = [][]string{{% range $i, $ex :=  .examples %}{`{% $ex.In %}`, `{% $ex.Out %}`},{% end %}}" + `

func TestMain(t *testing.T) {
	for i, test := range tests {
		b := &bytes.Buffer{}
		cmain(strings.NewReader(test[0]), b)

		expect, got := strings.TrimSpace(test[1]), strings.TrimSpace(string(b.Bytes()))
		if expect != got {
			t.Errorf("Test case %d failed.\nExpected:\t%s\nGot:\t\t\t%s", i+1, expect, got)
		}
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

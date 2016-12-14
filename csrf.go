package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// type rtf func(*http.Request) (*http.Response, error)
type csrfTransport struct {
	Referrer string
	CSRF     string
}

func (r *csrfTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.85 Safari/537.36")
	req.Header.Set("Origin", "https://www.hackerrank.com")
	req.Header.Set("Referrer", r.Referrer)

	if r.CSRF != "" {
		req.Header.Set("X-CSRF-Token", r.CSRF)
	}

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if debug {
		fmt.Println("Req")
		req.Write(os.Stdout)
		fmt.Println(req.URL.String())
		req.Header.Write(os.Stdout)
		// fmt.Println("Req headers")
		// req.Header.Write(os.Stdout)
		fmt.Println()

		fmt.Println("Res")
		res.Write(os.Stdout)
		fmt.Println()
	}

	return res, err
}

func getBodyCSRF(r io.ReadCloser) (string, error) {
	defer r.Close()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		log.Fatal(err)
	}

	csrf, ok := doc.Find(`meta[name="csrf-token"]`).First().Attr("content")
	if !ok {
		return "", fmt.Errorf("No csrf meta")
	}
	return csrf, nil
}

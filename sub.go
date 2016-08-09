package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/puerkitobio/goquery"
)

const (
	loginURL   = "https://www.hackerrank.com/auth/login"
	subRestFmt = "https://www.hackerrank.com/rest/contests/%s/challenges/%s/submissions?"
	subFmt     = "https://www.hackerrank.com/contests/%s/challenges/%s/submissions"
	refFmt     = "https://www.hackerrank.com/challenges/%s"
)

type submission struct {
	ContestSlug string `json:"contest_slug"`
	Code        string `json:"code"`
	Language    string `json:"language"`
}

type session struct {
	cookies map[string]string
	rtrip   *csrfTransport
}

func (c *session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if c.cookies == nil {
		c.cookies = map[string]string{}
	}
	for _, ck := range cookies {
		c.cookies[ck.Name] = ck.Value
	}

	f, err := os.OpenFile("/home/andrew/.local/hr.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0640)
	if err != nil {
		return
	}

	json.NewEncoder(f).Encode(*c)
	f.Close()
}

func (c *session) Cookies(u *url.URL) []*http.Cookie {
	if len(c.cookies) == 0 {
		f, err := os.OpenFile("~/.local/hr.json", os.O_RDONLY, 0640)
		if err == nil {
			json.NewDecoder(f).Decode(&c)
			f.Close()
		}
	}

	cookies := []*http.Cookie{}
	for k, v := range c.cookies {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	return cookies
}

func (c *session) login() error {
	cli := http.Client{Transport: c.rtrip, Jar: c}
	res, err := cli.Get("https://www.hackerrank.com/login")
	if err != nil {
		return err
	}

	csrf, err := getBodyCSRF(res.Body)
	if err != nil {
		return err
	}
	c.rtrip.CSRF = csrf

	v := url.Values{
		"login":       []string{os.Getenv("HR_USER")},
		"password":    []string{os.Getenv("HR_PASS")},
		"fallback":    []string{"true"},
		"remember_me": []string{"false"},
	}

	req, err := http.NewRequest("POST", loginURL, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}

	res, err = cli.Do(req)
	if err != nil {
		return err
	}

	return nil
}

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

	fmt.Println("Req headers")
	req.Header.Write(os.Stdout)
	fmt.Println()

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	fmt.Println("Res headers")
	res.Header.Write(os.Stdout)
	fmt.Println()

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

func submit(name string, code io.Reader) error {
	if cl, ok := code.(io.Closer); ok {
		defer cl.Close()
	}

	jar := &session{
		rtrip: &csrfTransport{Referrer: fmt.Sprintf(refFmt, name)},
	}
	if len(jar.Cookies(nil)) == 0 {
		err := jar.login()
		if err != nil {
			return err
		}
	}

	cli := &http.Client{Jar: jar, Transport: jar.rtrip}

	bs, err := ioutil.ReadAll(code)
	if err != nil {
		return err
	}

	jsbs, err := json.Marshal(submission{ContestSlug: *contest, Language: "go", Code: string(bs)})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf(subRestFmt, *contest, name), bytes.NewReader(jsbs))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := cli.Do(req)
	if err != nil {
		return err
	}

	log.Println(res.Status)
	io.Copy(os.Stdout, res.Body)

	fmt.Println()
	fmt.Printf("View results at: "+subFmt+"\n", *contest, name)

	return nil
}

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

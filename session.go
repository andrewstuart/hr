package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	prompt "github.com/segmentio/go-prompt"
)

type session struct {
	CStore map[string]string
	RTer   http.RoundTripper
}

const hrDir = ".local/hr.json"

var (
	home, hrCachePath string
	cli               http.Client
)

func init() {
	var err error
	home, err = homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}

	hrCachePath = path.Join(home, hrDir)
}

func (c *session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	if c.CStore == nil {
		c.CStore = map[string]string{}
	}
	for _, ck := range cookies {
		c.CStore[ck.Name] = ck.Value
	}

	f, err := os.OpenFile(hrCachePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return
	}

	json.NewEncoder(f).Encode(*c)
	f.Close()
}

//Cookies implements cookie store
func (c *session) Cookies(u *url.URL) []*http.Cookie {
	if c.CStore == nil || len(c.CStore) == 0 {
		f, err := os.OpenFile(hrCachePath, os.O_RDONLY, 0600)
		if err != nil {
			log.Println(err)
		}

		err = json.NewDecoder(f).Decode(&c)
		if err != nil {
			log.Println(err)
		}
		err = f.Close()
		if err != nil {
			log.Println(err)
		}
	}

	cookies := []*http.Cookie{}
	for k, v := range c.CStore {
		cookies = append(cookies, &http.Cookie{Name: k, Value: v})
	}
	return cookies
}

func (c *session) login() error {
	cli := http.Client{Transport: c.RTer, Jar: c}
	res, err := cli.Get("https://www.hackerrank.com/login")
	if err != nil {
		return err
	}

	csrf, err := getBodyCSRF(res.Body)
	if err != nil {
		return err
	}
	c.RTer.(*csrfTransport).CSRF = csrf

	v := url.Values{
		"login":       []string{os.Getenv("HR_USER")},
		"password":    []string{os.Getenv("HR_PASS")},
		"fallback":    []string{"true"},
		"remember_me": []string{"false"},
	}
	if v.Get("login") == "" {
		v.Set("login", prompt.String("Enter your HackerRank username"))
	}

	if v.Get("password") == "" {
		v.Set("password", prompt.Password("Please enter your HackerRank password"))
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

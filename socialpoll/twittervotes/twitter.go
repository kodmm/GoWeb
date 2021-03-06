package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/garyburd/go-oauth/oauth"
	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

var conn net.Conn

func dial(netw, addr string) (net.Conn, error) {
	if conn != nil {
		conn.Close()
		conn = nil
	}

	netc, err := net.DialTimeout(netw, addr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	conn = netc
	return netc, nil
}

var reader io.ReadCloser

func closeConn() {
	if conn != nil {
		conn.Close()
	}
	if reader != nil {
		reader.Close()
	}
}

var (
	authClient *oauth.Client
	creds      *oauth.Credentials
)

func setupTwitterAuth() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	var ts struct {
		ConsumerKey    string `env:"SP_TWITTER_KEY,required"`
		ConsumerSecret string `env:"SP_TWITTER_SECRET, required"`
		AccessToken    string `env:"SP_TWITTER_ACCESSTOKEN, required"`
		AccessSecret   string `env:"SP_TWITTER_ACCESSSECRET, required"`
	}
	if err := envdecode.Decode(&ts); err != nil {
		log.Fatalln(err)
	}
	creds = &oauth.Credentials{
		Token:  ts.AccessToken,
		Secret: ts.AccessSecret,
	}
	authClient = &oauth.Client{
		Credentials: oauth.Credentials{
			Token:  ts.ConsumerKey,
			Secret: ts.ConsumerSecret,
		},
	}
}

var (
	authSetupOnce sync.Once
	httpClient    *http.Client
)

func makeRequest(query url.Values) (*http.Request, error) {
	authSetupOnce.Do(
		func() {
			setupTwitterAuth()
		})
	const endpoint = "https://stream.twitter.com/1.1/statuses/filter.json"
	req, err := http.NewRequest("POST", endpoint,
		strings.NewReader(query.Encode()))
	if err != nil {
		return nil, err
	}

	formEnc := query.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(formEnc)))
	ah := authClient.AuthorizationHeader(creds, "POST", req.URL, query)
	req.Header.Set("Authorization", ah)

	return req, nil
}

type tweet struct {
	Text string
}

func readFromTwitter(ctx context.Context, votes chan<- string) {
	options, err := loadOptions()
	if err != nil {
		log.Println("?????????????????????????????????????????????:", err)
		return
	}
	query := make(url.Values)
	query.Set("track", strings.Join(options, ","))
	req, err := makeRequest(query)
	if err != nil {
		log.Println("?????????????????????????????????????????????????????????:", err)
		return
	}
	const timeout = 1 * time.Minute
	client := &http.Client{}
	if deadline, ok := ctx.Deadline(); ok {
		client.Timeout = deadline.Sub(time.Now())
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("?????????????????????????????????????????????:", err)
		return
	}
	done := make(chan struct{})
	// ?????????goroutine?????????????????????readFromTwitter?????????????????????
	defer func() { <-done }()

	// ?????????????????????resp.Body??????????????????
	defer resp.Body.Close()

	go func() {
		defer close(done)
		log.Println("resp:", resp.StatusCode)
		if resp.StatusCode != 200 {
			var buf bytes.Buffer
			io.Copy(&buf, resp.Body)
			log.Println("resp body: %s", buf.String())
			return
		}

		decoder := json.NewDecoder(resp.Body)
		for {
			var tweet tweet
			if err := decoder.Decode(&tweet); err != nil {
				break
			}
			for _, option := range options {
				if strings.Contains(strings.ToLower(tweet.Text), strings.ToLower(option)) {
					log.Println("??????:", option)
					votes <- option
				}
			}
		}
	}()

	select {
	case <-time.After(timeout):
	case <-ctx.Done():
	case <-done:
	}

}

func readFromTwitterWithTimeout(ctx context.Context,
	timeout time.Duration, votes chan<- string) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	readFromTwitter(ctx, votes)
}

func twitterStream(ctx context.Context, votes chan<- string) {
	defer close(votes)
	for {
		log.Println("Twitter?????????????????????????????????")
		readFromTwitterWithTimeout(ctx, 1*time.Minute, votes)
		log.Println("?????????")
		select {
		case <-ctx.Done():
			log.Println("Twitter??????????????????????????????????????????")
			return
		case <-time.After(10 * time.Second):
		}
	}
}

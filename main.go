package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ping struct {
	Ping string `json:"message"`
}

type gip struct {
	Gip string `json:"gip"`
}

type Client struct {
	ping
	BaseURL   *url.URL
	UserAgent string

	httpClient *http.Client
}

func main() {
	httpClient := &http.Client{}
	c := Client{httpClient: httpClient}
	u, _ := url.Parse("http://localhost:8080")
	c.BaseURL = u
	c.UserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.121 Safari/537.36"
	c.getPing()
	c.getGip()
}

func (c *Client) getGip() {
	req, err := c.newRequest("GET", "/test", nil)
	if err != nil {
		return
	}

	record := gip{}

	_, err = c.do(req, &record)
	fmt.Println(record.Gip)
}

func (c *Client) getPing() {
	req, err := c.newRequest("GET", "/ping", nil)
	if err != nil {
		return
	}

	record := ping{}

	_, err = c.do(req, &record)
	fmt.Println(record.Ping)

}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	rel := &url.URL{Path: path}
	u := c.BaseURL.ResolveReference(rel)
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(v)
	return resp, err
}

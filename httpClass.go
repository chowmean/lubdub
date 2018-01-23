package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
)

const baseURL string = "http://127.0.0.1:5005/"

type Client struct {
	Token string
}

func BasicAuthClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

type Status struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func (s *Client) PostStatus(status *Status) error {
	url := fmt.Sprintf(baseURL)
	fmt.Println(url)
	j, err := json.Marshal(status)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	_, err = s.doRequest(req)
	return err
}

func (s *Client) doRequest(req *http.Request) ([]byte, error) {
	//req.SetBasicAuth(s.Token)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

func main() {
	client := BasicAuthClient("Token")
	status := Status{
		Content: "New Todo",
		ID:    12,
	}
	// Add a todo
	client.PostStatus(&status)
}

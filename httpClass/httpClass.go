package httpClass

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"bytes"
)

type Client struct {
	Token string
}

func BasicAuthClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

type Content struct {
	Hostname string `json:"hostname"`
	ID      int    `json:"id"`
	Content string `json:"content"`
}

func (s *Client) PostStatus(content *Content, baseurl string) error {
	url := fmt.Sprintf(baseurl)
	j, err := json.Marshal(content)
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


package httpclient

import (
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

var client = resty.New().SetTimeout(30 * time.Second)

type Header map[string]string

func Get(url string) ([]byte, error) {
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func PostJSON(url string, body interface{}) ([]byte, error) {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func PostForm(url string, data map[string]string) ([]byte, error) {
	resp, err := client.R().
		SetFormData(data).
		Post(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

// GetJSON is a convenience wrapper: GET + JSON unmarshal into result.
func GetJSON(url string, result interface{}) error {
	body, err := Get(url)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}

// PostJSONWith is a convenience wrapper: POST JSON body + JSON unmarshal response.
func PostJSONWith(url string, body interface{}, result interface{}) error {
	resp, err := PostJSON(url, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp, result)
}

// New creates a reusable client with custom config.
func New() *Client {
	return &Client{
		client: resty.New().SetTimeout(30 * time.Second),
	}
}

type Client struct {
	client *resty.Client
}

func (c *Client) SetTimeout(d time.Duration) *Client {
	c.client.SetTimeout(d)
	return c
}

func (c *Client) SetBaseURL(url string) *Client {
	c.client.SetBaseURL(url)
	return c
}

func (c *Client) SetAuthToken(token string) *Client {
	c.client.SetAuthToken(token)
	return c
}

func (c *Client) SetHeader(key, value string) *Client {
	c.client.SetHeader(key, value)
	return c
}

func (c *Client) Get(url string) ([]byte, error) {
	resp, err := c.client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func (c *Client) GetJSON(url string, result interface{}) error {
	body, err := c.Get(url)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}

func (c *Client) PostJSON(url string, body interface{}) ([]byte, error) {
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func (c *Client) PostJSONWith(url string, body interface{}, result interface{}) error {
	resp, err := c.PostJSON(url, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(resp, result)
}

func (c *Client) PostForm(url string, data map[string]string) ([]byte, error) {
	resp, err := c.client.R().
		SetFormData(data).
		Post(url)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

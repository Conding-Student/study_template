package http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// Client ...
var Client C

// C ...
type C struct {
	Client *http.Client
}

// New ...
func (c *C) New(client *http.Client) {
	c.Client = client
}

// ReadBodyRequest ...
func ReadBodyRequest(method, url string, input interface{}, headers http.Header) ([]byte, error) {
	message, err := json.Marshal(input)

	if err == nil {
		headers.Set("Content-Type", "application/json")
		resp, err := Request(method, url, bytes.NewBuffer(message), headers)

		if err == nil {
			defer resp.Body.Close()
			return ioutil.ReadAll(resp.Body)
		}
	}

	return nil, err
}

// RequestJSON ...
func RequestJSON(method, url string, input, output, outputErr interface{}, headers http.Header) (*http.Response, error) {
	message, err := json.Marshal(input)

	if err == nil {
		headers.Set("Content-Type", "application/json")
		resp, err := Request(method, url, bytes.NewBuffer(message), headers)

		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err == nil {

				if resp.StatusCode == 400 {

					return resp, json.Unmarshal(body, &outputErr)
				} else {
					return resp, json.Unmarshal(body, &output)
				}

			}
		}
	}

	return nil, err
}

// Request ...
func Request(method, url string, body io.Reader, headers http.Header) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header = headers

	if err == nil {
		return Client.Client.Do(req)
	}

	return nil, err
}

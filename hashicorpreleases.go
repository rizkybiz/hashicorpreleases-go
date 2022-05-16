package hashicorpreleases

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Client represents an HTTP client for interfacing with the
// HashiCorp Releases API
type Client struct {
	URL        string
	HTTPClient *http.Client
}

type errorResponse struct {
	Status  int    `json:"code"`
	Message string `json:"message"`
}

// NewClient returns a new hashicorpreleases client. Provide a
// custom releases endpoint by setting RELEASES_URL in the
// environment
func NewClient() *Client {

	// Check if a URL is provided via ENV VARS
	url := os.Getenv("RELEASES_URL")
	if url == "" {
		url = "https://api.releases.hashicorp.com/v1"
	}

	// Setup the client and return
	return &Client{
		URL: url,
		HTTPClient: &http.Client{
			Timeout: time.Minute * 1,
		},
	}
}

// sendRequest assumes proper "content-type" header is set
// and that a body is attached if necessary to the http request
func (c *Client) sendRequest(req *http.Request, v interface{}) error {

	// Set the appropriate headers
	req.Header.Set("Accept", "application/json; charset=utf-8")

	// execute the http request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check for non OK status code and attempt to decode into errorResponse
	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err = json.NewDecoder(res.Body).Decode(&errRes); err == nil {
			return fmt.Errorf("error: %s; status code: %d", errRes.Message, res.StatusCode)
		}
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	// Attempt to decode response into whichever interface was provided
	err = json.NewDecoder(res.Body).Decode(&v)
	if err != nil {
		return fmt.Errorf("error decoding response body: %s", err)
	}
	return nil
}

func setJSONHeader(r *http.Request) {
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
}

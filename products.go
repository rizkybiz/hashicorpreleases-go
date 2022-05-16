package hashicorpreleases

import (
	"fmt"
	"net/http"
)

// ProductResponse is a list of all HashiCorp products
type ProductResponse []string

// GetProducts retrieves a list of all of the HashiCorp products
func (c *Client) GetProducts() (ProductResponse, error) {

	// Start by creating request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/products", c.URL), nil)
	if err != nil {
		return nil, err
	}
	setJSONHeader(req)

	// Issue the request against the API
	res := ProductResponse{}
	if err = c.sendRequest(req, &res); err != nil {
		return nil, err
	}
	return res, nil
}

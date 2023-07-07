package outline

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

// Client is per server top level client which acts as entry point and stores common configuration (like base url) for
// resource level clients. It is preferred to reuse same client while communicating to the same server as this makes
// better utilization of resources.
type Client struct {
	base *sling.Sling
}

// New creates and returns a new per server client.
func New(baseURL string, hc *http.Client, apiKey string) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	sl := sling.New().Client(hc).Base(baseURL)
	sl.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	sl.Set("Content-Type", "application/json")
	sl.Set("Accept", "application/json")

	return &Client{base: sl}
}

// Documents creates client for operating on documents.
func (cl *Client) Documents() *DocumentsClient { return nil }

// Collections creates client for operating on collections.
func (cl *Client) Collections() *CollectionsClient {
	return &CollectionsClient{sl: cl.base.New()}
}

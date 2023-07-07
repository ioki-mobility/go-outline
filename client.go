package outline

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
	"github.com/rsjethani/secret/v2"
)

// Client is per server top level client which acts as entry point and stores common configuration (like base url) for
// resource level clients. It is preferred to reuse same client while communicating to the same server as this makes
// better utilization of resources.
type Client struct {
	base *sling.Sling
}

// New creates and returns a new per server client.
func New(hc *http.Client, baseURL string, apiKey secret.Text) *Client {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL = baseURL + "/"
	}

	cl := &Client{
		base: sling.New().Client(hc).Set("Authorization", fmt.Sprintf("Bearer %s", apiKey)).Base(baseURL),
	}

	return cl
}

// Documents creates client for operating on documents.
func (cl *Client) Documents() *DocumentsClient { return nil }

// Collections creates client for operating on collections.
func (cl *Client) Collections() *CollectionsClient { return nil }

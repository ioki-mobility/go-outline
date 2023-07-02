package outline

import "github.com/rsjethani/secret/v2"

// Client is per server top level client which acts as entry point and stores common configuration (like base url) for
// resource level clients. It is preferred to reuse same client while communicating to the same server as this makes
// better utilization of resources.
type Client struct{}

// New creates and returns a new per server client.
func New(baseURL string, apiKey secret.Text) *Client { return nil }

// Documents creates client for operating on documents.
func (cl *Client) Documents() *DocumentsClient { return nil }

// Documents creates client for operating on collections.
func (cl *Client) Collections() *CollectionsClient { return nil }

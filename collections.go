package outline

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
)

// CollectionsClient exposes CRUD operations around the collections resource.
type CollectionsClient struct {
	sl *sling.Sling
}

func (cl *CollectionsClient) Get(id CollectionID) *CollectionsGetClient {
	data := struct {
		ID CollectionID `json:"id"`
	}{ID: id}

	sl := cl.sl.New().Post("collections.info").BodyJSON(&data)
	return &CollectionsGetClient{sl: sl}
}

type CollectionsGetClient struct {
	sl *sling.Sling
}

func (cl *CollectionsGetClient) Do(ctx context.Context) (*Collection, error) {
	success := struct {
		Data Collection `json:"data"`
	}{}

	br, err := request(ctx, cl.sl, &success)
	if err != nil {
		return nil, fmt.Errorf("failed making HTTP request: %w", err)
	}
	if br != nil {
		return nil, fmt.Errorf("bad response: %w", &apiError{br: *br})
	}

	return &success.Data, nil
}

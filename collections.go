package outline

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/ioki-mobility/go-outline/internal/common"
)

// CollectionsClient exposes CRUD operations around the collections resource.
type CollectionsClient struct {
	sl *sling.Sling
}

func (cl *CollectionsClient) Get(id CollectionID) *CollectionsGetClient {
	data := struct {
		ID CollectionID `json:"id"`
	}{ID: id}

	copy := cl.sl.New()
	copy.BodyJSON(&data).Post(common.CollectionsGetEndpoint())

	return &CollectionsGetClient{sl: copy}
}

func (cl *CollectionsClient) List() *CollectionsListClient {
	copy := cl.sl.New()
	copy.Post(common.CollectionsListEndpoint())

	return &CollectionsListClient{sl: copy}
}

type CollectionsGetClient struct {
	sl *sling.Sling
}

func (cl *CollectionsGetClient) Do(ctx context.Context) (*Collection, error) {
	success := &struct {
		Data *Collection `json:"data"`
	}{}

	br, err := request(ctx, cl.sl, success)
	if err != nil {
		return nil, fmt.Errorf("failed making HTTP request: %w", err)
	}
	if br != nil {
		return nil, fmt.Errorf("bad response: %w", &apiError{br: *br})
	}

	return success.Data, nil
}

type CollectionsListClient struct {
	sl *sling.Sling
}

// CollectionsListFn is the type of function called by [CollectionsListClient.Do] for every new collecion it finds.
type CollectionsListFn func(*Collection, error) (bool, error)

// Do makes the actual request for listing all collections. If the request is successful then fn is called sequentially
// with every collection received. But if there is some error/bad response then fn is called with the error. If fn
// returns false then the whole process is aborted otherwise the request is retried. NOTE: Policies if any returned are
// ignored as of now. Later if we find them important then we can include them too.
func (cl *CollectionsListClient) Do(ctx context.Context, fn CollectionsListFn) error {
	success := &struct {
		Data       []*Collection `json:"data"`
		Pagination pagination    `json:"pagination"`
	}{}

	params := &paginationQueryParams{}
	for {
		copy := cl.sl.New().QueryStruct(params)

		// Make the request and see if there is an error/bad response. If there is one then give fn the error ask for
		// its intention. If fn still wants to continue the we abort further processing in current iteration and
		// basically retry the same request again.
		br, err := request(ctx, copy, success)
		if err != nil {
			err = fmt.Errorf("failed making HTTP request: %w", err)
		}
		if br != nil {
			err = fmt.Errorf("bad response: %w", &apiError{br: *br})
		}
		if err != nil {
			if ok, e := fn(nil, err); !ok {
				return e
			}
			continue
		}

		// If we are here then it means there was no error/bad response while fetching current page
		// so lets iterate over page items.
		for _, col := range success.Data {
			if ok, e := fn(col, nil); !ok {
				return e
			}
		}

		// If there are more than one items in current list then there could be more items remaining to be fetched. In
		// that case we adjust offset for next request. If there are no items or just a single item in the list that
		// means there are no more items to be fetched and we are done.
		if len(success.Data) <= 1 {
			return nil
		}
		params.Offset += len(success.Data)
	}
}

// pagination represents pagination logic related metadata usually part of responses containing list of items.
type pagination struct {
	Limit    int    `json:"limit"`
	Offset   int    `json:"offset"`
	NextPath string `json:"nextPath"`
}

// paginationQueryParams contains valid query paramters for pagination logic.
// Reference: https://www.getoutline.com/developers#section/Pagination
type paginationQueryParams struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

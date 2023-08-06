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

func newCollectionsClient(sl *sling.Sling) *CollectionsClient {
	return &CollectionsClient{sl: sl}
}

// Structure gives access to id's document structure.
func (cl *CollectionsClient) Structure(id CollectionID) *CollectionsStructureClient {
	return newCollectionsStructureClient(cl.sl, id)
}

func (cl *CollectionsClient) Get(id CollectionID) *CollectionsGetClient {
	return newCollectionsGetClient(cl.sl, id)
}

func (cl *CollectionsClient) List() *CollectionsListClient {
	return newCollectionListClient(cl.sl)
}

// Create returns a client for creating a collection.
// API reference: https://www.getoutline.com/developers#tag/Collections/paths/~1collections.create/post
func (cl *CollectionsClient) Create(name string) *CollectionsCreateClient {
	return newCollectionsCreateClient(cl.sl, name)
}

type CollectionsStructureClient struct {
	sl *sling.Sling
}

func newCollectionsStructureClient(sl *sling.Sling, id CollectionID) *CollectionsStructureClient {
	data := struct {
		ID CollectionID `json:"id"`
	}{ID: id}

	copy := sl.New()
	copy.Post(common.CollectionsStructureEndpoint()).BodyJSON(&data)

	return &CollectionsStructureClient{sl: copy}
}

// Do makes the actual request for fetching the collection's document structure. The structure is essentially a summary
// of all associated documents.
func (cl *CollectionsStructureClient) Do(ctx context.Context) ([]DocumentSummary, error) {
	success := &struct {
		Data []DocumentSummary `json:"data"`
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

type CollectionsGetClient struct {
	sl *sling.Sling
}

func newCollectionsGetClient(sl *sling.Sling, id CollectionID) *CollectionsGetClient {
	data := struct {
		ID CollectionID `json:"id"`
	}{ID: id}

	copy := sl.New()
	copy.Post(common.CollectionsGetEndpoint()).BodyJSON(&data)

	return &CollectionsGetClient{sl: copy}
}

// Do makes the actual request for fetching said collection info.
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

func newCollectionListClient(sl *sling.Sling) *CollectionsListClient {
	copy := sl.New()
	copy.Post(common.CollectionsListEndpoint())

	return &CollectionsListClient{sl: copy}
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
		// Create a fresh copy of original request for every page then set query parameters accordingly.
		copy := cl.sl.New().QueryStruct(params)

		// Make the request and see if there is an error/bad response. If there is one then give fn the error ask for
		// its intention. If fn still wants to continue then we abort further processing in current iteration and
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
		// so let's iterate over page items.
		for _, col := range success.Data {
			if ok, e := fn(col, nil); !ok {
				return e
			}
		}

		// If there is more than one item in current list then there could be more items remaining to be fetched. In
		// that case we adjust offset for next request. If there are no items or just a single item in the list that
		// means there are no more items to be fetched, and we are done.
		if len(success.Data) <= 1 {
			return nil
		}
		params.Offset += len(success.Data)
	}
}

// collectionsCreateParams represents the Outline Collections.create parameters
type collectionsCreateParams struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Permission  string `json:"permission,omitempty"`
	Color       string `json:"color,omitempty"`
	Private     bool   `json:"private,omitempty"`
}

type CollectionsCreateClient struct {
	sl     *sling.Sling
	params collectionsCreateParams
}

func newCollectionsCreateClient(sl *sling.Sling, name string) *CollectionsCreateClient {
	copy := sl.New()
	params := collectionsCreateParams{Name: name}
	return &CollectionsCreateClient{sl: copy, params: params}
}

func (cl *CollectionsCreateClient) Description(desc string) *CollectionsCreateClient {
	cl.params.Description = desc
	return cl
}

func (cl *CollectionsCreateClient) PermissionRead() *CollectionsCreateClient {
	cl.params.Permission = "read"
	return cl
}

func (cl *CollectionsCreateClient) PermissionReadWrite() *CollectionsCreateClient {
	cl.params.Permission = "read_write"
	return cl
}

func (cl *CollectionsCreateClient) Color(color string) *CollectionsCreateClient {
	cl.params.Color = color
	return cl
}

func (cl *CollectionsCreateClient) Private(private bool) *CollectionsCreateClient {
	cl.params.Private = private
	return cl
}

// Do make the actual request to create a collection.
func (cl *CollectionsCreateClient) Do(ctx context.Context) (*Collection, error) {
	cl.sl.Post(common.CollectionsCreateEndpoint()).BodyJSON(&cl.params)

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

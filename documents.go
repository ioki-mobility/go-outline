package outline

import (
	"context"
	"fmt"

	"github.com/dghubble/sling"
	"github.com/ioki-mobility/go-outline/internal/common"
)

// DocumentsClient exposes CRUD operations around the documents resource.
type DocumentsClient struct {
	sl *sling.Sling
}

// newDocumentsClient creates a new instance of DocumentsClient.
func newDocumentsClient(sl *sling.Sling) *DocumentsClient {
	return &DocumentsClient{sl: sl}
}

// Get returns a client for retriving a single document.
func (cl *DocumentsClient) Get() *DocumentsClientGet {
	return nil
}

// GetAll returns a client for retrieving multiple documents at once.
func (cl *DocumentsClient) GetAll() *DocumentsClientGetAll {
	return nil
}

// Create returns a client for creating a single document in the specified collection.
// API reference: https://www.getoutline.com/developers#tag/Documents/paths/~1documents.create/post
func (cl *DocumentsClient) Create(title string, id CollectionID) *DocumentsCreateClient {
	return newDocumentsCreateClient(cl.sl, title, id)
}

// Update returns a client for updating a single document in the specified collection.
// API reference: https://www.getoutline.com/developers#tag/Documents/paths/~1documents.update/post
func (cl *DocumentsClient) Update(id DocumentID) *DocumentsUpdateClient {
	return newDocumentsUpdateClient(cl.sl, id)
}

// DocumentsClientGet gets a single document.
type DocumentsClientGet struct{}

// ByID configures that document be retrieved by its id.
func (cl *DocumentsClientGet) ByID(id DocumentID) *DocumentsClientGet {
	return nil
}

// ByShareID configures that document be retrieved by its share id.
func (cl *DocumentsClientGet) ByShareID(id DocumentShareID) *DocumentsClientGet {
	return nil
}

// Do makes the actual request and returns the document.
func (cl *DocumentsClientGet) Do(ctx context.Context) (*Document, error) { return nil, nil }

// DocumentsClientGetAll can be used to retrieve more than one document. Use available configuration options to select
// the documents you want to retrieve then finally call [DocumentsClientGetAll.Do].
type DocumentsClientGetAll struct{}

// Collection selects documents belonging to the collection identified by id.
func (cl *DocumentsClientGetAll) Collection(id CollectionID) *DocumentsClientGetAll { return nil }

// Parent selects documents having the parent document identified by id.
func (cl *DocumentsClientGetAll) Parent(id DocumentID) *DocumentsClientGetAll { return nil }

// Do makes the actual request and retrieves selected documents. The user provided callback fn is called for every such
// document. If there is any error during the process then fn is given the error so that it can decide whether to
// continue or not. The callback can return false in case it wants to abort getting documents.
func (cl *DocumentsClientGetAll) Do(ctx context.Context, fn func(*Document, error) bool) error {
	return nil
}

// documentsCreateParams represents the Outline Documents.create parameters
type documentsCreateParams struct {
	CollectionID     CollectionID `json:"collectionId"`
	ParentDocumentId DocumentID   `json:"parentDocumentId,omitempty"`
	Publish          bool         `json:"publish,omitempty"`
	Template         bool         `json:"template,omitempty"`
	TemplateID       TemplateID   `json:"templateId,omitempty"`
	Text             string       `json:"text,omitempty"`
	Title            string       `json:"title"`
}

// DocumentsCreateClient is a client for creating a single document.
type DocumentsCreateClient struct {
	sl     *sling.Sling
	params documentsCreateParams
}

func newDocumentsCreateClient(sl *sling.Sling, title string, collectionId CollectionID) *DocumentsCreateClient {
	copy := sl.New()
	params := documentsCreateParams{Title: title, CollectionID: collectionId}
	return &DocumentsCreateClient{sl: copy, params: params}
}

func (cl *DocumentsCreateClient) Publish(publish bool) *DocumentsCreateClient {
	cl.params.Publish = publish
	return cl
}

func (cl *DocumentsCreateClient) Text(text string) *DocumentsCreateClient {
	cl.params.Text = text
	return cl
}

func (cl *DocumentsCreateClient) ParentDocumentID(id DocumentID) *DocumentsCreateClient {
	cl.params.ParentDocumentId = id
	return cl
}

func (cl *DocumentsCreateClient) TemplateID(id TemplateID) *DocumentsCreateClient {
	cl.params.TemplateID = id
	return cl
}

func (cl *DocumentsCreateClient) Template(template bool) *DocumentsCreateClient {
	cl.params.Template = template
	return cl
}

// Do makes the actual request to create a document.
func (cl *DocumentsCreateClient) Do(ctx context.Context) (*Document, error) {
	cl.sl.Post(common.DocumentsCreateEndpoint()).BodyJSON(&cl.params)

	success := &struct {
		Data *Document `json:"data"`
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

// documentsUpdateParams represents the Outline Documents.update parameters
type documentsUpdateParams struct {
	Id      DocumentID `json:"id"`
	Title   string     `json:"title,omitempty"`
	Text    string     `json:"text,omitempty"`
	Append  bool       `json:"append,omitempty"`
	Publish bool       `json:"publish,omitempty"`
	Done    bool       `json:"done,omitempty"`
}

// DocumentsUpdateClient is a client for updating a single document.
type DocumentsUpdateClient struct {
	sl     *sling.Sling
	params documentsUpdateParams
}

func newDocumentsUpdateClient(sl *sling.Sling, id DocumentID) *DocumentsUpdateClient {
	copy := sl.New()
	params := documentsUpdateParams{Id: id}
	return &DocumentsUpdateClient{sl: copy, params: params}
}

func (cl *DocumentsUpdateClient) Title(title string) *DocumentsUpdateClient {
	cl.params.Title = title
	return cl
}

func (cl *DocumentsUpdateClient) Text(text string) *DocumentsUpdateClient {
	cl.params.Text = text
	return cl
}

func (cl *DocumentsUpdateClient) Publish(publish bool) *DocumentsUpdateClient {
	cl.params.Publish = publish
	return cl
}

func (cl *DocumentsUpdateClient) Append(append bool) *DocumentsUpdateClient {
	cl.params.Append = append
	return cl
}

func (cl *DocumentsUpdateClient) Done(done bool) *DocumentsUpdateClient {
	cl.params.Done = done
	return cl
}

// Do makes the actual request to update a document.
func (cl *DocumentsUpdateClient) Do(ctx context.Context) (*Document, error) {
	cl.sl.Post(common.DocumentsUpdateEndpoint()).BodyJSON(&cl.params)

	success := &struct {
		Data *Document `json:"data"`
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

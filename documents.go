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

// GetAll returns a client for retriving multiple documents at once.
func (cl *DocumentsClient) GetAll() *DocumentsClientGetAll {
	return nil
}

// DocumentsClientGet gets a single document.
type DocumentsClientGet struct{}

// ByID configures that document be retrieved by its id.
func (cl *DocumentsClientGet) ByID(id DocumentID) *DocumentsClientGet {
	return nil
}

// GetByID configures that document be retrieved by its share id.
func (cl *DocumentsClientGet) ByShareID(id DocumentShareID) *DocumentsClientGet {
	return nil
}

// Do makes the actual request and returns the document.
func (cl *DocumentsClientGet) Do(ctx context.Context) (*Document, error) { return nil, nil }

// DocumentsClientGetAll can be used to retrieve more than one document. Use available configuration options to select
// the documents you want to retrive then finall call [DocumentsClientGetAll.Do].
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

// Create returns a client for creating a single document in the specified collection.
//
// API reference: https://www.getoutline.com/developers#tag/Documents/paths/~1documents.create/post
func (cl *DocumentsClient) Create(collectionID CollectionID, title string, text string, publish bool, parentDocumentId ParentDocumentID, templateId TemplateID, template bool) *DocumentsClientCreate {
	return newDocumentsClientCreate(cl.sl, collectionID, title, text, publish, parentDocumentId, templateId, template)
}

// DocumentsClientCreate is a client for creating a single document.
type DocumentsClientCreate struct {
	sl *sling.Sling
}

// newDocumentsClientCreate creates a new DocumentsClientCreate instance to create a single document in the specified collection.
func newDocumentsClientCreate(sl *sling.Sling, collectionID CollectionID, title string, text string, publish bool, parentDocumentId ParentDocumentID, templateId TemplateID, template bool) *DocumentsClientCreate {
	data := struct {
		CollectionID     CollectionID     `json:"collectionId"`
		Title            string           `json:"title"`
		Publish          bool             `json:"publish,omitempty"`
		Text             string           `json:"text,omitempty"`
		ParentDocumentId ParentDocumentID `json:"parentDocumentId,omitempty"`
		TemplateID       TemplateID       `json:"templateId,omitempty"`
		Template         bool             `json:"template,omitempty"`
	}{
		CollectionID:     collectionID,
		Title:            title,
		Publish:          publish,
		Text:             text,
		ParentDocumentId: parentDocumentId,
		TemplateID:       templateId,
		Template:         template,
	}

	copy := sl.New()
	copy.Post(common.DocumentsCreateEndpoint()).BodyJSON(&data)
	return &DocumentsClientCreate{sl: copy}
}

// Do makes the actual request to create a document.
func (cl *DocumentsClientCreate) Do(ctx context.Context) (*Document, error) {
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

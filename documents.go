package outline

import "context"

// DocumentsClient exposes CRUD operations around the documents resource.
type DocumentsClient struct{}

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

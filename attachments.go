package outline

import (
	"context"
	"fmt"

	"github.com/ioki-mobility/go-outline/internal/common"
	"github.com/rsjethani/rsling"
)

// AttachmentsClient exposes CRUD operations around the attachments resource.
type AttachmentsClient struct {
	sl *rsling.Sling
}

// newAttachmentsClient creates a new instance of AttachmentsClient.
func newAttachmentsClient(sl *rsling.Sling) *AttachmentsClient {
	return &AttachmentsClient{sl: sl}
}

// Create returns a client for creating a single attachment in the specified collection.
// API reference: https://www.getoutline.com/developers#tag/Attachments/paths/~1attachments.create/post
func (cl *AttachmentsClient) Create(name string, contentType string, size int) *AttachmentCreateClient {
	return newAttachmentCreateClient(cl.sl, name, contentType, size)
}

// attachmentsCreateParams represents the Outline Attachment.create parameters
type attachmentsCreateParams struct {
	Name        string     `json:"name"`
	ContentType string     `json:"contentType"`
	Size        int        `json:"size"`
	DocumentID  DocumentID `json:"documentId,omitempty"`
}

// AttachmentCreateClient is a client for creating a single attachment.
type AttachmentCreateClient struct {
	sl     *rsling.Sling
	params attachmentsCreateParams
}

func newAttachmentCreateClient(sl *rsling.Sling, name string, contentType string, size int) *AttachmentCreateClient {
	copy := sl.New()
	params := attachmentsCreateParams{Name: name, ContentType: contentType, Size: size}

	return &AttachmentCreateClient{sl: copy, params: params}
}

func (cl *AttachmentCreateClient) DocumentID(id DocumentID) *AttachmentCreateClient {
	cl.params.DocumentID = id
	return cl
}

// Do makes the actual request to create a attachment.
func (cl *AttachmentCreateClient) Do(ctx context.Context) (*Attachment, error) {
	cl.sl.Post(common.AttachmentsCreateEndpoint()).BodyJSON(&cl.params)

	success := &struct {
		Data *Attachment `json:"data"`
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

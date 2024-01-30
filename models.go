package outline

import "time"

type (
	DocumentID      string
	DocumentShareID string
	DocumentUrlID   string
	CollectionID    string
	TemplateID      string
)

// DocumentSummary represents summary of a document (and its children) that is part of a collection.
type DocumentSummary struct {
	ID       DocumentID        `json:"id"`
	Title    string            `json:"title"`
	URL      string            `json:"url"`
	Children []DocumentSummary `json:"children"`
}

// Document represents an outline document.
type Document struct {
	ID               DocumentID   `json:"id"`
	CollectionID     CollectionID `json:"collectionId"`
	ParentDocumentID DocumentID   `json:"parentDocumentId"`
	Title            string       `json:"title"`
	FullWidth        bool         `json:"fullWidth"`
	Emoji            string       `json:"emoji"`
	Text             string       `json:"text"`
	URLID            string       `json:"urlId"`
	Collaborators    []User       `json:"collaborators"`
	Pinned           bool         `json:"pinned"`
	Template         bool         `json:"template"`
	TemplateID       TemplateID   `json:"templateId"`
	Revision         int          `json:"revision"`
	CreatedAt        time.Time    `json:"createdAt"`
	CreatedBy        User         `json:"createdBy"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	UpdatedBy        User         `json:"updatedBy"`
	PublishedAt      time.Time    `json:"publishedAt"`
	ArchivedAt       time.Time    `json:"archivedAt"`
	DeletedAt        time.Time    `json:"deletedAt"`
}

// User represents an outline user.
type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	AvatarURL    string    `json:"avatarUrl"`
	Email        string    `json:"email"`
	IsAdmin      bool      `json:"isAdmin"`
	IsSuspended  bool      `json:"isSuspended"`
	LastActiveAt time.Time `json:"lastActiveAt"`
	CreatedAt    time.Time `json:"createdAt"`
}

// Collection represents an outline collection.
type Collection struct {
	ID          CollectionID   `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Sort        map[string]any `json:"sort"`
	Index       string         `json:"index"`
	Color       string         `json:"color"`
	Icon        string         `json:"icon"`
	Permission  string         `json:"permission"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   time.Time      `json:"deletedAt"`
}

type Attachment struct {
	MaxUploadSize  int                    `json:"maxUploadSize"`
	UploadURL      string                 `json:"uploadUrl"`
	Form           map[string]interface{} `json:"form"`
	AttachmentData AttachmentData         `json:"attachment"`
}

type AttachmentData struct {
	ContentType string `json:"contentType"`
	Size        int    `json:"size"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	DocumentID  string `json:"documentId"`
}

package outline

import "time"

type (
	DocumentID      string
	DocumentShareID string
	DocumentUrlID   string
	CollectionID    string
	TemplateID      string
)

// CollectionDocument represents summary of a document (and its children) that is part of a collection.
type CollectionDocument struct {
	ID       DocumentID
	Title    string
	URL      string
	Children []CollectionDocument
}

// Document represents an outline document.
type Document struct {
	ID               DocumentID    `json:"id"`
	CollectionID     CollectionID  `json:"collectionId"`
	ParentDocumentID DocumentID    `json:"parentDocumentId"`
	Title            string        `json:"title"`
	FullWidth        bool          `json:"fullWidth"`
	Emoji            string        `json:"emoji"`
	Text             string        `json:"text"`
	URLID            string        `json:"urlId"`
	Collaborators    []interface{} `json:"collaborators"`
	Pinned           bool          `json:"pinned"`
	Template         bool          `json:"template"`
	TemplateID       TemplateID    `json:"templateId"`
	Revision         int           `json:"revision"`
	CreatedAt        time.Time     `json:"createdAt"`
	CreatedBy        interface{}   `json:"createdBy"`
	UpdatedAt        time.Time     `json:"updatedAt"`
	UpdatedBy        interface{}   `json:"updatedBy"`
	PublishedAt      time.Time     `json:"publishedAt"`
	ArchivedAt       time.Time     `json:"archivedAt"`
	DeletedAt        time.Time     `json:"deletedAt"`
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

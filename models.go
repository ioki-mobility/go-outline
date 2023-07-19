package outline

import "time"

type (
	DocumentID      string
	DocumentShareID string
	DocumentUrlID   string
	CollectionID    string
	CollectionName  string
)

// Document represents an outline document.
type Document struct{}

// Collection represents an outline collection.
type Collection struct {
	ID          DocumentID     `json:"id"`
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

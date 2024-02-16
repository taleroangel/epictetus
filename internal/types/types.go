package types

// Special string type that must be treated as markdown during rendering
type Markdown string

// Special string type that represents a server resource
type Uri string

// Author of a quote
type Author struct {
	Name        string   `json:"name"`
	Description Markdown `json:"description"`
	// Image to show with the author (optional)
	Image Uri `json:"image"`
}

// User stored quote
type Quote struct {
	Id      int32    `json:"id"`
	Content Markdown `json:"content"`
	Author  Author   `json:"author"`
	// Source of the quote (optional)
	Source   Markdown `json:"source,omitempty"`
	Tags     []string `json:"tags"`
	Favorite bool     `json:"favorite"`
}

// User information
type User struct {
	User     string `json:"user"`
	Name     string `json:"name"`
	HashPass string `json:"-"`
	Sudo     bool   `json:"sudo"`
}

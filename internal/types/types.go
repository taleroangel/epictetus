package types

// Special string type that must be treated as markdown during rendering
type Markdown string

// Special string type that represents a server resource
type Uri string

// Author of a quote - 1:1 with database representation
type Author struct {
	Id          int      `json:"id"`
	User        string   `json:"user"`
	Name        string   `json:"name"`
	Description Markdown `json:"description"`
	// Image to show with the author (optional)
	Image Uri `json:"image"`
}

// User stored quote - Not the same stored in the database
type Quote struct {
	Id      int32    `json:"id"`
	User    string   `json:"user"`
	Content Markdown `json:"content"`
	Author  Author   `json:"author"`
	// Source of the quote (optional)
	Source   Markdown `json:"source,omitempty"`
	Favorite bool     `json:"favorite"`
	Tags     []string `json:"tags"`
}

// User information - 1:1 with database representation
type User struct {
	User     string `json:"user"`
	HashPass string `json:"-"`
	Name     string `json:"name"`
	Sudo     bool   `json:"sudo"`
}

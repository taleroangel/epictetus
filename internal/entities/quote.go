package entities

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

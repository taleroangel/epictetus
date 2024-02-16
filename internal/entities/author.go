package entities

// Author of a quote - 1:1 with database representation
type Author struct {
	Id          int      `json:"id"`
	User        string   `json:"user"`
	Name        string   `json:"name"`
	Description Markdown `json:"description"`
	// Image to show with the author (optional)
	Image Uri `json:"image"`
}

package entities

// User information - 1:1 with database representation
type User struct {
	User     string `json:"user"`
	HashPass string `json:"-"`
	Name     string `json:"name"`
	Sudo     bool   `json:"sudo"`
}

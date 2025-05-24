package models

// Client represents a client in the system.
type Client struct {
	ID       int    `json:"id"`
	UUID     string `json:"uuid"`
	UserName string `json:"username"`
	Token    string `json:"token"`
}

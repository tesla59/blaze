package types

// MessageType is a struct that represents the type of message being sent.
// Embedding this struct allows us to easily add more fields in the future if needed.
type MessageType struct {
	Type string `json:"type"`
}

type IdentityMessage struct {
	MessageType
	ClientID string `json:"id"`
}

type Message struct {
	MessageType
	Value string `json:"value"`
}

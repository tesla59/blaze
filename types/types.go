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
	Message string `json:"message"`
}

type DisconnectedMessage struct {
	MessageType
}

type ErrorMessage struct {
	MessageType
	Error string `json:"error"`
}

type MatchedMessage struct {
	MessageType
	ClientID string `json:"client_id"`
}

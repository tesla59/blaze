package types

import "github.com/tesla59/blaze/models"

// MessageType is a struct that represents the type of message being sent.
// Embedding this struct allows us to easily add more fields in the future if needed.
type MessageType struct {
	Type string `json:"type"`
}

type IdentityMessage struct {
	MessageType
	Client *models.Client `json:"client"`
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
	Client models.Client `json:"client"`
}

type State int

const (
	Connected State = iota
	Waiting
	Matched
	Disconnected
)

func (s State) String() string {
	return [...]string{"Connected", "Waiting", "Matched", "Disconnected"}[s]
}

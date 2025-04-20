package matchmaker

import (
	"encoding/json"
	"github.com/tesla59/blaze/types"
)

func DisconnectedMessage() []byte {
	msg := types.DisconnectedMessage{
		MessageType: types.MessageType{
			Type: "disconnected",
		},
	}
	message, _ := json.Marshal(msg)
	return message
}

func ErrorByte(err error) []byte {
	msg := types.ErrorMessage{
		MessageType: types.MessageType{
			Type: "error",
		},
		Error: err.Error(),
	}
	message, _ := json.Marshal(msg)
	return message
}

func MatchedMessage(clientID string) []byte {
	msg := types.MatchedMessage{
		MessageType: types.MessageType{
			Type: "matched",
		},
		ClientID: clientID,
	}
	message, _ := json.Marshal(msg)
	return message
}

package matchmaker

import (
	"encoding/json"
	"github.com/tesla59/blaze/models"
	"github.com/tesla59/blaze/types"
)

func disconnectedMessage() []byte {
	msg := types.DisconnectedMessage{
		MessageType: types.MessageType{
			Type: "disconnected",
		},
	}
	message, _ := json.Marshal(msg)
	return message
}

func errorByte(err error) []byte {
	msg := types.ErrorMessage{
		MessageType: types.MessageType{
			Type: "error",
		},
		Error: err.Error(),
	}
	message, _ := json.Marshal(msg)
	return message
}

func matchedMessage(peer *Client) []byte {
	peerClient := models.Client{
		ID:       peer.ID,
		UUID:     peer.UUID,
		UserName: peer.UserName,
	}
	msg := types.MatchedMessage{
		MessageType: types.MessageType{
			Type: "matched",
		},
		Client: peerClient,
	}
	message, _ := json.Marshal(msg)
	return message
}

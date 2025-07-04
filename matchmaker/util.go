package matchmaker

import (
	"encoding/json"
	"github.com/tesla59/blaze/models"
	"github.com/tesla59/blaze/types"
)

func PeerDisconnectedMessage() []byte {
	msg := types.PeerDisconnectedMessage{
		MessageType: types.MessageType{
			Type: "peer_disconnected",
		},
	}
	message, _ := json.Marshal(msg)
	return message
}

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

func MatchedMessage(peer *Client) []byte {
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

package matchmaker

import (
	"github.com/google/uuid"
	"github.com/tesla59/blaze/config"
	"github.com/tesla59/blaze/log"
	"github.com/tesla59/blaze/models"
	"github.com/tesla59/blaze/types"
	"github.com/tesla59/blaze/utils"
	"strconv"
	"testing"
)

func BenchmarkEnqueue(b *testing.B) {
	config.GetConfig()
	log.Init()
	mm := NewMatchmaker(1000)
	hub := NewHub(mm)
	for i := 0; i < b.N; i++ {
		uuid := uuid.New().String()
		client := &Client{
			Client: &models.Client{
				ID:       i,
				UUID:     uuid,
				Token:    strconv.Itoa(i),
				UserName: utils.GenerateName(uuid),
			},
			State: types.Waiting,
			Hub:   hub,
			Send: make(chan []byte),
		}
		// Simulate a websocket connection
		go func(ch chan []byte) {
			for range ch {
				// Discard received values
			}
		}(client.Send)
		mm.enqueue(client)
	}
}

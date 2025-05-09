package matchmaker

import (
	"encoding/json"
	"github.com/tesla59/blaze/matchmaker"
	"net/http"
)

type QueueStateHandler struct {
	Matchmaker *matchmaker.Matchmaker
}

func NewQueueStateHandler(m *matchmaker.Matchmaker) *QueueStateHandler {
	return &QueueStateHandler{Matchmaker: m}
}

func (h *QueueStateHandler) Handle() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		h.queueHandler(w, r)
	}
}

func (h *QueueStateHandler) queueHandler(w http.ResponseWriter, r *http.Request) {
	queueState := h.Matchmaker.GetQueueState()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(queueState); err != nil {
		http.Error(w, "Failed to encode queue state", http.StatusInternalServerError)
		return
	}
}

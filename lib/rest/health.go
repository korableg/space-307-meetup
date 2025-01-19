package rest

import (
	"net/http"
	"sync/atomic"
)

type health struct {
	shutdown atomic.Bool
}

func (h *health) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if h.shutdown.Load() {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *health) Shutdown() {
	h.shutdown.Store(true)
}

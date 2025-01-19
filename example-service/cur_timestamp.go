package main

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type curTimestamp struct {
}

func (*curTimestamp) Handler() (string, http.Handler) {
	return "GET /time", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(strconv.FormatInt(time.Now().Unix(), 10)))
		if err != nil {
			slog.Warn("err on write response", "err", err)
		}
	})
}

package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/korableg/space-307-meetup/lib/config"
	"github.com/korableg/space-307-meetup/lib/rest"
)

func main() {
	cfg := config.NewConfig()

	srvr := rest.New(cfg.Rest, &curTimestamp{})
	go func() {
		sErr := srvr.Serve()
		if sErr != nil {
			slog.Error(sErr.Error())
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	<-sig

	err := srvr.Shutdown(context.Background())
	if err != nil {
		slog.Error("shutdown error", "err", err)
	}

	slog.Info("bye bye")
}

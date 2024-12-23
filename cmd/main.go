package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"grpc_method_extension/internal/app"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	appContext, cancelFunc := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		<-signalChan
		slog.Info("received interrupt signal")
		cancelFunc()
	}()
	err := app.ServerGrpc(appContext)
	if err != nil {
		slog.Error("failed to start grpc server", slog.String("err", err.Error()))
		return
	}
}

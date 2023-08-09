package main

import (
	"context"
	app2 "github.com/BAN1ce/Tree/app"
	"github.com/BAN1ce/Tree/logger"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		app         = app2.NewApp()
		ctx, cancel = context.WithCancel(context.TODO())
	)
	if err := app.StartTopicCluster(ctx); err != nil {
		log.Fatal(err)
	}
	defer cancel()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit
	cancel()
	time.Sleep(1 * time.Second)
	logger.Logger.Info("exist successfully")
	os.Exit(0)

}

package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"scrapper/internal/controller"
	"scrapper/internal/repository"
	"scrapper/internal/service"
	"syscall"
)

func main() {
	db, err := repository.NewSqliteDb()
	if err != nil {
		log.Fatal(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-sigs
		cancel()
	}()

	repo := repository.NewNodesRepository(db)
	parser := service.NewParserService()
	svc := service.NewNodesService(repo, parser)
	cli := controller.NewClient(svc)

	cli.StartParsing(ctx)

}

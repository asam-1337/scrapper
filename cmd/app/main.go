package main

import (
	"log"
	"scrapper/internal/controller"
	"scrapper/internal/repository"
	"scrapper/internal/service"
)

func main() {
	db, err := repository.NewSqliteDb()
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewNodesRepository(db)
	svc := service.NewNodesService(repo)
	client := controller.NewClient(svc)

	client.Parse("https://oidref.com/")
}
package main

import (
	"S3FileStorage/server"
	"S3FileStorage/server/config"
	"github.com/21Timofei/UI-Web-Interface"
)

func main() {
	cfg := config.LoadConfig()

	logger, err := config.ZapConfig().Build()
	if err != nil {
		panic(err)
	}

	tmpl, err := templates.LoadTemplates()
	if err != nil {
		logger.Fatal(err.Error())
	}

	srv := server.NewServer(cfg, tmpl, logger)
	srv.InitializeRoutes()
	srv.Start("8080")
}

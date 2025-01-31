package main

import (
	"S3FileStorage/server"
	"S3FileStorage/server/config"
	"html/template"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		log.Fatalf("Ошибка при загрузке шаблона: %v", err)
	}

	logger, _ := config.ZapConfig().Build()
	srv := server.NewServer(cfg, tmpl, logger)
	srv.InitializeRoutes()
	srv.Start("8080")
}

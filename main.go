package main

import (
	"html/template"
	"log"

	"S3FileStorage/server"
	"S3FileStorage/server/config"
)

func main() {
	cfg := config.LoadConfig()

	tmpl, err := template.ParseFiles("templates/template.html")
	if err != nil {
		log.Fatalf("Ошибка при загрузке шаблона: %v", err)
	}

	srv := server.NewServer(cfg, tmpl)
	srv.InitializeRoutes()
	srv.Start("8080")
}

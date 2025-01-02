package server

import (
	"S3FileStorage/server/config"
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	Config *config.Config
	Tmpl   *template.Template
	Mux    *http.ServeMux
}

func NewServer(cfg *config.Config, tmpl *template.Template) *Server {
	return &Server{
		Config: cfg,
		Tmpl:   tmpl,
		Mux:    http.NewServeMux(),
	}
}

func (s *Server) InitializeRoutes() {
	s.Mux.HandleFunc("/", s.ListFilesHandler)
	s.Mux.HandleFunc("/upload", s.UploadFileHandler)
	s.Mux.HandleFunc("/download", s.DownloadFileHandler)

	s.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func (s *Server) Start(port string) {
	log.Printf("Сервер запущен на http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, s.Mux))
}

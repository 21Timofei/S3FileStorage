package server

import (
	"S3FileStorage/server/config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
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

func (s *Server) ListFilesHandler(w http.ResponseWriter, _ *http.Request) {
	resp, err := s.Config.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.Config.Bucket),
	})
	if err != nil {
		log.Printf("Ошибка при получении списка файлов: %v", err)
		http.Error(w, "Ошибка при получении списка файлов", http.StatusInternalServerError)
		return
	}

	var files []string
	for _, item := range resp.Contents {
		files = append(files, *item.Key)
	}

	data := struct {
		Files []string
	}{
		Files: files,
	}

	if err = s.Tmpl.Execute(w, data); err != nil {
		log.Printf("Ошибка при рендеринге шаблона: %v", err)
		http.Error(w, "Ошибка при рендере шаблона", http.StatusInternalServerError)
	}
}

func (s *Server) UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Ошибка при получении файла", http.StatusBadRequest)
		return
	}
	defer func(file multipart.File) {
		err = file.Close()
		if err != nil {

		}
	}(file)

	fileName := strings.TrimSpace(handler.Filename)
	if fileName == "" {
		http.Error(w, "Некорректное имя файла", http.StatusBadRequest)
		return
	}

	_, err = s.Config.S3Client.PutObject(&s3.PutObjectInput{
		Bucket:        aws.String(s.Config.Bucket),
		Key:           aws.String(fileName),
		Body:          file,
		ContentLength: aws.Int64(handler.Size),
		ContentType:   aws.String(handler.Header.Get("Content-Type")),
	})
	if err != nil {
		http.Error(w, "Ошибка при загрузке файла", http.StatusInternalServerError)
		return
	}

	go http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "Не указано имя файла", http.StatusBadRequest)
		return
	}

	resp, err := s.Config.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {

		}
	}(resp.Body)

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", *resp.ContentType)

	if _, err = io.Copy(w, resp.Body); err != nil {
		http.Error(w, "Ошибка при скачивании файла", http.StatusInternalServerError)
	}
}

package server

import (
	"S3FileStorage/server/config"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type Server struct {
	Config *config.Config
	Tmpl   *template.Template
	Mux    *http.ServeMux
	Logger *zap.SugaredLogger
}

func NewServer(cfg *config.Config, tmpl *template.Template, logger *zap.Logger) *Server {
	return &Server{
		Config: cfg,
		Tmpl:   tmpl,
		Mux:    http.NewServeMux(),
		Logger: logger.Sugar(),
	}
}

func (s *Server) InitializeRoutes() {
	s.Mux.HandleFunc("/", s.ListFilesHandler)
	s.Mux.HandleFunc("/upload", s.UploadFileHandler)
	s.Mux.HandleFunc("/download", s.DownloadFileHandler)
	s.Mux.HandleFunc("/delete", s.DeleteFilesHandler)
	s.Mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func (s *Server) Start(port string) {
	s.Logger.Info("Сервер запущен на http://localhost:", port)
	s.Logger.Fatal(http.ListenAndServe(":"+port, s.Mux))
}

func (s *Server) ListFilesHandler(w http.ResponseWriter, _ *http.Request) {
	if s.Config.S3Client == nil {
		s.Logger.Fatal("S3 клиент не инициализирован!")
	}
	s.Logger.Infof("S3 подключен к бакету: %s", s.Config.Bucket)

	s.Logger.Infof("Запрашиваем список файлов из бакета: %s", s.Config.Bucket)
	resp, err := s.Config.S3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(s.Config.Bucket),
	})
	if err != nil {
		s.Logger.Errorf("Ошибка: %v", err)
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
		s.Logger.Info("Ошибка при рендеринге шаблона: %v", err)
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
	defer func() {
		if err := file.Close(); err != nil {
			s.Logger.Errorf("Ошибка закрытия файла: %v", err)
		}
	}()

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

	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

func (s *Server) DownloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("name")
	if fileName == "" {
		http.Error(w, "Не указано имя файла", http.StatusBadRequest)
		return
	}

	s.Logger.Infof("Запрос на скачивание файла: %s", fileName)

	resp, err := s.Config.S3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Config.Bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		s.Logger.Errorf("Ошибка при скачивании %s: %v", fileName, err)
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			s.Logger.Errorf("Ошибка закрытия файла %s: %v", fileName, err)
		}
	}()

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", *resp.ContentType)

	if _, err = io.Copy(w, resp.Body); err != nil {
		s.Logger.Errorf("Ошибка при передаче файла %s: %v", fileName, err)
		http.Error(w, "Ошибка при скачивании файла", http.StatusInternalServerError)
	}
}
func (s *Server) DeleteFilesHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Files []string `json:"files"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, file := range request.Files {
		_, err := s.Config.S3Client.DeleteObject(&s3.DeleteObjectInput{
			Bucket: aws.String(s.Config.Bucket),
			Key:    aws.String(file),
		})
		if err != nil {
			http.Error(w, "Failed to delete some files", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

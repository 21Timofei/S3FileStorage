# **S3Storage**

📦 **S3Storage** — это HTTP-сервер на **Go**, предназначенный для работы с объектным хранилищем **Amazon S3 / Yandex Object Storage**.  
Сервер позволяет загружать, скачивать, просматривать и удалять файлы через веб-интерфейс и API.

## **🚀 Функциональность**
- 📂 **Загрузка файлов** в S3 через веб-интерфейс или API
- 📥 **Скачивание файлов** из S3
- 📃 **Просмотр списка файлов**
- 🗑️ **Удаление файлов**
- 🔒 **Логирование через Uber Zap**
- 🏗 **Поддержка контейнеризации (Docker)**
- ⚡ **Асинхронная обработка HTTP-запросов**

---

## **🛠 Используемые технологии**
| Технология  | Описание  |
|------------|-----------|
| **Go** | Основной язык разработки |
| **Amazon S3 / Yandex Object Storage** | Хранение файлов |
| **AWS SDK for Go** | Взаимодействие с S3 API |
| **Yandex Go SDK** | Интеграция с Yandex Object Storage |
| **Uber Zap** | Логирование событий |
| **HTML + JavaScript** | Веб-интерфейс |
| **Docker** | Контейнеризация |

---

## **📥 Установка**
### **1️⃣ Клонирование репозитория**
```sh
git clone https://github.com/21Timofei/S3Storage.git
cd S3Storage
```

### **2️⃣ Установка зависимостей**
```sh
go mod tidy
```

### **3️⃣ Настройка переменных окружения**
Создайте `.env` файл и добавьте:
```ini
S3_BUCKET=your-bucket-name
S3_REGION=your-region
S3_ACCESS_KEY=your-access-key
S3_SECRET_KEY=your-secret-key
```

---

## **🚀 Запуск**
### **Локальный запуск**
```sh
go run main.go
```
Сервер запустится на `http://localhost:8080`.

### **Запуск в Docker**
```sh
docker build -t s3storage .
docker run -p 8080:8080 --env-file .env s3storage
```

---

## **📡 API Эндпоинты**
| Метод  | URL  | Описание |
|--------|------|----------|
| **GET** | `/` | Главная страница с веб-интерфейсом |
| **POST** | `/upload` | Загрузка файла |
| **GET** | `/download?name={file}` | Скачивание файла |
| **DELETE** | `/delete` | Удаление файла |

Пример запроса на удаление файла:
```sh
curl -X DELETE http://localhost:8080/delete -H "Content-Type: application/json" -d '{"files":["example.txt"]}'
```

---

## **🛠 Структура проекта**
```
S3Storage/
│── server/
│   ├── config/            # Конфигурация сервера
│   ├── handlers/          # HTTP-обработчики
│   ├── templates/         # HTML-шаблоны
│── static/                # Статические файлы (CSS, JS)
│── main.go                # Точка входа
│── Dockerfile             # Конфигурация для Docker
│── go.mod                 # Модуль Go
```

---

## **📝 To-Do**
- [ ] Поддержка аутентификации
- [ ] Версионность файлов
- [ ] Поддержка Google Cloud Storage
- [ ] Подключение CI/CD

---

## **🤝 Контакты**
📧 Email: [timverhos@gmail.com](mailto:timverhos@gmail.com)  
🟦 Linkedin: [Timofei Verkhososov](https://www.linkedin.com/feed/)


**✨ Если проект вам понравился — ставьте ⭐ на GitHub!** 🚀


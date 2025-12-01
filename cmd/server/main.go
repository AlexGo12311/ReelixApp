package main

import (
	"log"
	"net/http"
	"os"
	"reelixapp/internal/handlers"
	"reelixapp/internal/storage"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	// Создаем необходимые папки
	createUploadDirs()

	// Инициализируем хранилище
	videoStorage := storage.NewMemoryStorage()

	// Инициализируем обработчики
	videoHandlers := handlers.NewVideoHandlers(videoStorage)

	// Настраиваем маршрутизатор
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/videos", videoHandlers.GetVideos).Methods("GET")
	api.HandleFunc("/videos", videoHandlers.UploadVideo).Methods("POST")
	api.HandleFunc("/videos/{id}", videoHandlers.GetVideo).Methods("GET")
	api.HandleFunc("/test", videoHandlers.TestUpload).Methods("GET") // Тестовый маршрут

	// Статические файлы
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Главная страница
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	// Настраиваем CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(r)

	log.Println("Reelix сервер запущен на http://localhost:8080")
	log.Println("Откройте в браузере: http://localhost:8080")
	log.Println("Для теста API перейдите на: http://localhost:8080/api/v1/videos")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func createUploadDirs() {
	dirs := []string{
		"uploads/videos",
		"uploads/thumbnails",
		"web/static/css",
		"web/static/js",
		"web/static/images",
		"web/templates",
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Printf("Warning: Error creating directory %s: %v", dir, err)
		} else {
			log.Printf("Created directory: %s", dir)
		}
	}
}

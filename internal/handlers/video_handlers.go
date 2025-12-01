package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reelixapp/internal/models"
	"reelixapp/internal/storage"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type VideoHandlers struct {
	storage *storage.MemoryStorage
}

func NewVideoHandlers(storage *storage.MemoryStorage) *VideoHandlers {
	return &VideoHandlers{
		storage: storage,
	}
}

func (h *VideoHandlers) GetVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := h.storage.GetAllVideos()
	if err != nil {
		http.Error(w, "Error getting videos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func (h *VideoHandlers) GetVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["id"]

	video, err := h.storage.GetVideo(videoID)
	if err != nil {
		http.Error(w, "Error getting video", http.StatusInternalServerError)
		return
	}

	if video == nil {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}

	// Увеличиваем счетчик просмотров
	go h.storage.IncrementViews(videoID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}

func (h *VideoHandlers) UploadVideo(w http.ResponseWriter, r *http.Request) {
	// Парсим multipart форму (максимальный размер 500MB)
	err := r.ParseMultipartForm(500 << 20) // 500 MB
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Получаем файл
	file, header, err := r.FormFile("video")
	if err != nil {
		http.Error(w, "Error getting video file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем данные формы
	title := r.FormValue("title")
	description := r.FormValue("description")
	author := r.FormValue("author")

	if title == "" || author == "" {
		http.Error(w, "Title and author are required", http.StatusBadRequest)
		return
	}

	// Создаем уникальное имя файла
	filename := uuid.New().String() + filepath.Ext(header.Filename)
	filePath := filepath.Join("uploads", "videos", filename)

	// Создаем директорию если не существует
	err = os.MkdirAll("uploads/videos", 0755)
	if err != nil {
		http.Error(w, "Error creating directory: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем файл на диске
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error creating file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Копируем содержимое
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем видео объект
	video := models.Video{
		Title:       title,
		Description: description,
		FilePath:    filePath,
		Thumbnail:   "/static/images/default-thumbnail.jpg",
		Duration:    "0:00",
		Author:      author,
	}

	err = h.storage.CreateVideo(video)
	if err != nil {
		http.Error(w, "Error saving video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id":      video.ID,
		"message": "Video uploaded successfully",
	})
}

// Простая заглушка для тестирования
func (h *VideoHandlers) TestUpload(w http.ResponseWriter, r *http.Request) {
	video := models.Video{
		Title:       "Тестовое видео",
		Description: "Это тестовое видео",
		FilePath:    "uploads/videos/test.mp4",
		Thumbnail:   "/static/images/default-thumbnail.jpg",
		Duration:    "5:30",
		Author:      "Тестовый автор",
	}

	err := h.storage.CreateVideo(video)
	if err != nil {
		http.Error(w, "Error creating test video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Test video created",
		"id":      video.ID,
	})
}

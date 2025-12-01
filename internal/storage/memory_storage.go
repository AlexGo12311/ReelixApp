package storage

import (
	"reelixapp/internal/models"
	"sync"
	"time"

	"github.com/google/uuid"
)

type MemoryStorage struct {
	videos map[string]models.Video
	mu     sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		videos: make(map[string]models.Video),
	}
}

func (s *MemoryStorage) CreateVideo(video models.Video) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	video.ID = uuid.New().String()
	video.CreatedAt = time.Now()
	video.Views = 0
	video.Likes = 0
	video.Dislikes = 0

	s.videos[video.ID] = video
	return nil
}

func (s *MemoryStorage) GetVideo(id string) (*models.Video, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	video, exists := s.videos[id]
	if !exists {
		return nil, nil
	}

	return &video, nil
}

func (s *MemoryStorage) GetAllVideos() ([]models.Video, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	videos := make([]models.Video, 0, len(s.videos))
	for _, video := range s.videos {
		videos = append(videos, video)
	}

	// Сортируем по дате создания (новые сначала)
	for i, j := 0, len(videos)-1; i < j; i, j = i+1, j-1 {
		videos[i], videos[j] = videos[j], videos[i]
	}

	return videos, nil
}

func (s *MemoryStorage) IncrementViews(videoID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if video, exists := s.videos[videoID]; exists {
		video.Views++
		s.videos[videoID] = video
	}

	return nil
}

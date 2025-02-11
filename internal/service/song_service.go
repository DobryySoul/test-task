package service

import (
	"strings"

	"github.com/DobryySoul/test-task/internal/models"
)

type SongRepository interface {
	CreateSong(song *models.CreateSongInput) error
	GetByGroupAndSongName(group, songName string) (*models.Song, error)
	UpdateSong(song *models.Song, ID int) error
	UpdateFieldSong(updateField *models.UpdateSongInput, song *models.Song) error
	Delete(id int) error
	GetByID(id int) (*models.Song, error)
	GetAllSongs(filter models.SongFilter, pagination models.Pagination) ([]models.Song, int, error)
}

type SongService struct {
	repo SongRepository
}

func NewSongService(repo SongRepository) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) CreateSong(song *models.CreateSongInput) error {
	return s.repo.CreateSong(song)
}

func (s *SongService) GetByGroupAndSongName(group, songName string) (*models.Song, error) {
	return s.repo.GetByGroupAndSongName(group, songName)
}

func (s *SongService) UpdateSong(song *models.Song, ID int) error {
	return s.repo.UpdateSong(song, ID)
}

func (s *SongService) UpdateFieldSong(updateField *models.UpdateSongInput, song *models.Song) error {
	return s.repo.UpdateFieldSong(updateField, song)
}

func (s *SongService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *SongService) GetSongByID(id int) (*models.Song, error) {
	return s.repo.GetByID(id)
}

func (s *SongService) GetSongText(id int) ([]string, error) {
	song, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(song.Text, "\n")
	return verses, nil
}

func (s *SongService) GetAllSongs(filter models.SongFilter, pagination models.Pagination) ([]models.Song, int, error) {
	return s.repo.GetAllSongs(filter, pagination)
}

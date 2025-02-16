package service

import (
	"strings"

	"github.com/DobryySoul/test-task/internal/entity"
)

type Repository interface {
	CreateSong(song *entity.CreateSongInput) error
	GetByGroupAndSongName(group, songName string) (*entity.Song, error)
	// UpdateSong(song *entity.Song, ID int) error
	UpdateFieldSong(updateField *entity.UpdateSongInput, song *entity.Song) error
	Delete(id int) error
	GetByID(id int) (*entity.Song, error)
	GetAllSongs(filter entity.SongFilter, pagination entity.Pagination) ([]entity.Song, int, error)
}

type Service struct {
	repo Repository
}

func NewSongService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateSong(song *entity.CreateSongInput) error {
	return s.repo.CreateSong(song)
}

func (s *Service) GetByGroupAndSongName(group, songName string) (*entity.Song, error) {
	return s.repo.GetByGroupAndSongName(group, songName)
}

// func (s *Service) UpdateSong(song *entity.Song, ID int) error {
// 	return s.repo.UpdateSong(song, ID)
// }

func (s *Service) UpdateFieldSong(updateField *entity.UpdateSongInput, song *entity.Song) error {
	return s.repo.UpdateFieldSong(updateField, song)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *Service) GetSongByID(id int) (*entity.Song, error) {
	return s.repo.GetByID(id)
}

func (s *Service) GetSongText(id int) ([]string, error) {
	song, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	verses := strings.Split(song.SongText, "\n")
	return verses, nil
}

func (s *Service) GetAllSongs(filter entity.SongFilter, pagination entity.Pagination) ([]entity.Song, int, error) {
	return s.repo.GetAllSongs(filter, pagination)
}

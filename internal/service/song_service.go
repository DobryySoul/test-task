package service

import (
	"strings"

	"github.com/DobryySoul/test-task/internal/models"
	"github.com/DobryySoul/test-task/internal/repo/postgres"
)

type SongService struct {
	repo *postgres.SongRepository
}

func NewSongService(repo *postgres.SongRepository) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) CreateSong(song *models.Song) error {
	return s.repo.CreateSong(song)
}

func (s *SongService) GetByGroupAndSongName(group, songName string) (*models.Song, error) {
	return s.repo.GetByGroupAndSongName(group, songName)
}

func (s *SongService) UpdateSong(song *models.Song, ID int) error {
	return s.repo.UpdateSong(song, ID)
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
	return s.repo.GetAll(filter, pagination)
}

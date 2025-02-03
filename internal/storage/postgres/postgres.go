package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/DobryySoul/test-task/internal/models"
)

type SongRepository struct {
	db  *sql.DB
	log *logrus.Logger
}

func NewSongRepository(db *sql.DB, log *logrus.Logger) *SongRepository {
	return &SongRepository{
		db:  db,
		log: log,
	}
}

func (s *SongRepository) CreateSong(song *models.Song) error {
	stmt, err := s.db.Prepare("INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5) RETURNING id")
	if err != nil {
		return fmt.Errorf("failed to insert song: %v", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(song.Group, song.Song, song.ReleaseDate, song.Text, song.Link).Scan(&song.ID)
	if err != nil {
		return fmt.Errorf("failed to insert song: %v", err)
	}

	return nil
}

func (s *SongRepository) GetByID(id int) (*models.Song, error) {
	var song models.Song

	stmt, err := s.db.Prepare("SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		return nil, err
	}

	return &song, nil
}

func (s *SongRepository) UpdateSong(song *models.Song, ID int) error {
	stmt, err := s.db.Prepare("UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}

	song.ID = ID

	_, err = stmt.Exec(song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID)
	if err != nil {
		return fmt.Errorf("failed to update song: %v", err)
	}

	return nil
}

func (s *SongRepository) Delete(id int) error {
	stmt, err := s.db.Prepare("DELETE FROM songs WHERE id = $1")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %v", err)
	}

	return nil
}

func (s *SongRepository) GetAll(filter models.SongFilter, pagination models.Pagination) ([]models.Song, int, error) {
	var whereClauses []string
	var args []interface{}
	paramIdx := 1

	buildCondition := func(filterValue *string, column string, exactMatch bool) {
		if filterValue != nil && *filterValue != "" {
			if exactMatch {
				whereClauses = append(whereClauses, fmt.Sprintf("%s = $%d", column, paramIdx))
			} else {
				whereClauses = append(whereClauses, fmt.Sprintf("%s ILIKE $%d", column, paramIdx))
				*filterValue = "%" + *filterValue + "%"
			}
			args = append(args, *filterValue)
			paramIdx++
		}
	}

	buildCondition(filter.Group, "group_name", true)
	buildCondition(filter.Song, "song_name", true)
	buildCondition(filter.ReleaseDate, "release_date", true)
	buildCondition(filter.Text, "text", false)
	buildCondition(filter.Link, "link", true)

	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM songs" + where
	var total int
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		s.log.Errorf("Count query failed: %v", err)
		return nil, 0, err
	}

	mainQuery := fmt.Sprintf(`
        SELECT 
            id, 
            group_name, 
            song_name, 
            release_date, 
            text, 
            link
        FROM songs
        %s
		ORDER BY group_name, song_name, release_date, text, link
        LIMIT $%d OFFSET $%d`,
		where, paramIdx, paramIdx+1)

	args = append(args, pagination.Limit, (pagination.Page-1)*pagination.Limit)

	rows, err := s.db.Query(mainQuery, args...)
	if err != nil {
		s.log.Errorf("Main query failed: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var s models.Song
		err := rows.Scan(
			&s.ID,
			&s.Group,
			&s.Song,
			&s.ReleaseDate,
			&s.Text,
			&s.Link,
		)
		if err != nil {
			return nil, 0, err
		}
		songs = append(songs, s)
	}

	return songs, total, nil
}

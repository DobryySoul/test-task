package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/DobryySoul/test-task/internal/entity"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

var ErrNotFound = errors.New("record not found")

func (s *Repository) CreateSong(song *entity.CreateSongInput) error {

	query := `
		INSERT INTO Songs(song_name, release_date, song_text, link, artist_id)
		VALUES($1, $2, $3, $4, $5)
	`

	_, err := s.db.Exec(query, song.SongName, song.ReleaseDate, song.SongText, song.Link, song.ArtistID)
	if err != nil {
		return fmt.Errorf("возникла ошибка в добавлении песни: %w", err)
	}

	return nil
}

func (s *Repository) GetByGroupAndSongName(group, songName string) (*entity.Song, error) {
	const methodName = "GetByGroupAndSongName"

	var song entity.Song
	query := `SELECT s.song_id, s.song_name, s.release_date, s.song_text, s.link, s.artist_id
			  FROM Songs s
			  JOIN Artists a ON s.artist_id = a.artist_id
			  WHERE a.group_name = $1 AND s.song_name = $2`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка в подготовке stmt: %w", methodName, err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			log.Printf("%s: ошибка закрытия stmt: %v", methodName, err)
		}
	}()

	err = stmt.QueryRow(group, songName).Scan(
		&song.SongID,
		&song.SongName,
		&song.ReleaseDate,
		&song.SongText,
		&song.Link,
		&song.ArtistID,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", methodName, ErrNotFound)
		}

		return nil, fmt.Errorf("%s: ошибка при выполнении запроса: %w", methodName, err)
	}

	return &song, nil
}

func (s *Repository) GetByID(id int) (*entity.Song, error) {
	const methodName = "GetByID"

	var song entity.Song
	query := "SELECT song_id, song_name, release_date, song_text, link FROM Songs WHERE song_id = $1"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", methodName, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			log.Printf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
		}
	}()

	err = stmt.QueryRow(id).Scan(
		&song.SongID,
		&song.SongName,
		&song.ReleaseDate,
		&song.SongText,
		&song.Link,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", methodName, ErrNotFound)
		}

		return nil, fmt.Errorf("%s: %w", methodName, err)
	}

	return &song, nil
}

func (s *Repository) UpdateFieldSong(updateField *entity.UpdateSongInput, song *entity.Song) error {
	const methodName = "UpdateFieldSong"

	query := `UPDATE Songs
             SET song_name = $1,
                 release_date = $2,
                 song_text = $3,
                 link = $4
             WHERE song_id = $5`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: ошибка в подготовке stmt: %w", methodName, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			log.Printf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
		}
	}()

	switch {
	case updateField.Song != "":
		song.SongName = updateField.Song
	case updateField.ReleaseDate != "":
		song.ReleaseDate = updateField.ReleaseDate
	case updateField.Text != "":
		song.SongText = updateField.Text
	case updateField.Link != "":
		song.Link = updateField.Link
	}

	_, err = stmt.Exec(
		song.SongName,
		song.ReleaseDate,
		song.SongText,
		song.Link,
		song.SongID,
	)

	if err != nil {
		return fmt.Errorf("%s: ошибка при выполнении запроса: %w", methodName, err)
	}

	return nil
}

// func (s *Repository) UpdateSong(song *entity.Song, ID int) error {
// 	const methodName = "UpdateSong"

// 	query := `UPDATE songs
//              SET group_name = $1,
//                  song_name = $2,
//                  release_date = $3,
//                  text = $4,
//                  link = $5
//              WHERE id = $6`

// 	stmt, err := s.db.Prepare(query)
// 	if err != nil {
// 		return fmt.Errorf("%s: ошибка подготовки: %w", methodName, err)
// 	}
// 	defer func() {
// 		if closeErr := stmt.Close(); closeErr != nil {
// 			log.Printf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
// 		}
// 	}()

// 	song.ID = ID

// 	_, err = stmt.Exec(
// 		song.Group,
// 		song.Song,
// 		song.ReleaseDate,
// 		song.Text,
// 		song.Link,
// 		song.ID,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("%s: ошибка выполнения: %w", methodName, err)
// 	}

// 	return nil
// }

func (s *Repository) Delete(id int) error {
	const methodName = "Delete"

	query := "DELETE FROM Songs WHERE song_id = $1"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			log.Printf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
		}
	}()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", methodName, err)
	}

	return nil
}

func (s *Repository) GetAllSongs(filter entity.SongFilter, pagination entity.Pagination) ([]entity.Song, int, error) {
	const methodName = "GetAll"

	var whereClauses []string
	var args []interface{}
	paramIdx := 1

	buildCondition := func(filterValue *string, column string, exactMatch bool) {
		if filterValue != nil && *filterValue != "" {
			clause := ""
			if exactMatch {
				clause = fmt.Sprintf("%s = $%d", column, paramIdx)
			} else {
				clause = fmt.Sprintf("%s ILIKE $%d", column, paramIdx)
				*filterValue = "%" + *filterValue + "%"
			}
			whereClauses = append(whereClauses, clause)
			args = append(args, *filterValue)
			paramIdx++
		}
	}

	buildCondition(filter.Song, "song_name", true)
	buildCondition(filter.ReleaseDate, "release_date", true)
	buildCondition(filter.Text, "song_text", false)
	buildCondition(filter.Link, "link", true)

	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM Songs" + where

	var total int

	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", methodName, err)
	}

	mainQuery := fmt.Sprintf(`
        SELECT
            song_id,
            song_name,
            release_date,
            song_text,
            link,
			artist_id
        FROM songs
        %s
		ORDER BY song_name, release_date, song_text, link, artist_id
        LIMIT $%d OFFSET $%d`,
		where, paramIdx, paramIdx+1)

	args = append(args, pagination.Limit, (pagination.Page-1)*pagination.Limit)

	rows, err := s.db.Query(mainQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%s: %w", methodName, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("%s: ошибка закрытия rows: %v", methodName, err)
		}
	}()

	var songs []entity.Song

	for rows.Next() {
		var song entity.Song
		err := rows.Scan(
			&song.SongID,
			&song.SongName,
			&song.ReleaseDate,
			&song.SongText,
			&song.Link,
			&song.ArtistID,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("%s: %w", methodName, err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%s: %w", methodName, err)
	}

	return songs, total, nil
}

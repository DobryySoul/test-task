package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/DobryySoul/test-task/internal/models"
	"github.com/DobryySoul/test-task/pkg/logger"
)

type SongRepository struct {
	db     *sql.DB
	logger logger.Logger
}

func NewSongRepository(db *sql.DB, log logger.Logger) *SongRepository {
	return &SongRepository{
		db:     db,
		logger: log,
	}
}

var ErrNotFound = errors.New("record not found")

func (s *SongRepository) CreateSong(song *models.CreateSongInput) error {

	query := `
		INSERT INTO songs(group_name, song_name, release_date, text, link)
		VALUES($1, $2, $3, $4, $5) RETURNING id
	`

	s.logger.Debugf("начало выполнения запроса: %s", query)
	start := time.Now()

	var id int

	err := s.db.QueryRow(query,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
	).Scan(&id)

	if err != nil {
		s.logger.Errorf("ошибка в добавлении песни: %v", err)

		return fmt.Errorf("возникла ошибка в добавлении песни: %w", err)
	}

	s.logger.Infof("песня успешно добавлена. ID: %d. Время выполнения: %v", id, time.Since(start))

	return nil
}

func (s *SongRepository) GetByGroupAndSongName(group, songName string) (*models.Song, error) {
	const methodName = "GetByGroupAndSongName"

	startTime := time.Now()

	var song models.Song
	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE group_name = $1 AND song_name = $2"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		s.logger.Debugf("%s: ошибка в подготовке stmt: %v", methodName, err)

		return nil, fmt.Errorf("%s: ошибка в подготовке stmt: %w", methodName, err)
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Debugf("%s: ошибка закрытия stmt: %v", methodName, err)
		}
	}()

	s.logger.Debugf("%s: выполнение запроса с параметрами: group='%s', songName='%s'",
		methodName, group, songName)

	err = stmt.QueryRow(group, songName).Scan(
		&song.ID,
		&song.Group,
		&song.Song,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Errorf("%s: параметры не найдены - group: '%s', song: '%s'",
				methodName, group, songName)

			return nil, fmt.Errorf("%s: %w", methodName, ErrNotFound)
		}
		s.logger.Errorf("%s: неудалось выполнить запрос: %v", methodName, err)

		return nil, fmt.Errorf("%s: ошибка при выполнении запроса: %w", methodName, err)
	}

	s.logger.Infof("%s: успешное выполнение запроса - ID: %d, Group: %s, Sing: %s (Время выполнения: %v)",
		methodName,
		song.ID,
		song.Group,
		song.Song,
		time.Since(startTime))

	return &song, nil
}

func (s *SongRepository) GetByID(id int) (*models.Song, error) {
	const methodName = "GetByID"

	s.logger.Debugf("%s: началось выполнение запроса по ID: %d", methodName, id)
	startTime := time.Now()

	var song models.Song
	query := "SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id = $1"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		s.logger.Errorf("%s: ошибка в подготовке stmt: %v", methodName, err)

		return nil, fmt.Errorf("%s: %w", methodName, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			s.logger.Errorf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
		}
	}()

	s.logger.Debugf("%s: выполнение запроса для ID %d", methodName, id)

	err = stmt.QueryRow(id).Scan(
		&song.ID,
		&song.Group,
		&song.Song,
		&song.ReleaseDate,
		&song.Text,
		&song.Link,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Errorf("%s: песня с ID: %d не найдена", methodName, id)

			return nil, fmt.Errorf("%s: %w", methodName, ErrNotFound)
		}
		s.logger.Errorf("%s: не удалось выполнить запрос: %v", methodName, err)

		return nil, fmt.Errorf("%s: %w", methodName, err)
	}

	s.logger.Infof("%s: успешное выполнение запроса - ID: %d, Group: %s, Song: %s (took %v)",
		methodName,
		song.ID,
		song.Group,
		song.Song,
		time.Since(startTime))

	return &song, nil
}

func (s *SongRepository) UpdateFieldSong(updateField *models.UpdateSongInput, song *models.Song) error {
	const methodName = "UpdateFieldSong"

	s.logger.Debugf("%s: начало обновления поля песни по ID: %d", methodName, updateField)
	startTime := time.Now()

	query := `UPDATE songs 
             SET group_name = $1, 
                 song_name = $2, 
                 release_date = $3, 
                 text = $4, 
                 link = $5 
             WHERE id = $6`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		s.logger.Errorf("%s: ошибка в подготовке stmt: %v", methodName, err)

		return fmt.Errorf("%s: ошибка в подготовке stmt: %w", methodName, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			s.logger.Errorf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
		}
	}()

	s.logger.Debugf("%s: параметры обновления - Группа: '%s', Песня: '%s', Дата: %s, Ссылка: %s",
		methodName,
		updateField.Group,
		updateField.Song,
		updateField.ReleaseDate,
		updateField.Link,
	)

	switch {
	case updateField.Group != "":
		song.Group = updateField.Group
	case updateField.Song != "":
		song.Song = updateField.Song
	case updateField.ReleaseDate != "":
		song.ReleaseDate = updateField.ReleaseDate
	case updateField.Text != "":
		song.Text = updateField.Text
	case updateField.Link != "":
		song.Link = updateField.Link
	}

	_, err = stmt.Exec(
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.ID,
	)

	if err != nil {
		s.logger.Errorf("%s: не удалось обновить данные, запрос: %v", methodName, err)

		return fmt.Errorf("%s: ошибка при выполнении запроса: %w", methodName, err)
	}

	s.logger.Infof("%s: успешное выполнение запроса - ID: %d, Group: %s, Song: %s (took %v)",
		methodName,
		updateField.Group,
		updateField.Song,
		time.Since(startTime))

	return nil
}

func (s *SongRepository) UpdateSong(song *models.Song, ID int) error {
	const methodName = "UpdateSong"

	s.logger.Debugf("%s: начало обновления песни ID: %d", methodName, ID)
	startTime := time.Now()

	query := `UPDATE songs 
             SET group_name = $1, 
                 song_name = $2, 
                 release_date = $3, 
                 text = $4, 
                 link = $5 
             WHERE id = $6`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		s.logger.Errorf("%s: ошибка подготовки запроса: %v", methodName, err)

		return fmt.Errorf("%s: ошибка подготовки: %w", methodName, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			s.logger.Errorf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
		}
	}()

	song.ID = ID

	s.logger.Debugf("%s: параметры обновления - Группа: '%s', Песня: '%s', Дата: %s, Ссылка: %s",
		methodName,
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Link,
	)

	_, err = stmt.Exec(
		song.Group,
		song.Song,
		song.ReleaseDate,
		song.Text,
		song.Link,
		song.ID,
	)
	if err != nil {
		s.logger.Errorf("%s: ошибка выполнения запроса: %v", methodName, err)

		return fmt.Errorf("%s: ошибка выполнения: %w", methodName, err)
	}

	s.logger.Infof("%s: обновлённые данные - %+v, время выполнения: %v", methodName, song, time.Since(startTime))

	return nil
}

func (s *SongRepository) Delete(id int) error {
	const methodName = "Delete"

	s.logger.Debugf("%s: начало удаления песни ID: %d", methodName, id)
	startTime := time.Now()

	query := "DELETE FROM songs WHERE id = $1"

	stmt, err := s.db.Prepare(query)
	if err != nil {
		s.logger.Errorf("%s: ошибка подготовки запроса: %v", methodName, err)

		return fmt.Errorf("%s: %w", methodName, err)
	}
	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			s.logger.Errorf("%s: ошибка закрытия stmt: %v", methodName, closeErr)
		}
	}()

	s.logger.Debugf("%s: выполнение удаления для ID %d", methodName, id)

	_, err = stmt.Exec(id)
	if err != nil {
		s.logger.Errorf("%s: ошибка удаления по ID: %v", methodName, err)

		return fmt.Errorf("%s: %w", methodName, err)
	}

	s.logger.Infof("%s: успешное удаление - ID: %d, время выполнения: %v", methodName, id, time.Since(startTime))

	return nil
}

func (s *SongRepository) GetAll(filter models.SongFilter, pagination models.Pagination) ([]models.Song, int, error) {
	const methodName = "GetAll"

	startTime := time.Now()

	s.logger.Debugf("%s: начало выполнения. Фильтр: %+v, Пагинация: %+v",
		methodName, filter, pagination)

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
			s.logger.Debugf("%s: добавлено условие - %s", methodName, clause)
			paramIdx++
		}
	}

	s.logger.Debugf("%s: построение условий фильтрации", methodName)

	buildCondition(filter.Group, "group_name", true)
	buildCondition(filter.Song, "song_name", true)
	buildCondition(filter.ReleaseDate, "release_date", true)
	buildCondition(filter.Text, "text", false)
	buildCondition(filter.Link, "link", true)

	where := ""
	if len(whereClauses) > 0 {
		where = " WHERE " + strings.Join(whereClauses, " AND ")
		s.logger.Debugf("%s: итоговое условие WHERE: %s", methodName, where)
	}

	countQuery := "SELECT COUNT(*) FROM songs" + where

	var total int

	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		s.logger.Errorf("%s: ошибка функции агрегации (COUNT) запроса: %v", methodName, err)
		return nil, 0, fmt.Errorf("%s: %w", methodName, err)
	}
	s.logger.Infof("%s: найдено всего записей: %d", methodName, total)

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
		s.logger.Errorf("%s: ошибка основного запроса: %v", methodName, err)
		return nil, 0, fmt.Errorf("%s: %w", methodName, err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Errorf("%s: ошибка закрытия rows: %v", methodName, err)
		}
	}()

	var songs []models.Song

	for rows.Next() {
		var song models.Song
		err := rows.Scan(
			&song.ID,
			&song.Group,
			&song.Song,
			&song.ReleaseDate,
			&song.Text,
			&song.Link,
		)
		if err != nil {
			s.logger.Errorf("%s: ошибка сканирования строки: %v", methodName, err)

			return nil, 0, fmt.Errorf("%s: %w", methodName, err)
		}
		songs = append(songs, song)
	}

	if err = rows.Err(); err != nil {
		s.logger.Errorf("%s: ошибка при обработке результатов: %v", methodName, err)

		return nil, 0, fmt.Errorf("%s: %w", methodName, err)
	}

	s.logger.Infof("%s: успешно получено %d/%d записей (время выполнения: %v)",
		methodName,
		len(songs),
		total,
		time.Since(startTime))

	return songs, total, nil
}

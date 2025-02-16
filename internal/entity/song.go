package entity

type Artist struct {
	ArtistID  int    `json:"artist_id"`
	GroupName string `json:"group_name"`
}

// Song model info
// @Description Информация о песне
type Song struct {
	ArtistID    int    `json:"artist_id"`
	SongID      int    `json:"song_id"`
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date"`
	SongText    string `json:"song_text"`
	Link        string `json:"link"`
}

// CreateSongInput model info
// @Description Информация о песне после создания
type CreateSongInput struct {
	SongName    string `json:"song_name"`
	ReleaseDate string `json:"release_date"`
	SongText    string `json:"song_text"`
	Link        string `json:"link"`
	ArtistID    int    `json:"artist_id"`
}

// UpdateSongInput model info
// @Description Изменение параметра существующей песни
type UpdateSongInput struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// SongsResponse model info
// @Description Ответ со списком песен и пагинацией
type SongsResponse struct {
	Data       []Song `json:"data"`
	Page       int    `json:"page"`
	TotalPages int    `json:"total_pages"`
	TotalItems int    `json:"total_items"`
}

// CreateSongResponse model info
// @Description Ответ с информацией о песне
type GetSongResponse struct {
	SongName    string `json:"songName"`
	Group       string `json:"group"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// ErrorResponse model info
// @Description Ответ об ошибке
type ErrorResponse struct {
	Error string `json:"error"`
}

// SongFilter model info
// @Description Фильтр по параметрам
type SongFilter struct {
	Group       *string `form:"group"`
	Song        *string `form:"song"`
	ReleaseDate *string `form:"release_date"`
	Text        *string `form:"text"`
	Link        *string `form:"link"`
}

// Pagination model info
// @Description Пагинация
type Pagination struct {
	Page  int `form:"page" validate:"min=1"`
	Limit int `form:"limit" validate:"min=1,max=100"`
}

package models

// Song model info
// @Description Информация о песне
type Song struct {
	ID          int    `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type CreateSongInput struct {
	Group       string `json:"group" binding:"required`
	Song        string `json:"song" binding:"required`
	ReleaseDate string `json:"releaseDate" binding:"required`
	Text        string `json:"text" binding:"required`
	Link        string `json:"link" binding:"required`
}

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

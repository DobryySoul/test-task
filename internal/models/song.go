package models

type Song struct {
	ID          int    `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SongsResponse struct {
	Data       []Song `json:"data"`
	Page       int    `json:"page"`
	TotalPages int    `json:"total_pages"`
	TotalItems int    `json:"total_items"`
}

type SongFilter struct {
    Group       *string `form:"group"`
    Song        *string `form:"song"`
    ReleaseDate *string `form:"release_date"`
    Text        *string `form:"text"`
    Link        *string `form:"link"`
}

type Pagination struct {
    Page  int `form:"page" validate:"min=1"`
    Limit int `form:"limit" validate:"min=1,max=100"`
}


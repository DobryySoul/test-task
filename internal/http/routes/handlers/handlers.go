package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/DobryySoul/test-task/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Service interface {
	CreateSong(song *entity.CreateSongInput) error
	GetByGroupAndSongName(group, songName string) (*entity.Song, error)
	// UpdateSong(song *entity.Song, ID int) error
	UpdateFieldSong(updateField *entity.UpdateSongInput, song *entity.Song) error
	Delete(id int) error
	GetSongByID(id int) (*entity.Song, error)
	GetSongText(id int) ([]string, error)
	GetAllSongs(filter entity.SongFilter, pagination entity.Pagination) ([]entity.Song, int, error)
}

type Handler struct {
	service   Service
	validator validator.Validate
	// log       *logger.Logger // implement logger in handle func
}

func NewHandler(service Service, validator validator.Validate) *Handler {
	return &Handler{
		service:   service,
		validator: validator,
	}
}

func (h *Handler) CreateSong(c *gin.Context) {
	var song entity.CreateSongInput

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := h.service.CreateSong(&song); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"created data": song})
}

// Handler godoc
// @Summary Получить песню по группе и названию
// @Description Возвращает информацию о песне по группе и названию
// @Tags songs
// @Produce  json
// @Param group query string true "Название группы"
// @Param song query string true "Название песни"
// @Success 200 {object} entity.GetSongResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /info [get]
func (h *Handler) GetSongByQuery(c *gin.Context) {
	var response entity.GetSongResponse

	group := c.Query("group")
	songName := c.Query("song")

	if group == "" || songName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song parameters are required"})
		return
	}

	song, err := h.service.GetByGroupAndSongName(group, songName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}
	response = entity.GetSongResponse{
		SongName:    song.SongName,
		Group:       group,
		ReleaseDate: song.ReleaseDate,
		Text:        song.SongText,
		Link:        song.Link,
	}

	c.JSON(http.StatusOK, response)
}

// // Handler godoc
// // @Summary Обновить песню
// // @Description Обновляет существующую запись песни
// // @Tags songs
// // @Accept  json
// // @Produce  json
// // @Param id path int true "ID песни"
// // @Param song body entity.Song true "Обновленные данные песни"
// // @Success 200 {object} entity.Song
// // @Failure 400 {object} entity.ErrorResponse
// // @Failure 500 {object} entity.ErrorResponse
// // @Router /update-song/{id} [put]
// func (h *Handler) UpdateSong(c *gin.Context) {
// 	var song entity.Song

// 	ID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

// 		return
// 	}

// 	if err := c.ShouldBindJSON(&song); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

// 		return
// 	}

// 	err = h.service.UpdateSong(&song, ID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

// 		return
// 	}

// 	c.JSON(http.StatusOK, song)
// }

func (h *Handler) UpdateFieldSong(c *gin.Context) {
	var song *entity.Song

	group := c.Query("group")
	songName := c.Query("song_name")

	if group == "" || songName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song parameters are required"})

		return
	}

	song, err := h.service.GetByGroupAndSongName(group, songName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	var UpdateFieldSong *entity.UpdateSongInput

	if err := c.ShouldBindJSON(&UpdateFieldSong); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = h.service.UpdateFieldSong(UpdateFieldSong, song)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"update data song": song})
}

// Handler godoc
// @Summary Удалить песню
// @Description Удаляет запись песни по ID
// @Tags songs
// @Produce  json
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]string
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /delete-song/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	group := c.Query("group")
	songName := c.Query("song_name")

	if group == "" || songName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song parameters are required"})

		return
	}

	song, err := h.service.GetByGroupAndSongName(group, songName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	err = h.service.Delete(song.SongID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "song deleted", "song": song.SongName, "song_id": song.SongID})
}

// Handler godoc
// @Summary Получить текст песни
// @Description Возвращает текст песни с пагинацией
// @Tags songs
// @Produce  json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит элементов на странице" default(2)
// @Success 200 {object} map[string]interface{} "Пример ответа: {"song": "название", "verses": [...]}"
// @Failure 400 {object} entity.ErrorResponse
// @Failure 404 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /song-text/{id}/text [get]
func (h *Handler) GetSongText(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid song ID"})

		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "2"))

	if page < 1 || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pagination parameters"})

		return
	}

	text, err := h.service.GetSongText(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	totalText := len(text)
	start := (page - 1) * limit
	end := start + limit

	if start > totalText {
		start = totalText
	}
	if end > totalText {
		end = totalText
	}

	song, _ := h.service.GetSongByID(id)
	response := gin.H{
		"song": song.SongName,
		"text": text[start:end],
	}

	c.JSON(http.StatusOK, response)
}

// Handler godoc
// @Summary Получить список песен
// @Description Возвращает список песен с фильтрацией и пагинацией
// @Tags songs
// @Produce  json
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param release_date query string false "Фильтр по дате выпуска"
// @Param text query string false "Фильтр по тексту"
// @Param link query string false "Фильтр по ссылке"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит элементов на странице" default(10)
// @Success 200 {object} entity.SongsResponse
// @Failure 400 {object} entity.ErrorResponse
// @Failure 500 {object} entity.ErrorResponse
// @Router /songs-with-filter [get]
func (h *Handler) GetSongs(c *gin.Context) {
	var filter entity.SongFilter
	var pagination entity.Pagination

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid query parameters"})

		return
	}

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{Error: "Invalid pagination parameters"})

		return
	}

	if err := h.validator.Struct(pagination); err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorResponse{
			Error: "Validation error: " + err.Error(),
		})

		return
	}

	songs, totalItems, err := h.service.GetAllSongs(filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, entity.ErrorResponse{
			Error: "Failed to retrieve songs",
		})

		return
	}

	totalPages := totalItems / pagination.Limit
	if totalItems%pagination.Limit != 0 {
		totalPages++
	}

	response := entity.SongsResponse{
		Data:       songs,
		Page:       pagination.Page,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}

	c.JSON(http.StatusOK, response)
}

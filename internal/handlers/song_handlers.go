package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/DobryySoul/test-task/internal/models"
	"github.com/DobryySoul/test-task/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type SongHandler struct {
	service   *service.SongService
	validator *validator.Validate
}

func NewHandler(service *service.SongService) *SongHandler {
	return &SongHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreateSong добавляет новую песню в библиотеку.
// @Summary Добавить новую песню
// @Description Добавляет новую песню в библиотеку, используя данные из JSON.
// @Tags info
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные песни"
// @Success 201 {object} models.Song
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [post]
func (sh *SongHandler) CreateSong(c *gin.Context) {
	var song models.Song

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := sh.service.CreateSong(&song); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, song)
}

// GetSongByID возвращает информацию о песне по её ID.
// @Summary Получить песню по ID
// @Description Возвращает информацию о песне по её ID.
// @Tags info
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id} [get]
func (sh *SongHandler) GetSongByID(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	song, err := sh.service.GetSongByID(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, song)
}

func (sh *SongHandler) UpdateSong(c *gin.Context) {
	var song models.Song

	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = sh.service.UpdateSong(&song, ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, song)
}

func (sh *SongHandler) DeleteSong(c *gin.Context) {
	ID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	err = sh.service.Delete(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "song deleted"})
}

// GetSongText возвращает текст песни с пагинацией по куплетам.
// @Summary Получить текст песни
// @Description Возвращает текст песни с пагинацией по куплетам
// @Tags songs
// @Produce json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Количество куплетов на странице" default(2)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
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

	verses, err := h.service.GetSongText(id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	totalVerses := len(verses)
	start := (page - 1) * limit
	end := start + limit

	if start > totalVerses {
		start = totalVerses
	}
	if end > totalVerses {
		end = totalVerses
	}

	song, _ := h.service.GetSongByID(id)
	response := gin.H{
		// "id":          id,
		// "group":       song.Group,
		"song":   song.Song,
		"verses": verses[start:end],
		// "page":        page,
		// "limit":       limit,
		// "totalVerses": totalVerses,
	}

	c.JSON(http.StatusOK, response)
}

func (sh *SongHandler) GetSongs(c *gin.Context) {
	var filter models.SongFilter
	var pagination models.Pagination

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid query parameters"})
		return
	}

	if err := c.ShouldBindQuery(&pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Invalid pagination parameters"})
		return
	}

	if err := sh.validator.Struct(pagination); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Validation error: " + err.Error(),
		})
		return
	}

	songs, totalItems, err := sh.service.GetAllSongs(filter, pagination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to retrieve songs",
		})
		return
	}

	totalPages := totalItems / pagination.Limit
	if totalItems%pagination.Limit != 0 {
		totalPages++
	}

	response := models.SongsResponse{
		Data:       songs,
		Page:       pagination.Page,
		TotalPages: totalPages,
		TotalItems: totalItems,
	}

	c.JSON(http.StatusOK, response)
}

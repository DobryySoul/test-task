package routes

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/DobryySoul/test-task/internal/models"
	"github.com/DobryySoul/test-task/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	// swagger
	_ "github.com/DobryySoul/test-task/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SongHandler struct {
	service   *service.SongService
	validator *validator.Validate
}

func NewRouter(h *gin.Engine, serv *service.SongService) {
	r := SongHandler{service: serv, validator: validator.New()}

	h.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: log.Writer(),
		Formatter: func(param gin.LogFormatterParams) string {
			return fmt.Sprintf("[GIN] %s |%s %d %s| %s |%s %s %s %s %s %s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				param.StatusCodeColor(),
				param.StatusCode,
				param.ResetColor(),
				param.ClientIP,
				param.MethodColor(),
				param.Method,
				param.ResetColor(),
				param.Path,
				param.Request.UserAgent(),
				param.Path,
			)
		},
	}))

	h.Use(gin.Recovery())

	// swagger
	h.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
	h.GET("/songs-with-filter", r.GetSongs)
	// Получение текста песни с пагинацией по куплетам
	h.GET("/song-text/:id/text", r.GetSongText)
	// Удаление песни по ID
	h.DELETE("/delete-song/:id", r.DeleteSong)
	// Изменение данных песни
	h.PUT("/update-song/:id", r.UpdateSong)
	// Добавление новой песни в формате JSON
	h.POST("/create-song", r.CreateSong)
	// метод из АПИ
	h.GET("/info", r.GetSongByQuery)

}

// SongHandler godoc
// @Summary Создать новую песню
// @Description Создает новую запись песни
// @Tags songs
// @Accept  json
// @Produce  json
// @Param song body models.Song true "Данные песни"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /create-song [post]
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

	c.JSON(http.StatusOK, gin.H{"message": "song sucssesfully created"})
}

// SongHandler godoc
// @Summary Получить песню по группе и названию
// @Description Возвращает информацию о песне по группе и названию
// @Tags songs
// @Produce  json
// @Param group query string true "Название группы"
// @Param song query string true "Название песни"
// @Success 200 {object} models.GetSongResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /info [get]
func (sh *SongHandler) GetSongByQuery(c *gin.Context) {
	var response models.GetSongResponse

	group := c.Query("group")
	songName := c.Query("song")

	if group == "" || songName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "group and song parameters are required"})
		return
	}

	song, err := sh.service.GetByGroupAndSongName(group, songName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

		return
	}

	response.ReleaseDate = song.ReleaseDate
	response.Text = song.Text
	response.Link = song.Link

	c.JSON(http.StatusOK, response)
}

// SongHandler godoc
// @Summary Обновить песню
// @Description Обновляет существующую запись песни
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Обновленные данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /update-song/{id} [put]
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

// SongHandler godoc
// @Summary Удалить песню
// @Description Удаляет запись песни по ID
// @Tags songs
// @Produce  json
// @Param id path int true "ID песни"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /delete-song/{id} [delete]
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

// SongHandler godoc
// @Summary Получить текст песни
// @Description Возвращает текст песни с пагинацией
// @Tags songs
// @Produce  json
// @Param id path int true "ID песни"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит элементов на странице" default(2)
// @Success 200 {object} map[string]interface{} "Пример ответа: {"song": "название", "verses": [...]}"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /song-text/{id}/text [get]
func (sh *SongHandler) GetSongText(c *gin.Context) {
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

	text, err := sh.service.GetSongText(id)
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

	// song, _ := sh.service.GetSongByID(id)
	response := gin.H{
		"text": text[start:end],
	}

	c.JSON(http.StatusOK, response)
}

// SongHandler godoc
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
// @Success 200 {object} models.SongsResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs-with-filter [get]
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

package router

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	// swagger
	_ "github.com/DobryySoul/test-task/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler interface {
	CreateSong(c *gin.Context)
	GetSongByQuery(c *gin.Context)
	DeleteSong(c *gin.Context)
	UpdateFieldSong(c *gin.Context)
	GetSongText(c *gin.Context)
	GetSongs(c *gin.Context)
}

type Router struct {
	Router  *gin.Engine
	Handler Handler
}

func NewRouter(h Handler) *Router {
	r := gin.Default()
	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
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

	r.Use(gin.Recovery())

	// swagger
	r.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Добавление новой песни в формате JSON
	r.POST("/create-song", h.CreateSong)
	// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
	r.GET("/songs-with-filter", h.GetSongs)
	// метод для получения информации о песне по названию группы и песни
	r.GET("/info", h.GetSongByQuery)
	// Получение текста песни и назвыания с пагинацией по куплетам по ID
	r.GET("/song-text/:id/text", h.GetSongText)
	// Удаление песни по названию группы и песни
	r.DELETE("/delete-song", h.DeleteSong)
	// Частичное изменение данных песни по названию группы и песни
	r.PATCH("/update-song", h.UpdateFieldSong)
	// // Изменение данных песни
	// r.PUT("/update-song/:id", h.UpdateSong)

	return &Router{Router: r}
}

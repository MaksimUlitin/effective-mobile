package routes

import (
	_ "effectiveMobileTask/docs"
	"effectiveMobileTask/internal/controllers"
	"effectiveMobileTask/lib/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router godoc
// @title Song Management API
// @version 1.0
// @description API for managing song information
// @host localhost:8080
// @BasePath /
func Router() *gin.Engine {
	r := gin.Default()
	// Info endpoint
	// @Tags Songs
	// @Summary Add song information
	r.POST("/info", controllers.AddSongInfo)
	// Songs list endpoint
	// @Tags Songs
	// @Summary List songs
	r.GET("/songs", controllers.GetSongs)
	// Song text endpoint
	// @Tags Songs
	// @Summary Get song text
	r.GET("/songs/:id/text", controllers.GetSongText)
	// Update song endpoint
	// @Tags Songs
	// @Summary Update a song
	r.PUT("/songs/:id", controllers.UpdateSong)
	// Delete song endpoint
	// @Tags Songs
	// @Summary Delete a song
	r.DELETE("/songs/:id", controllers.DeleteSong)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info("docs documentation is available at http://localhost:8080/swagger/index.html")

	return r
}

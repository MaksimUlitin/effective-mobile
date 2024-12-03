package main

import (
	"effectiveMobileTask/config"
	"effectiveMobileTask/internal/controllers"
	"effectiveMobileTask/internal/routes"
	"effectiveMobileTask/internal/storage/database"
	"effectiveMobileTask/lib/logger"
	"github.com/gin-gonic/gin"
	"log"
	"log/slog"
	"net/http"
)

// @title Song Management API
// @version 1.0
// @description API for managing song information
// @host localhost:8080
// @BasePath /
func main() {

	config.LoadConfigEnv()
	logger.Info("environment variables loaded")

	db := database.DbConnect()
	logger.Info("database connect success")

	database.Migrate(db)
	logger.Info("database migrate success")

	go mockServer()
	logger.Info("mock server start success")
	log.Fatal(routes.Router().Run(":8080"))
}

func mockServer() {
	testRouter := gin.Default()

	testRouter.GET("/info", func(c *gin.Context) {
		group := c.Query("group")
		song := c.Query("song")

		if group == "" || song == "" {
			logger.Debug("group or song is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "missing parameters"})
			return
		}

		songDetail, err := controllers.GetSongDetailJSON(group, song)
		if err != nil {
			logger.Debug("get song detail fail", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "song not found"})
			return
		}

		logger.Info("request to info succeeded", slog.Any("group", group), slog.Any("song", song))
		c.JSON(http.StatusOK, songDetail)
	})

	if err := testRouter.Run(":8088"); err != nil {
		log.Fatal("error starting http server", err)
	}
}

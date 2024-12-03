package controllers

import (
	"effectiveMobileTask/internal/models"
	"effectiveMobileTask/internal/storage/database"
	"effectiveMobileTask/lib/logger"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type SongEnriched struct {
	ReleaseDate string `json:"release_date"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type songRequest struct {
	Group string `json:"group" example:"Muse"` // Пример значения
	Title string `json:"song" example:"Supermassive Black Hole"`
}

// AddSongInfo godoc
// @Summary Get song information
// @Description Retrieve song information by group and title
// @Tags Songs
// @Accept json
// @Produce json
// @Param request body songRequest true "Request Body"
// @Success 200 {object} models.SongDetail "Song details successfully retrieved"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /info [post]
func AddSongInfo(c *gin.Context) {
	// Param group query string true "Group Name"
	// Param title query string true "Song Title"
	var requestBody songRequest

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		logger.Error("invalid request body", slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid request body"})
		return
	}

	group := requestBody.Group
	songs := requestBody.Title

	db := database.DbConnect()

	var songDB models.Song
	if err := db.Where("\"group\" = ? AND song = ?", group, songs).First(&songDB).Error; err != nil {
		logger.Info("song not found", slog.Any("params", map[string]string{"group": group, "song": songs}))
		songDetail, boolReturn := GetSongDetailAPI(group, songs, c)
		if boolReturn {
			return
		}
		newSong := models.Song{
			Group:       group,
			Song:        songs,
			ReleaseDate: songDetail.ReleaseDate,
			Text:        songDetail.Text,
			Link:        songDetail.Link,
		}

		if err := db.Create(&newSong).Error; err != nil {
			logger.Error("failed to add new song", slog.Any("error", err), slog.Any("params", map[string]string{"group": group, "song": songs}))
			c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
			return
		}
		logger.Info("added new song", slog.Any("params", map[string]string{"group": group, "song": songs}))
		songDB = newSong

	}

	songDetail := models.SongDetail{
		ReleaseDate: songDB.ReleaseDate,
		Text:        songDB.Text,
		Link:        songDB.Link,
	}

	songEnrichFromJSON(&songDetail, group, songs)
	c.JSON(http.StatusOK, songDetail)
}

func GetSongDetailAPI(group, song string, c *gin.Context) (models.SongDetail, bool) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)

	urlAPI := fmt.Sprintf("http://localhost:8088/info?group=%s&song=%s", encodedGroup, encodedSong)
	resp, err := http.Get(urlAPI)
	if err != nil {
		return handleAPIError(c, "failed to get song detail", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return handleAPIError(c, "failed to get song detail with status code", resp.StatusCode)
	}

	var dataAPI models.SongDetail
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return handleAPIError(c, "failed to read song detail", err)
	}

	if err := json.Unmarshal(body, &dataAPI); err != nil {
		return handleAPIError(c, "failed to unmarshal song detail", err)
	}

	return dataAPI, false
}

func handleAPIError(c *gin.Context, message string, err interface{}) (models.SongDetail, bool) {
	logger.Error(message, slog.Any("error", err))
	c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
	return models.SongDetail{}, true
}

func GetSongDetailJSON(group, song string) (models.SongDetail, error) {
	jsonFileEnrich, err := os.Open("enrichInfoSong.json")
	if err != nil {
		logger.Error("could not open song enrichment file", slog.Any("error", err))
		return models.SongDetail{}, err
	}
	defer jsonFileEnrich.Close()

	jsonVal, err := io.ReadAll(jsonFileEnrich)

	if err != nil {
		logger.Error("failed to read song enrichment file", slog.Any("error", err))
		return models.SongDetail{}, err
	}

	var enrichmentData SongEnriched
	if err := json.Unmarshal(jsonVal, &enrichmentData); err != nil {
		logger.Error("failed to unmarshal song enrichment file", slog.Any("error", err))
		return models.SongDetail{}, err
	}

	if enrichmentData.Group == group && enrichmentData.Song == song {
		return models.SongDetail{
			ReleaseDate: enrichmentData.ReleaseDate,
			Text:        enrichmentData.Text,
			Link:        enrichmentData.Link,
		}, nil
	}

	return models.SongDetail{}, errors.New("invalid song enrichment")
}

func songEnrichFromJSON(songDetail *models.SongDetail, group, song string) {
	jsonFileEnrich, err := os.Open("enrichInfoSong.json")
	if err != nil {
		logger.Error("failed to open enrichInfoSong.json", slog.Any("error", err))
		return
	}
	defer jsonFileEnrich.Close()

	jsonVal, err := io.ReadAll(jsonFileEnrich)

	if err != nil {
		logger.Error("failed to read enrichInfoSong.json", slog.Any("error", err))
	}

	var enrichmentData SongEnriched
	if err := json.Unmarshal(jsonVal, &enrichmentData); err != nil {
		logger.Error("failed to unmarshal enrichInfoSong.json", slog.Any("error", err))

	} else {
		if enrichmentData.Group == group && enrichmentData.Song == song {
			songDetail.ReleaseDate = enrichmentData.ReleaseDate
			songDetail.Text = enrichmentData.Text
			songDetail.Link = enrichmentData.Link

		}
	}
}

// GetSongs godoc
// @Summary List songs with optional filtering
// @Description Retrieve a list of songs with optional filtering and pagination
// @Tags Songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by Group Name"
// @Param song query string false "Filter by Song Title"
// @Param releaseDate query string false "Filter by Release Date"
// @Param text query string false "Filter by Text"
// @Param link query string false "Filter by Link"
// @Param page query int false "Page number for pagination" default(1)
// @Param limit query int false "Number of items per page" default(10)
// @Success 200 {array} models.Song "Songs retrieved successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /songs [get]
func GetSongs(c *gin.Context) {
	db := database.DbConnect()
	var songs []models.Song
	group := c.Query("group")
	song := c.Query("song")
	releaseDate := c.Query("release_date")
	text := c.Query("text")
	link := c.Query("link")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNumber, err := strconv.Atoi(page)

	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	limitNumber, err := strconv.Atoi(limit)

	if err != nil || limitNumber < 1 {
		limitNumber = 10
	}

	query := db.Model(&models.Song{})
	if group != "" {
		query = query.Where("\"group\" ILIKE ?", "%"+group+"%")
	}

	if song != "" {
		query = query.Where("\"song\" ILIKE ?", "%"+song+"%")
	}
	if releaseDate != "" {
		query = query.Where("release_date = ?", releaseDate)
	}
	if text != "" {
		query = query.Where("text ILIKE ?", "%"+text+"%")
	}
	if link != "" {
		query = query.Where("link ILIKE ?", "%"+link+"%")
	}

	offset := (pageNumber - 1) * limitNumber
	query = query.Offset(offset).Limit(limitNumber)
	if err := query.Find(&songs).Error; err != nil {
		logger.Error("failed to query songs", slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	logger.Info("Songs retrieved successfully", slog.Int("count", len(songs)))
	c.JSON(http.StatusOK, songs)

}

// GetSongText godoc
// @Summary Get song text by ID with pagination
// @Description Retrieve song text for a specific song ID with pagination support
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param page query int false "Page number for text pagination" default(1)
// @Param limit query int false "Number of text lines per page" default(10)
// @Success 200 {object} map[string]interface{} "Song text retrieved successfully"
// @Failure 404 {object} map[string]string "Song or page not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /songs/{id}/text [get]
func GetSongText(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Error("failed to get song text", slog.Any("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	db := database.DbConnect()
	var song models.Song
	if err := db.Unscoped().First(&song, id).Error; err != nil {
		logger.Error("failed to query song", slog.Any("id", id))
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}
	text := strings.Split(song.Text, "\n\n")

	totalText := len(text)
	if totalText == 0 {
		logger.Error("no text found for song with id ", slog.Any("id", id))
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	beginOfIndex := (page - 1) * limit
	endOfIndex := beginOfIndex + limit

	if beginOfIndex >= totalText {
		logger.Error("page out of range for song id", slog.Any("id", id), slog.Any("page", page))
		c.JSON(http.StatusNotFound, gin.H{"message": "no text found for req page"})
		return
	}
	if endOfIndex > totalText {
		endOfIndex = totalText
	}
	selectText := text[beginOfIndex:endOfIndex]
	resp := map[string]interface{}{
		"songId":    id,
		"page":      page,
		"text":      selectText,
		"total":     totalText,
		"limit":     limit,
		"totalPage": (totalText + limit - 1) / limit,
	}
	logger.Info("retrieved text for song id ", slog.Any("id", id), slog.Any("page", page))
	c.JSON(http.StatusOK, resp)
}

// UpdateSong godoc
// @Summary Update an existing song
// @Description Update song information by ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Song Update Information"
// @Success 200 {object} map[string]string "Song updated successfully"
// @Failure 400 {object} map[string]string "Invalid song data"
// @Failure 404 {object} map[string]string "Song not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	db := database.DbConnect()
	if err := db.First(&song, id).Error; err != nil {
		logger.Error("song not found", slog.Any("id", id))
		c.JSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	}
	if err := c.ShouldBind(&song); err != nil {
		logger.Error("invalid song data", slog.Any("id", id))
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid song data"})
		return
	}
	if err := db.Save(&song).Error; err != nil {
		logger.Error("failed to update song", slog.Any("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	logger.Info("song updated successfully", slog.Any("id", id))
	c.JSON(http.StatusOK, gin.H{"message": "song updated successfully"})
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song by its ID
// @Tags Songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]interface{} "Song deleted successfully"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	id := c.Param("id")
	db := database.DbConnect()
	if err := db.Delete(&models.Song{}, id).Error; err != nil {
		logger.Error("failed to delete song", slog.Any("id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}
	logger.Info("song deleted successfully", slog.Any("id", id))
	c.JSON(http.StatusOK, map[string]interface{}{"id #" + id: "deleted"})

}

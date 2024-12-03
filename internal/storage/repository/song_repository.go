package repository

import (
	"effectiveMobileTask/internal/models"
	"effectiveMobileTask/lib/logger"
	"gorm.io/gorm"
	"log/slog"
)

type SongRepository struct {
	DB *gorm.DB
}

func NewSongRepository(db *gorm.DB) *SongRepository {
	logger.Info("create song repository")
	return &SongRepository{DB: db}
}

func (repo *SongRepository) SaveSong(song *models.Song) (*models.Song, error) {
	if err := repo.DB.Create(song).Error; err != nil {
		return nil, err
	}
	logger.Info("save song successfully", slog.Any("song", song))
	return song, nil
}

func (repo *SongRepository) GetAllSongs(page, limit int) ([]models.Song, error) {
	logger.Info("retrieving all songs. page: %d, Limit: %d\\n\", page, limit")
	var songs []models.Song
	offset := (page - 1) * limit

	if err := repo.DB.Limit(limit).Offset(offset).Find(&songs).Error; err != nil {
		logger.Error("failed to retrieve songs", slog.Any("page", page), slog.Any("limit", limit), slog.Any("error", err))
		return nil, err
	}
	logger.Info("retrieved songs", slog.Any("songs", songs))
	return songs, nil
}

func (repo *SongRepository) GetSongById(id uint) (*models.Song, error) {
	logger.Info("retrieving song", slog.Any("id", id))
	var song models.Song
	if err := repo.DB.First(&song, id).Error; err != nil {
		logger.Error("failed to retrieve song ", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}
	logger.Info("retrieved song", slog.Any("id", song))
	return &song, nil
}

func (repo *SongRepository) UpdateSong(song *models.Song) (*models.Song, error) {
	logger.Info("updating song", slog.Any("song", song))
	if err := repo.DB.Save(song).Error; err != nil {
		logger.Error("failed to update song ", slog.Any("song", song), slog.Any("error", err))
		return nil, err

	}
	logger.Info("updated song", slog.Any("song", song))
	return song, nil
}

func (repo *SongRepository) DeletedSong(id uint) error {
	logger.Info("deleting song", slog.Any("id", id))
	if err := repo.DB.Delete(&models.Song{}, id).Error; err != nil {
		logger.Error("failed to delete song", slog.Any("id", id), slog.Any("error", err))
		return err
	}
	logger.Info("deleted song", slog.Any("id", id))
	return nil
}

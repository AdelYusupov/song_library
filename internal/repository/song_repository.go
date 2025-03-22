package repository

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"music-api/internal/model"
)

type SongRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewSongRepository(db *gorm.DB, logger *logrus.Logger) *SongRepository {
	return &SongRepository{db: db, logger: logger}
}

func (r *SongRepository) Create(song *model.Song) error {
	return r.db.Create(song).Error
}

func (r *SongRepository) GetByID(id uint) (*model.Song, error) {
	var song model.Song
	err := r.db.First(&song, id).Error
	return &song, err
}

func (r *SongRepository) Update(song *model.Song) error {
	return r.db.Save(song).Error
}

func (r *SongRepository) Delete(id uint) error {
	return r.db.Delete(&model.Song{}, id).Error
}

func (r *SongRepository) GetAll(filter map[string]interface{}, limit, offset int) ([]model.Song, error) {
	var songs []model.Song
	query := r.db.Model(&model.Song{})
	for key, value := range filter {
		query = query.Where(key, value)
	}
	err := query.Limit(limit).Offset(offset).Find(&songs).Error
	return songs, err
}

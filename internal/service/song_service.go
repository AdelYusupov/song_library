package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"music-api/internal/model"
	"music-api/internal/repository"
	"net/http"
	"time"
)

type SongService struct {
	repo   *repository.SongRepository
	logger *logrus.Logger
}

func NewSongService(repo *repository.SongRepository, logger *logrus.Logger) *SongService {
	return &SongService{repo: repo, logger: logger}
}

func (s *SongService) AddSong(song *model.Song) error {
	// Запрос к внешнему API для получения дополнительной информации
	url := fmt.Sprintf("%s/info?group=%s&song=%s", s.apiURL, song.Group, song.Title)
	resp, err := http.Get(url)
	if err != nil {
		s.logger.Errorf("Failed to fetch song info: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Errorf("API returned non-200 status code: %d", resp.StatusCode)
		return errors.New("failed to fetch song info")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Errorf("Failed to read response body: %v", err)
		return err
	}

	var songDetail struct {
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}
	if err := json.Unmarshal(body, &songDetail); err != nil {
		s.logger.Errorf("Failed to unmarshal song detail: %v", err)
		return err
	}

	releaseDate, err := time.Parse("02.01.2006", songDetail.ReleaseDate)
	if err != nil {
		s.logger.Errorf("Failed to parse release date: %v", err)
		return err
	}

	song.ReleaseDate = releaseDate
	song.Text = songDetail.Text
	song.Link = songDetail.Link

	return s.repo.Create(song)
}

func (s *SongService) GetSongs(filter map[string]interface{}, limit, offset int) ([]model.Song, error) {
	return s.repo.GetAll(filter, limit, offset)
}

func (s *SongService) GetSongText(id uint) (string, error) {
	song, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	return song.Text, nil
}

func (s *SongService) UpdateSong(song *model.Song) error {
	return s.repo.Update(song)
}

func (s *SongService) DeleteSong(id uint) error {
	return s.repo.Delete(id)
}

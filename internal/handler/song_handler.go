package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"music-api/internal/model"
	"music-api/internal/service"
	"net/http"
	"strconv"
)

type SongHandler struct {
	service *service.SongService
	logger  *logrus.Logger
}

func NewSongHandler(service *service.SongService, logger *logrus.Logger) *SongHandler {
	return &SongHandler{service: service, logger: logger}
}

// @Summary Get all songs
// @Description Get a list of songs with optional filtering and pagination
// @Tags songs
// @Accept json
// @Produce json
// @Param group query string false "Filter by group"
// @Param title query string false "Filter by title"
// @Param limit query int false "Limit the number of results"
// @Param offset query int false "Offset for pagination"
// @Success 200 {array} model.Song
// @Failure 500 {object} map[string]string
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	filter := make(map[string]interface{})
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			filter[key] = values[0]
		}
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	songs, err := h.service.GetSongs(filter, limit, offset)
	if err != nil {
		h.logger.Errorf("Failed to get songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get song text
// @Description Get the text of a song by ID with pagination by verses
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	text, err := h.service.GetSongText(uint(id))
	if err != nil {
		h.logger.Errorf("Failed to get song text: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get song text"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"text": text})
}

// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	if err := h.service.DeleteSong(uint(id)); err != nil {
		h.logger.Errorf("Failed to delete song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Update a song
// @Description Update song details by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body model.Song true "Song data"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.logger.Errorf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	song.ID = uint(id)
	if err := h.service.UpdateSong(&song); err != nil {
		h.logger.Errorf("Failed to update song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Add a new song
// @Description Add a new song to the library
// @Tags songs
// @Accept json
// @Produce json
// @Param song body model.Song true "Song data"
// @Success 201 {object} model.Song
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var song model.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.logger.Errorf("Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.AddSong(&song); err != nil {
		h.logger.Errorf("Failed to add song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add song"})
		return
	}

	c.JSON(http.StatusCreated, song)
}

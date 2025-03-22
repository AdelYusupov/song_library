package model

import "time"

type Song struct {
	ID          uint      `gorm:"primary_key"`
	Band        string    `gorm:"not null"` // Используем "Band" вместо "Group"
	Title       string    `gorm:"not null"`
	ReleaseDate time.Time `gorm:"not null"`
	Text        string    `gorm:"not null"`
	Link        string    `gorm:"not null"`
}

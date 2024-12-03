package models

import "time"

type Song struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Group       string     `json:"group"`
	Song        string     `json:"song"`
	ReleaseDate string     `json:"release_date"`
	Text        string     `json:"text"`
	Link        string     `json:"link"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}

type SongDetail struct {
	Text        string `json:"text"`
	Link        string `json:"link"`
	ReleaseDate string `json:"release_date"`
}

type AddNewSong struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

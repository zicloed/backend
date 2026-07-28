package models

import (
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Name        string `json:"name" binding:"required"`
	Description string `json:"name" binding:"required"`
	Location    string `json:"name" binding:"required"`
	UserId      int    `json: "UserId"`
}

var events []Event = []Event{}

// Fungsi Save event
func (e Event) Save() {
	events = append(events, e)
}

// fungsi show semua event
func GetAllEvent() []Event {
	return events
}

package db

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Website struct {
	Name        string   `json:"name"`
	Id          string   `json:"id"`
	IconURL     string   `json:"icon_url"`
	Description string   `json:"description"`
	Tags        pq.StringArray `json:"tags" gorm:"type:text[]"`
	Owner int `json:"owner"`
	Created     int `json:"created"`
  LastBumped  int `json:"last_bumped"`
	Bumps       int `json:"bumps"`
  Upvotes int `json:"upvotes"`
  UpvotesToday int `json:"upvotes_today"`
}

func Bump(id string, db gorm.DB) bool {
	website := Website{}

	db.Where("id = ?", id).First(&website)

	if len(website.Id) == 0 {
		return false
	}

	website.LastBumped = int(time.Now().Unix())
	website.Bumps = website.Bumps + 1
	
	db.Save(&website)

	return true
}


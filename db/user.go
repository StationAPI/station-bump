package db

import "gorm.io/gorm"

type User struct {
	GithubId int
	Bumps int
} 

func GetUser(id int, db gorm.DB) (bool, User) {
	user := User{}

	db.Where("id = ?", id).First(&user)

	if user.GithubId <= 0 {
		return false, User{}
	}

	return true, user
}

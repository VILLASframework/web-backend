package common

import (
//"github.com/jinzhu/gorm"
//"github.com/jinzhu/gorm/dialects/postgres"
)

func FindAllUsers() ([]User, int, error) {
	db := GetDB()
	var users []User
	err := db.Find(&users).Error
	return users, len(users), err
}

func FindUserProjects(user *User) ([]Project, int, error) {
	db := GetDB()
	var projects []Project
	err := db.Model(user).Related(&projects, "Projects").Error
	return projects, len(projects), err
}

func FindUserSimulations(user *User) ([]Simulation, int, error) {
	db := GetDB()
	var simulations []Simulation
	err := db.Model(user).Related(&simulations, "Simulations").Error
	return simulations, len(simulations), err
}

func FindUserFiles(user *User) ([]File, int, error) {
	db := GetDB()
	var files []File
	err := db.Model(user).Related(&files, "Files").Error
	return files, len(files), err
}

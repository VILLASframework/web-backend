package user

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllUsers() ([]common.User, int, error) {
	db := common.GetDB()
	var users []common.User
	err := db.Find(&users).Error
	return users, len(users), err
}

func FindUserProjects(user *common.User) ([]common.Project, int, error) {
	db := common.GetDB()
	var projects []common.Project
	err := db.Model(user).Related(&projects, "Projects").Error
	return projects, len(projects), err
}

func FindUserSimulations(user *common.User) ([]common.Simulation, int, error) {
	db := common.GetDB()
	var simulations []common.Simulation
	err := db.Model(user).Related(&simulations, "Simulations").Error
	return simulations, len(simulations), err
}

func FindUserFiles(user *common.User) ([]common.File, int, error) {
	db := common.GetDB()
	var files []common.File
	err := db.Model(user).Related(&files, "Files").Error
	return files, len(files), err
}

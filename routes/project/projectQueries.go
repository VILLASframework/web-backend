package project

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllProjects() ([]common.Project, int, error) {
	db := common.GetDB()
	var projects []common.Project
	err := db.Find(&projects).Error
	return projects, len(projects), err
}

func FindUserProjects(user *common.User) ([]common.Project, int, error) {
	db := common.GetDB()
	var projects []common.Project
	err := db.Model(user).Related(&projects, "Projects").Error
	return projects, len(projects), err
}

func FindVisualizationProject(visualization *common.Visualization) (common.Project, int, error) {
	db := common.GetDB()
	var project common.Project
	err := db.Model(visualization).Related(&project, "Projects").Error
	return project, 1, err
}



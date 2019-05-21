package queries

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllUsers() ([]common.User, int, error) {
	db := common.GetDB()
	var users []common.User
	err := db.Find(&users).Error
	return users, len(users), err
}

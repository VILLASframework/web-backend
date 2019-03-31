package user

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/file"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/project"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/simulation"
	"github.com/jinzhu/gorm"
	// TODO: we need also bcrypt
)

type User struct {
	gorm.Model
	Username    string `gorm:"unique"`
	Password    string
	Mail        string      `gorm:"default:"`
	Role        string      `gorm:"default:user"`
	Projects    []Project   // TODO: association & foreign key
	Simulations []Simulaion // TODO: association & foreign key
	Files       []File      // TODO: association & foreign key
}

// TODO: callback for verifying password

// TODO: execute before each user.save()

// TODO: execute before each user.delete()

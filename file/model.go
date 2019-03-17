package file

import (
	"github.com/jinzhu/gorm"
)

type File struct {
	gorm.Model
	Name        string
	Path        string
	Type        string
	Size        uint
	Dimmensions string
	FileUser    User
}

package file

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/user"
	"github.com/jinzhu/gorm"
	"time"
)

type File struct {
	gorm.Model
	Name        string
	Path        string `gorm:"not null"`
	Type        string
	Size        uint
	Dimmensions string // TODO: Mixed Type
	FileUser    User   `gorm:"not null"`
	Date        Time   `gorm:"default:Time.Now"`
}

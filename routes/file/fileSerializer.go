package file

import (
	"time"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type FilesSerializer struct {
	Ctx   *gin.Context
	Files []common.File
}

func (self *FilesSerializer) Response() []FileResponse {
	response := []FileResponse{}
	for _, File := range self.Files {
		serializer := FileSerializer{self.Ctx, File}
		response = append(response, serializer.Response())
	}
	return response
}

type FileSerializer struct {
	Ctx *gin.Context
	common.File
}

type FileResponse struct {
	Name string `json:"Name"`
	ID   uint   `json:"FileID"`
	Path string `json:"Path"`
	Type string `json:"Type"` //MIME type?
	Size uint   `json:"Size"`
	H    uint   `json:"ImageHeight"`
	W    uint   `json:"ImageWidth"`
	Date time.Time `json:"Date"`
}

func (self *FileSerializer) Response() FileResponse {

	response := FileResponse{
		Name:    self.Name,
		Path:    self.Path,
		Type:    self.Type,
		Size:    self.Size,
		Date:    self.Date,
		H: 	self.ImageHeight,
		W:    self.ImageWidth,
	}
	return response
}

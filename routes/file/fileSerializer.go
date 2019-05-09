package file

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

// File/s Serializers

type FilesSerializerNoAssoc struct {
	Ctx   *gin.Context
	Files []common.File
}

func (self *FilesSerializerNoAssoc) Response() []FileResponseNoAssoc {
	response := []FileResponseNoAssoc{}
	for _, files := range self.Files {
		serializer := FileSerializerNoAssoc{self.Ctx, files}
		response = append(response, serializer.Response())
	}
	return response
}

type FileSerializerNoAssoc struct {
	Ctx *gin.Context
	common.File
}

type FileResponseNoAssoc struct {
	Name string `json:"Name"`
	ID   uint   `json:"FileID"`
	Path string `json:"Path"`
	Type string `json:"Type"`
	Size uint   `json:"Size"`
	H    uint   `json:"ImageHeight"`
	W    uint   `json:"ImageWidth"`
	// Date
}

func (self *FileSerializerNoAssoc) Response() FileResponseNoAssoc {
	response := FileResponseNoAssoc{
		Name: self.Name,
		ID:   self.ID,
		Path: self.Path,
		Type: self.Type,
		Size: self.Size,
		H:    self.ImageHeight,
		W:    self.ImageWidth,
		// Date
	}
	return response
}
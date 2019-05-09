package project

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type ProjectsSerializerNoAssoc struct {
	Ctx      *gin.Context
	Projects []common.Project
}

func (self *ProjectsSerializerNoAssoc) Response() []ProjectResponseNoAssoc {
	response := []ProjectResponseNoAssoc{}
	for _, project := range self.Projects {
		serializer := ProjectSerializerNoAssoc{self.Ctx, project}
		response = append(response, serializer.Response())
	}
	return response
}

type ProjectSerializerNoAssoc struct {
	Ctx *gin.Context
	common.Project
}

type ProjectResponseNoAssoc struct {
	Name string `json:"Name"`
	ID   uint   `json:"ProjectID"`
}

func (self *ProjectSerializerNoAssoc) Response() ProjectResponseNoAssoc {
	response := ProjectResponseNoAssoc{
		Name: self.Name,
		ID:   self.ID,
	}
	return response
}



package common

import (
	"github.com/gin-gonic/gin"
)

type UsersSerializer struct {
	Ctx   *gin.Context
	Users []User
}

func (self *UsersSerializer) Response() []UserResponse {
	response := []UserResponse{}
	for _, user := range self.Users {
		serializer := UserSerializer{self.Ctx, user}
		response = append(response, serializer.Response())
	}
	return response
}

type UserSerializer struct {
	Ctx *gin.Context
	User
}

type UserResponse struct {
	Username    string `json:"Username"`
	Password    string `json:"Password"` // XXX: ???
	Role        string `json:"Role"`
	Mail        string `json:"Mail"`
	Projects    []ProjectResponseNoAssoc
	Simulations []SimulationResponseNoAssoc
	Files       []FileResponseNoAssoc
}

func (self *UserSerializer) Response() UserResponse {
	// TODO: maybe all those should be made in one transaction
	projects, _, _ := FindUserProjects(&self.User)
	projectsSerializer := ProjectsSerializerNoAssoc{self.Ctx, projects}

	simulations, _, _ := FindUserSimulations(&self.User)
	simulationsSerializer := SimulationsSerializerNoAssoc{self.Ctx, simulations}

	files, _, _ := FindUserFiles(&self.User)
	filesSerializer := FilesSerializerNoAssoc{self.Ctx, files}

	response := UserResponse{
		Username:    self.Username,
		Password:    self.Password,
		Role:        self.Role,
		Mail:        self.Mail,
		Projects:    projectsSerializer.Response(),
		Simulations: simulationsSerializer.Response(),
		Files:       filesSerializer.Response(),
	}
	return response
}

// Project/s Serializers

type ProjectsSerializerNoAssoc struct {
	Ctx      *gin.Context
	Projects []Project
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
	Project
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

// Simulation/s Serializers

type SimulationsSerializerNoAssoc struct {
	Ctx         *gin.Context
	Simulations []Simulation
}

func (self *SimulationsSerializerNoAssoc) Response() []SimulationResponseNoAssoc {
	response := []SimulationResponseNoAssoc{}
	for _, simulation := range self.Simulations {
		serializer := SimulationSerializerNoAssoc{self.Ctx, simulation}
		response = append(response, serializer.Response())
	}
	return response
}

type SimulationSerializerNoAssoc struct {
	Ctx *gin.Context
	Simulation
}

type SimulationResponseNoAssoc struct {
	Name    string `json:"Name"`
	ID      uint   `json:"SimulationID"`
	Running bool   `json:"Running"`
	//StartParams postgres.Jsonb `json:"Starting Parameters"`
}

func (self *SimulationSerializerNoAssoc) Response() SimulationResponseNoAssoc {
	response := SimulationResponseNoAssoc{
		Name:    self.Name,
		ID:      self.ID,
		Running: self.Running,
		//StartParams: self.StartParameters,
	}
	return response
}

// File/s Serializers

type FilesSerializerNoAssoc struct {
	Ctx   *gin.Context
	Files []File
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
	File
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

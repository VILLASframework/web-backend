package user

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/project"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

type UsersSerializer struct {
	Ctx   *gin.Context
	Users []common.User
}

func (self *UsersSerializer) Response(assoc bool) []UserResponse {
	response := []UserResponse{}
	for _, user := range self.Users {
		serializer := UserSerializer{self.Ctx, user}
		response = append(response, serializer.Response(assoc))
	}
	return response
}

type UserSerializer struct {
	Ctx *gin.Context
	common.User
}

type UserResponse struct {
	Username    string `json:"Username"`
	Password    string `json:"Password"` // XXX: this is the hashed pw
	Role        string `json:"Role"`
	Mail        string `json:"Mail"`
	Projects    []project.ProjectResponseNoAssoc
	Simulations []simulation.SimulationResponseNoAssoc
	Files       []file.FileResponseNoAssoc
}

func (self *UserSerializer) Response(assoc bool) UserResponse {

	response := UserResponse{
		Username: self.Username,
		Password: self.Password,
		Role:     self.Role,
		Mail:     self.Mail,
	}

	// Associated models MUST NOT called with assoc=true otherwise we
	// will have an infinite loop due to the circular dependencies
	if assoc {

		// TODO: maybe all those should be made in one transaction
		projects, _, _ := project.FindUserProjects(&self.User)
		projectsSerializer :=
			project.ProjectsSerializerNoAssoc{self.Ctx, projects}

		simulations, _, _ := simulation.FindUserSimulations(&self.User)
		simulationsSerializer :=
			simulation.SimulationsSerializerNoAssoc{self.Ctx, simulations}

		files, _, _ := file.FindUserFiles(&self.User)
		filesSerializer := file.FilesSerializerNoAssoc{self.Ctx, files}

		// Add the associated models to the response
		response.Projects = projectsSerializer.Response()
		response.Simulations = simulationsSerializer.Response()
		response.Files = filesSerializer.Response()
	}

	return response
}

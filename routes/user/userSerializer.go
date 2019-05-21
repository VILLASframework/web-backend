package user

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

type UsersSerializer struct {
	Ctx   *gin.Context
	Users []common.User
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
	common.User
}

type UserResponse struct {
	Username    string `json:"Username"`
	Password    string `json:"Password"` // XXX: ???
	Role        string `json:"Role"`
	Mail        string `json:"Mail"`
	Simulations []simulation.SimulationResponseNoAssoc
}

func (self *UserSerializer) Response() UserResponse {
	// TODO: maybe all those should be made in one transaction

	simulations, _, _ := simulation.FindUserSimulations(&self.User)
	simulationsSerializer := simulation.SimulationsSerializerNoAssoc{self.Ctx, simulations}


	response := UserResponse{
		Username:    self.Username,
		Password:    self.Password,
		Role:        self.Role,
		Mail:        self.Mail,
		Simulations: simulationsSerializer.Response(),
	}
	return response
}





package serializers

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
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
	Simulations []SimulationResponseNoAssoc
}

func (self *UserSerializer) Response() UserResponse {
	// TODO: maybe all those should be made in one transaction

	simulations, _, _ := queries.FindUserSimulations(&self.User)
	simulationsSerializer := SimulationsSerializerNoAssoc{self.Ctx, simulations}


	response := UserResponse{
		Username:    self.Username,
		Password:    self.Password,
		Role:        self.Role,
		Mail:        self.Mail,
		Simulations: simulationsSerializer.Response(),
	}
	return response
}





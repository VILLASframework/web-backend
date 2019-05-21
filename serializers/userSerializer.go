package serializers

import (
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
	Role        string `json:"Role"`
	Mail        string `json:"Mail"`
}

func (self *UserSerializer) Response() UserResponse {
	response := UserResponse{
		Username:    self.Username,
		Role:        self.Role,
		Mail:        self.Mail,
	}
	return response
}





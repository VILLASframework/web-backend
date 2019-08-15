package user

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

var router *gin.Engine
var db *gorm.DB

func cleanUserTable() {
	db.DropTable(&common.User{})
	db.AutoMigrate(&common.User{})
	db.Create(&common.User0)
}

func TestMain(m *testing.M) {

	db = common.DummyInitDB()
	defer db.Close()

	common.DummyAddOnlyUserTableWithAdminDB(db)

	router = gin.Default()
	api := router.Group("/api")

	RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication(true))
	RegisterUserEndpoints(api.Group("/users"))

	os.Exit(m.Run())
}

func TestGetAllUsers(t *testing.T) {

	defer cleanUserTable()

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials, 200)
	assert.NoError(t, err)

	// test GET user/1 (the admin)
	err = common.NewTestEndpoint(router, token,
		"/api/users/1", "GET", nil,
		200, common.KeyModels{"user": common.User0})
	assert.NoError(t, err)

	// test POST user/ $newUser
	newUser := common.Request{
		Username: common.UserA.Username,
		Password: common.StrPasswordA,
		Mail:     common.UserA.Mail,
		Role:     common.UserA.Role,
	}
	// Get the number of alreday existing users so to know the expected
	// id of the new user
	maxid, err := common.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)
	err = common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser},
		200, common.KeyModels{"id": maxid + 1})
	assert.NoError(t, err)

	// test PUT user/1 $modifiedUser
	modifiedUser := common.Request{Role: "Admin"}
	err = common.NewTestEndpoint(router, token,
		"/api/users/2", "PUT", common.KeyModels{"user": modifiedUser},
		200, common.KeyModels{"id": 2})
	assert.NoError(t, err)
}

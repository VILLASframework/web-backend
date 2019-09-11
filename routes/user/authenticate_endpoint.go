package user

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type tokenClaims struct {
	UserID uint   `json:"id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func RegisterAuthenticate(r *gin.RouterGroup) {
	r.POST("", authenticate)
}

// authenticate godoc
// @Summary Authentication for user
// @ID authenticate
// @Accept json
// @Produce json
// @Tags authentication
// @Param inputUser body user.loginRequest true "loginRequest of user"
// @Success 200 {object} docs.ResponseAuthenticate "JSON web token, success status, message and authenticated user object"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 401 {object} docs.ResponseError "Unauthorized"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity."
// @Failure 500 {object} docs.ResponseError "Internal server error."
// @Router /authenticate [post]
func authenticate(c *gin.Context) {

	// Bind the response (context) with the loginRequest struct
	var credentials loginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Validate the login request
	if errs := credentials.validate(); errs != nil {
		helper.BadRequestError(c, errs.Error())
		return
	}

	// Find the username in the database
	var user User
	err := user.ByUsername(credentials.Username)
	if err != nil {
		helper.NotFoundError(c, "User not found")
		return
	}

	// Validate the password
	err = user.validatePassword(credentials.Password)
	if err != nil {
		helper.UnauthorizedError(c, "Invalid password")
		return
	}

	// create authentication token
	claims := tokenClaims{
		user.ID,
		user.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(weekHours).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://web.villas.fein-aachen.org/",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSigningSecret))
	if err != nil {
		helper.InternalServerError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authenticated",
		"token":   tokenString,
		"user":    user.User,
	})
}

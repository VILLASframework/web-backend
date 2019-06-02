package user

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserToContext(ctx *gin.Context, user_id uint) {
	var user common.User
	if user_id != 0 {
		db := common.GetDB()
		db.First(&user, user_id)
	}
	ctx.Set("user_id", user_id)
	ctx.Set("user", user)
}

func Authentication(unauthorized bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Initialize user_id and model in the context
		UserToContext(ctx, 0)

		// Authentication's access token extraction
		// XXX: if we have a multi-header for Authorization (e.g. in
		// case of OAuth2 use the request.OAuth2Extractor and make sure
		// that the argument is 'access-token' or provide a custom one
		token, err := request.ParseFromRequest(ctx.Request,
			request.AuthorizationHeaderExtractor,
			func(token *jwt.Token) (interface{}, error) {

				// validate alg for signing the jwt
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing alg: %v",
						token.Header["alg"])
				}

				// return secret in byte format
				secret := ([]byte(jwtSigningSecret))
				return secret, nil
			})

		// If the authentication extraction fails return HTTP CODE 401
		if err != nil {
			if unauthorized {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"succes":  false,
					"message": "Authentication failed",
				})
			}
			return
		}

		// If the token is ok, pass user_id to context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user_id, _ := strconv.ParseInt(claims["id"].(string), 10, 64)
			UserToContext(ctx, uint(user_id))
		}
	}
}

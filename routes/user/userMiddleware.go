package user

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
)

func userToContext(ctx *gin.Context, user_id uint) {

	var user User

	err := user.ByID(user_id)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"succes":  false,
			"message": "Authentication failed (user not found)",
		})
		return
	}

	ctx.Set(common.UserRoleCtx, user.Role)
	ctx.Set(common.UserIDCtx, user_id)
}

func Authentication(unauthorized bool) gin.HandlerFunc {

	return func(ctx *gin.Context) {

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
					"message": "Authentication failed (claims extraction)",
				})
			}
			return
		}

		// If the token is ok, pass user_id to context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			user_id, ok := claims["id"].(float64)

			if !ok {
				ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"succes":  false,
					"message": "Authentication failed (claims casting)",
				})
				return
			}

			userToContext(ctx, uint(user_id))
		}

	}
}

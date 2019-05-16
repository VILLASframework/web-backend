package user

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const signatureSecret = "_A_strong_password_as_enviromental_variable_"

func UserToContext(ctx *gin.Context, user_id uint) {
	var user common.User
	if user_id != 0 {
		db := common.GetDB()
		db.First(&user, user_id)
	}
	ctx.Set("user_id", user_id)
	ctx.Set("user", user)
}

// func stripBearerPrefixFromTokenString()
// Originally is supposed to remove the 'BEARER' token from the Auth
// header "Authorization: Bearer <token>". Currently use curl's header
// like "$ curl -H 'Authorization: TOKEN <token> ..."
func removeTokenPrefix(tok string) (string, error) {
	// if the prefix exists remove it from token
	if len(tok) > 5 && strings.ToUpper(tok[0:6]) == "TOKEN " {
		return tok[6:], nil
	}
	// otherwise return token
	return tok, nil
}

// Extractor of Authorization Header
// var AuthorizationHeaderExtractor
var GetAuthorizationHeader = &request.PostExtractionFilter{
	request.HeaderExtractor{"Authorization"},
	removeTokenPrefix,
}

// Extractor of OAuth2 tokens. Finds the 'access_token'
// var OAuth2Extractor
var GetAuth2 = &request.MultiExtractor{
	GetAuthorizationHeader,
	request.ArgumentExtractor{"access_token"},
}

func Authentication(unauthorized bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Initialize user_id and model in the context
		UserToContext(ctx, 0)

		// Authentication's access token extraction
		token, err := request.ParseFromRequest(ctx.Request, GetAuth2,
			func(token *jwt.Token) (interface{}, error) {
				// validate alg for signing the jwt
				// XXX: whis is the default signing method?
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing alg: %v",
						token.Header["alg"])
				}
				// return secret in byte format
				secret := ([]byte(signatureSecret))
				return secret, nil
			})

		// If the authentication extraction fails return HTTP CODE 401
		if err != nil {
			if unauthorized {
				ctx.AbortWithError(http.StatusUnauthorized, err)
			}
			return
		}

		// If the token is ok, pass user_id to context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			user_id := uint(claims["id"].(float64))
			UserToContext(ctx, user_id)
		}
	}
}

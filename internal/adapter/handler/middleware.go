package handler

import (
	"github.com/OzkrOssa/freeradius-api/internal/core/domain"
	"github.com/OzkrOssa/freeradius-api/internal/core/port"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationType       = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(token port.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)

		isEmpty := len(authorizationHeader) == 0
		if isEmpty {
			err := domain.EmptyAuthorizationHeaderError
			errorMessages := parseError(err)
			errResponse := newErrorResponse(errorMessages)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		fields := strings.Fields(authorizationHeader)
		isValid := len(fields) == 2
		if !isValid {
			err := domain.InvalidAuthorizationHeaderError
			errorMessages := parseError(err)
			errResponse := newErrorResponse(errorMessages)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		currentAuthorizationType := strings.ToLower(fields[0])
		if currentAuthorizationType != authorizationType {
			err := domain.InvalidAuthorizationTypeError
			errorMessages := parseError(err)
			errResponse := newErrorResponse(errorMessages)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		accessToken := fields[1]
		payload, err := token.VerifyToken(accessToken)
		if err != nil {
			errorMessages := parseError(err)
			errResponse := newErrorResponse(errorMessages)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errResponse)
			return
		}

		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}

//func adminMiddleware() gin.HandlerFunc {
//}

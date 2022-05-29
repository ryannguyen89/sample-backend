package api

import (
	"fmt"
	"net/http"

	"github.com/auth0/go-jwt-middleware/v2"
	"github.com/gin-gonic/gin"
)

func (api *API) authorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token
		token, err := jwtmiddleware.AuthHeaderTokenExtractor(c.Request)
		if err != nil || token == "" {
			c.Status(http.StatusUnauthorized)
			c.Abort()
		}
		fmt.Println("token:", token)
		err = api.userSvc.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.Status(http.StatusUnauthorized)
			c.Abort()
		}
	}
}

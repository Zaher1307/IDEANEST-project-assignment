package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/auth"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/types"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusBadRequest, types.MessageResp{
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusBadRequest, types.MessageResp{
				Message: "Invalid Authorization header",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]

		email, err := auth.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusBadRequest, types.MessageResp{
				Message: "Invalid Authorization header",
			})
			c.Abort()
			return

		}

		c.Set("email", email)

		c.Next()
	}
}

package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/signup", SignUpHandler)
	r.POST("/signin", SignInHandler)
	r.POST("/refresh-token", RefreshTokenHandler)

	r.Use(AuthMiddleware())

	r.POST("/organization", CreateOrgHandler)
	r.GET("/organization/:organization_id", ReadOrgHandler)
	r.GET("/organization", ReadAllOrgsHandler)
	r.PUT("/organization/:organization_id", UpdateOrgHandler)
	r.DELETE("/organization/:organization_id", DeleteOrgHandler)
	r.POST("/organization/:organization_id/invite", InviteUserToOrgHandler)
	r.POST("/revoke-refresh-token", RevokeRefreshTokenHandler)

	r.Run(":8080")
}

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/business"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/types"
)

func SignUpHandler(c *gin.Context) {
	signUpReq := types.SignUpReq{}
	if err := c.ShouldBindJSON(&signUpReq); err != nil {
		c.JSON(http.StatusBadRequest, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	user := types.User{
		UserInfo: types.UserInfo{
			Name:  signUpReq.Name,
			Email: signUpReq.Email,
		},
		Password: signUpReq.Password,
	}

	err := business.SignUp(user)
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.MessageResp{
		Message: "Succeeded",
	})
}

func SignInHandler(c *gin.Context) {
	signInReq := types.SignInReq{}
	if err := c.ShouldBindJSON(&signInReq); err != nil {
		c.JSON(http.StatusBadRequest, types.TokenResp{
			Message:      "Faild: " + err.Error(),
			AccessToken:  "",
			RefreshToken: "",
		})
		return
	}

	user := types.User{
		UserInfo: types.UserInfo{
			Email: signInReq.Email,
		},
		Password: signInReq.Password,
	}

	tokens, err := business.SignIn(user)
	if err != nil {
		c.JSON(http.StatusOK, types.TokenResp{
			Message:      "Faild: " + err.Error(),
			AccessToken:  "",
			RefreshToken: "",
		})
		return
	}

	c.JSON(http.StatusOK, types.TokenResp{
		Message:      "Succeeded",
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func RefreshTokenHandler(c *gin.Context) {
	refreshTokenReq := types.RefreshTokenReq{}
	if err := c.ShouldBindJSON(&refreshTokenReq); err != nil {
		c.JSON(http.StatusBadRequest, types.TokenResp{
			Message:      "Faild: " + err.Error(),
			AccessToken:  "",
			RefreshToken: "",
		})
		return
	}

	accessToken, err := business.RefreshAccessToken(refreshTokenReq.RefreshToken)
	if err != nil {
		c.JSON(http.StatusOK, types.TokenResp{
			Message:      "Faild: " + err.Error(),
			AccessToken:  "",
			RefreshToken: "",
		})
		return
	}

	c.JSON(http.StatusOK, types.TokenResp{
		Message:      "Succeeded",
		AccessToken:  accessToken,
		RefreshToken: refreshTokenReq.RefreshToken,
	})
}

func CreateOrgHandler(c *gin.Context) {
	createOrgReq := types.CreateOrgReq{}
	if err := c.ShouldBindJSON(&createOrgReq); err != nil {
		c.JSON(http.StatusBadRequest, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	orgInfo := types.OrgInfo{
		Name:        createOrgReq.Name,
		Description: createOrgReq.Description,
	}

	email, _ := c.Get("email")

	id, err := business.CreateOrg(orgInfo, email.(string))
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.CreateOrgResp{
		OrgId: id,
	})
}

func ReadOrgHandler(c *gin.Context) {
	orgId := c.Param("organization_id")
	email, _ := c.Get("email")

	org, err := business.ReadOrg(orgId, email.(string))
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	readOrgResp := types.ReadOrgResp{
		OrgId:       org.OrgId,
		Name:        org.Name,
		Description: org.Description,
	}

	for _, orgMember := range org.OrgMembers {
		orgMemberResp := types.OrgMemberResp{
			Name:        orgMember.Name,
			Email:       orgMember.Email,
			AccessLevel: orgMember.AccessLevel,
		}
		readOrgResp.OrgMembers = append(readOrgResp.OrgMembers, orgMemberResp)
	}

	c.JSON(http.StatusOK, readOrgResp)
}

func ReadAllOrgsHandler(c *gin.Context) {
	email, _ := c.Get("email")

	orgs, err := business.ReadAllOrgs(email.(string))
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	var allOrgsResp []types.ReadOrgResp
	for _, org := range orgs {
		readOrgResp := types.ReadOrgResp{
			OrgId:       org.OrgId,
			Name:        org.Name,
			Description: org.Description,
		}

		for _, orgMember := range org.OrgMembers {
			orgMemberResp := types.OrgMemberResp{
				Name:        orgMember.Name,
				Email:       orgMember.Email,
				AccessLevel: orgMember.AccessLevel,
			}
			readOrgResp.OrgMembers = append(readOrgResp.OrgMembers, orgMemberResp)
		}

		allOrgsResp = append(allOrgsResp, readOrgResp)
	}

	c.JSON(http.StatusOK, allOrgsResp)
}

func UpdateOrgHandler(c *gin.Context) {
	updateOrgReq := types.UpdateOrgReq{}
	if err := c.ShouldBindJSON(&updateOrgReq); err != nil {
		c.JSON(http.StatusBadRequest, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	email, _ := c.Get("email")
	orgId := c.Param("organization_id")
	orgInfo := types.OrgInfo{
		OrgId:       orgId,
		Name:        updateOrgReq.Name,
		Description: updateOrgReq.Description,
	}

	orgInfoMod, err := business.UpdateOrg(orgInfo, email.(string))
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}
	updateOrgResp := types.UpdateOrgResp{
		OrgId:       orgId,
		Name:        orgInfoMod.Name,
		Description: orgInfo.Description,
	}

	c.JSON(http.StatusOK, updateOrgResp)
}

func DeleteOrgHandler(c *gin.Context) {
	email, _ := c.Get("email")
	orgId := c.Param("organization_id")

	err := business.DeleteOrg(orgId, email.(string))
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.MessageResp{
		Message: "Succeeded",
	})
}

func InviteUserToOrgHandler(c *gin.Context) {
	inviteReq := types.InviteReq{}
	if err := c.ShouldBindJSON(&inviteReq); err != nil {
		c.JSON(http.StatusBadRequest, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	email, _ := c.Get("email")
	orgId := c.Param("organization_id")
	member := types.OrgMember{
		UserInfo: types.UserInfo{
			Email: inviteReq.Email,
		},
		AccessLevel: types.ACCESS_LEVEL_USER,
	}

	err := business.InviteUserToOrg(orgId, email.(string), member)
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.MessageResp{
		Message: "Succeeded",
	})
}

func RevokeRefreshTokenHandler(c *gin.Context) {
	refreshTokenReq := types.RefreshTokenReq{}
	if err := c.ShouldBindJSON(&refreshTokenReq); err != nil {
		c.JSON(http.StatusBadRequest, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	email, _ := c.Get("email")

	err := business.RevokeRefreshToken(refreshTokenReq.RefreshToken, email.(string))
	if err != nil {
		c.JSON(http.StatusOK, types.MessageResp{
			Message: "Faild: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, types.MessageResp{
		Message: "Succeeded",
	})
}

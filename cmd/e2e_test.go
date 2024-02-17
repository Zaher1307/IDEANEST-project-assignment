package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/business"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/types"
)

var (
	r   *gin.Engine
	c   *resty.Client
	url string
)

func init() {
	c = resty.New()
	r = gin.Default()
	url = httptest.NewServer(r).URL

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
}

func TestHandlers(t *testing.T) {
	t.Run("SignUpHanlder", func(t *testing.T) {
		// Bad request with wrong body with status 400
		resp, _ := c.R().
			SetBody(`{"username":"testuser", "password":"testpass"}`).
			Post(url + "/signup")

		if resp.StatusCode() != http.StatusBadRequest {
			t.Errorf("Expected status %d but got %d", http.StatusOK, resp.StatusCode())
		}

		// good request with status 200
		resp, _ = c.R().
			SetBody(`{"name":"testname", "email": "test@mail.mail", "password":"testpass"}`).
			Post(url + "/signup")

		if resp.StatusCode() != http.StatusOK {
			t.Errorf("Expected status %d but got %d", http.StatusOK, resp.StatusCode())
		}
	})

	t.Run("SignInHandler", func(t *testing.T) {
		// Bad request with wrong body with faild message
		resp, _ := c.R().
			SetBody(`{"email":"test@mail.com", "password":"testpass"}`).
			Post(url + "/signin")

		faildMessage := `{"message":"Faild: user doesn't exists","access_token":"","refresh_token":""}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}

		// good request with status 200
		business.SignUp(types.User{
			UserInfo: types.UserInfo{
				Name:  "ahmed",
				Email: "zaher@a.b",
			},
			Password: "123",
		})

		resp, _ = c.R().
			SetBody(`{"email":"zaher@a.b", "password":"123"}`).
			Post(url + "/signin")

		if string(resp.Body()) == faildMessage {
			t.Errorf("Expected request to success")
		}
	})

	t.Run("RefreshTokenHandler", func(t *testing.T) {
		business.SignUp(types.User{
			UserInfo: types.UserInfo{
				Name:  "ahmed",
				Email: "zaher@a.b",
			},
			Password: "123",
		})
		tokens, _ := business.SignIn(types.User{
			UserInfo: types.UserInfo{
				Email: "zaher@a.b",
			},
			Password: "123",
		})

		refreshToken := tokens.RefreshToken

		resp, _ := c.R().
			SetBody(`{"refresh_token":"` + refreshToken + `", "password":"123"}`).
			Post(url + "/refresh-token")

		var body map[string]string
		json.Unmarshal(resp.Body(), &body)

		if body["message"] != "Succeeded" {
			t.Errorf("Expected request to success")
		}
	})

	// all the following will not have full functionality testing only authentication

	t.Run("CreateOrgHandler", func(t *testing.T) {
		resp, _ := c.R().
			SetBody(`{"name":"org name", "description":"org description"}`).
			Post(url + "/organization")

		faildMessage := `{"message":"Unauthorized"}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}
	})

	t.Run("ReadOrgHandler", func(t *testing.T) {
		resp, _ := c.R().
			Get(url + "/organization/21324124")

		faildMessage := `{"message":"Unauthorized"}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}
	})

	t.Run("ReadAllOrgsHandler", func(t *testing.T) {
		resp, _ := c.R().
			Get(url + "/organization")

		faildMessage := `{"message":"Unauthorized"}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}
	})

	t.Run("UpdateOrgHandler", func(t *testing.T) {
		resp, _ := c.R().
			SetBody(`{"name":"org name", "description":"org description"}`).
			Put(url + "/organization")

		faildMessage := `{"message":"Unauthorized"}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}
	})

	t.Run("DeleteOrgHandler", func(t *testing.T) {
		resp, _ := c.R().
			Delete(url + "/organization")

		faildMessage := `{"message":"Unauthorized"}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}
	})

	t.Run("InviteUserToOrgHandler", func(t *testing.T) {
		resp, _ := c.R().
			SetBody(`{"name":"org name", "description":"org description"}`).
			Post(url + "/organization/:1234/invite")

		faildMessage := `{"message":"Unauthorized"}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}
	})

	t.Run("RevokeRefreshTokenHandler", func(t *testing.T) {
		resp, _ := c.R().
			SetBody(`{"refresh_token":"refresh_token"`).
			Post(url + "/revoke-refresh-token")

		faildMessage := `{"message":"Unauthorized"}`

		if string(resp.Body()) != faildMessage {
			t.Errorf("Expected faild message %s but got %s", faildMessage, string(resp.Body()))
		}
	})
}

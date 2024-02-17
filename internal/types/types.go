package types

const (
	ACCESS_LEVEL_ADMIN = "admin"
	ACCESS_LEVEL_USER  = "user"

	USER_COLL = "user"
	ORG_COLL  = "organization"
)

type UserInfo struct {
	Name  string `bson:"name"`
	Email string `bson:"email"`
}

type User struct {
	UserInfo `bson:",inline"`
	Password string   `bson:"password"`
	Orgs     []string `bson:"organizations"`
}

type OrgMember struct {
	UserInfo    `bson:",inline"`
	AccessLevel string `bson:"access_level"`
}

type OrgInfo struct {
	OrgId       string `bson:"-"`
	Name        string `bson:"name"`
	Description string `bson:"description"`
}

type Org struct {
	OrgInfo    `bson:",inline"`
	OrgMembers []OrgMember `bson:"organization_members"`
}

type Token struct {
	RefreshToken string
	AccessToken  string
}

// ===================== Consumer Request Structures ===================== //

type SignUpReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignInReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type CreateOrgReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateOrgReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type InviteReq struct {
	Email string `json:"user_email" binding:"required"`
}

// ===================== Consumer Response Structures ===================== //

type MessageResp struct {
	Message string `json:"message"`
}

type TokenResp struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type CreateOrgResp struct {
	OrgId string `json:"organization_id"`
}

type UpdateOrgResp struct {
	OrgId       string `json:"organization_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type OrgMemberResp struct {
	Name        string `json:"name"`
	Email       string `json:"user_email"`
	AccessLevel string `json:"access_level"`
}

type ReadOrgResp struct {
	OrgId       string          `json:"organization_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	OrgMembers  []OrgMemberResp `json:"organization_members"`
}

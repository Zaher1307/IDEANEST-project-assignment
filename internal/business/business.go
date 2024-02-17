package business

import (
	"errors"
	"log"

	"github.com/zaher1307/IDEANEST-project-assignment/internal/auth"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/database"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/types"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal(err)
	}
}

func SignUp(user types.User) error {
	var err error

	existedUser, err := database.ReadUser(user.Email)
	if err != nil {
		return err
	}

	if existedUser.Email == user.Email {
		return errors.New("email already exists")
	}

	user.Password, err = hashPassword(user.Password)
	if err != nil {
		return err
	}

	return database.CreateUser(user)
}

func SignIn(user types.User) (types.Token, error) {
	fetchedUser, err := database.ReadUser(user.Email)
	if err != nil {
		return types.Token{}, err
	}

	if fetchedUser.Email == "" {
		return types.Token{}, errors.New("user doesn't exists")
	}

	err = verifyPassword(user.Password, fetchedUser.Password)
	if err != nil {
		return types.Token{}, err
	}

	refreshToken, err := auth.GenerateRefreshToken(fetchedUser)
	if err != nil {
		return types.Token{}, err
	}

	accessToken, err := auth.GenerateAccessToken(refreshToken)
	if err != nil {
		return types.Token{}, err
	}

	return types.Token{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func RevokeRefreshToken(refreshToken, email string) error {
	userEmail, err := auth.GetRefreshTokenUserEmail(refreshToken)
	if err != nil {
		return err
	}

	if userEmail != email {
		return errors.New("cannot revoke unknown token")
	}

	return auth.RevokeRefreshToken(refreshToken)
}

func RefreshAccessToken(refreshToken string) (string, error) {
	return auth.GenerateAccessToken(refreshToken)
}

func CreateOrg(orgInfo types.OrgInfo, email string) (string, error) {
	user, err := database.ReadUser(email)
	if err != nil {
		return "", nil
	}
	return database.CreateOrg(orgInfo, user)
}

func ReadOrg(orgId, email string) (types.Org, error) {
	isMember := database.IsOrgMember(orgId, email)
	if !isMember {
		return types.Org{}, errors.New("this user is not an org member")
	}

	org, err := database.ReadOrg(orgId)
	if err != nil {
		return types.Org{}, err
	}

	return org, nil
}

func ReadAllOrgs(email string) ([]types.Org, error) {
	return database.ReadAllOrgsInfo(email)
}

func UpdateOrg(orgInfo types.OrgInfo, email string) (types.OrgInfo, error) {
	adminEmail, err := database.ReadOrgAdmin(orgInfo.OrgId)
	if err != nil {
		return types.OrgInfo{}, err
	}

	if adminEmail != email {
		return types.OrgInfo{}, errors.New("orgs can be updated via admins only")
	}

	return database.UpdateOrg(orgInfo)
}

func DeleteOrg(orgId, email string) error {
	adminEmail, err := database.ReadOrgAdmin(orgId)
	if err != nil {
		return err
	}

	if adminEmail != email {
		return errors.New("orgs can be deleted via admins only")
	}

	return database.DeleteOrg(orgId)
}

func InviteUserToOrg(orgId, email string, member types.OrgMember) error {
	adminEmail, err := database.ReadOrgAdmin(orgId)
	if err != nil {
		return err
	}

	if adminEmail != email {
		return errors.New("inviting users to orgs can done only be admins")
	}

	user, err := database.ReadUser(member.Email)
	if err != nil {
		return err
	}

	if user.Email != member.Email {
		return errors.New("user doesn't exists")
	}

	for _, org := range user.Orgs {
		if org == orgId {
			return errors.New("user already exists in this organization")
		}
	}

	member.Name = user.Name

	return database.InviteUserToOrg(orgId, member)
}

// ================ Private helper functions ================ //

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func verifyPassword(inputPassword, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
	return err
}

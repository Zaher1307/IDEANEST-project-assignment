package database

import (
	"context"
	"os"

	"github.com/joho/godotenv"
	"github.com/zaher1307/IDEANEST-project-assignment/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client    *mongo.Client
	ctx       context.Context
	mongoUser string
	mongoPass string
	mongoDB   string
	mongoHost string
)

func init() {
	godotenv.Load()

	mongoUser = os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoPass = os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	mongoDB = os.Getenv("DATABASE_NAME")
	mongoHost = os.Getenv("DATABASE_HOST")
}

func ConnectDB() error {
	uri := "mongodb://" + mongoUser + ":" + mongoPass + "@" + mongoHost + ":27017/"

	ctx = context.Background()

	clientOptions := options.Client().ApplyURI(uri)

	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}

func DisconnectDB() error {
	if err := client.Disconnect(ctx); err != nil {
		return err
	}

	return nil
}

func CreateUser(user types.User) error {
	collection := client.Database(mongoDB).Collection(types.USER_COLL)
	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func ReadUser(email string) (types.User, error) {
	collection := client.Database(mongoDB).Collection(types.USER_COLL)
	filter := bson.M{"email": email}

	var user types.User
	collection.FindOne(ctx, filter).Decode(&user)

	return user, nil
}

func CreateOrg(orgInfo types.OrgInfo, user types.User) (string, error) {
	collection := client.Database(mongoDB).Collection(types.ORG_COLL)

	result, err := collection.InsertOne(ctx, types.Org{
		OrgInfo: orgInfo,
	})
	if err != nil {
		return "", err
	}

	id := result.InsertedID.(primitive.ObjectID).Hex()

	member := types.OrgMember{
		UserInfo: types.UserInfo{
			Name:  user.Name,
			Email: user.Email,
		},
		AccessLevel: types.ACCESS_LEVEL_ADMIN,
	}

	InviteUserToOrg(id, member)

	return id, nil
}

func UpdateOrg(orgInfo types.OrgInfo) (types.OrgInfo, error) {
	collection := client.Database(mongoDB).Collection(types.ORG_COLL)
	id, err := primitive.ObjectIDFromHex(orgInfo.OrgId)
	if err != nil {
		return types.OrgInfo{}, err
	}
	filter := bson.M{"_id": id}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "name", Value: orgInfo.Name},
			{Key: "description", Value: orgInfo.Description},
		}},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return types.OrgInfo{}, err
	}

	return types.OrgInfo{
		OrgId:       orgInfo.OrgId,
		Name:        orgInfo.Name,
		Description: orgInfo.Description,
	}, nil
}

func ReadOrg(orgId string) (types.Org, error) {
	collection := client.Database(mongoDB).Collection(types.ORG_COLL)
	id, err := primitive.ObjectIDFromHex(orgId)
	if err != nil {
		return types.Org{}, err
	}

	filter := bson.M{"_id": id}

	var org types.Org
	err = collection.FindOne(ctx, filter).Decode(&org)
	if err != nil {
		return types.Org{}, err
	}

	org.OrgId = orgId

	return org, nil
}

func ReadOrgAdmin(orgId string) (string, error) {
	org, err := ReadOrg(orgId)
	if err != nil {
		return "", err
	}

	for _, member := range org.OrgMembers {
		if member.AccessLevel == types.ACCESS_LEVEL_ADMIN {
			return member.Email, nil
		}
	}

	return "", nil
}

func ReadAllOrgsInfo(email string) ([]types.Org, error) {
	user, err := ReadUser(email)
	if err != nil {
		return nil, err
	}

	orgsId := make([]primitive.ObjectID, len(user.Orgs))

	for i, orgId := range user.Orgs {
		orgsId[i], err = primitive.ObjectIDFromHex(orgId)
		if err != nil {
			return nil, err
		}
	}

	collection := client.Database(mongoDB).Collection(types.ORG_COLL)
	filter := bson.M{"_id": bson.M{"$in": orgsId}}

	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orgs []types.Org
	for cursor.Next(ctx) {
		var mapForExtractId map[string]interface{}
		err := cursor.Decode(&mapForExtractId)
		if err != nil {
			return nil, err
		}

		var org types.Org
		err = cursor.Decode(&org)
		if err != nil {
			return nil, err
		}

		org.OrgId = mapForExtractId["_id"].(primitive.ObjectID).Hex()

		orgs = append(orgs, org)
	}

	return orgs, nil
}

func DeleteOrg(orgId string) error {
	collection := client.Database(mongoDB).Collection(types.ORG_COLL)
	id, err := primitive.ObjectIDFromHex(orgId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": id}

	_, err = collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	err = removeOrgFromUsers(orgId)
	if err != nil {
		return err
	}

	return nil
}

func InviteUserToOrg(orgId string, member types.OrgMember) error {
	org, err := ReadOrg(orgId)
	if err != nil {
		return err
	}
	org.OrgMembers = append(org.OrgMembers, member)

	err = updateOrgMembers(org.OrgMembers, orgId)
	if err != nil {
		return err
	}

	addOrgToUser(member, orgId)

	return nil
}

func IsOrgMember(orgId, email string) bool {
	collection := client.Database(mongoDB).Collection(types.USER_COLL)
	filter := bson.M{"email": email}

	var fetchedUser types.User
	collection.FindOne(ctx, filter).Decode(&fetchedUser)

	for _, memberOrgId := range fetchedUser.Orgs {
		if orgId == memberOrgId {
			return true
		}
	}

	return false
}

// ====================== helper private function ====================== //

func updateOrgMembers(members []types.OrgMember, orgId string) error {
	collection := client.Database(mongoDB).Collection(types.ORG_COLL)
	id, err := primitive.ObjectIDFromHex(orgId)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": id}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "organization_members", Value: members},
		}},
	}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func addOrgToUser(member types.OrgMember, orgId string) error {
	collection := client.Database(mongoDB).Collection(types.USER_COLL)
	filter := bson.M{"email": member.Email}

	var fetchedUser types.User
	collection.FindOne(ctx, filter).Decode(&fetchedUser)

	fetchedUser.Orgs = append(fetchedUser.Orgs, orgId)

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "organizations", Value: fetchedUser.Orgs},
		}},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil
	}

	return nil
}

func removeOrgFromUsers(orgId string) error {
	collection := client.Database(mongoDB).Collection(types.USER_COLL)

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return err
	}

	var users []types.User
	for cursor.Next(context.Background()) {
		var user types.User
		err := cursor.Decode(&user)
		if err != nil {
			return err
		}
		users = append(users, user)
	}

	var updateModels []mongo.WriteModel

	for i, user := range users {
		for j, org := range user.Orgs {
			if org == orgId {
				orgs := append(users[i].Orgs[:j], users[i].Orgs[j+1:]...)
				update := bson.D{
					{Key: "$set", Value: bson.D{
						{Key: "organizations", Value: orgs},
					}},
				}
				updateModels = append(updateModels, mongo.NewUpdateOneModel().SetFilter(bson.M{"email": user.Email}).SetUpdate(update))
			}
		}
	}

	_, err = collection.BulkWrite(context.Background(), updateModels)
	if err != nil {
		return err
	}

	return nil
}

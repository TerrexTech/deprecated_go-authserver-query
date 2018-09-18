package schema

import (
	"github.com/TerrexTech/go-authserver-query/auth"
	"github.com/TerrexTech/go-authserver-query/model"
	"github.com/gofrs/uuid"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

var loginResolver = func(params graphql.ResolveParams) (interface{}, error) {
	rootValue := params.Info.RootValue.(map[string]interface{})
	redis := rootValue["redis"].(auth.TokenStoreI)
	db := rootValue["authDB"].(auth.DBI)

	// uid, _ := uuid.NewV4()
	// u := &model.User{
	// 	UUID:      uid,
	// 	Username:  "test",
	// 	FirstName: "Deadly",
	// 	LastName:  "Potato",
	// 	Password:  "testing",
	// 	Role:      "manager",
	// 	Email:     "asd",
	// }
	// b, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	// u.Password = string(b)
	// _, e := db.Collection().InsertOne(u)
	// log.Println(e)

	username := params.Args["username"].(string)
	password := params.Args["password"].(string)
	user := &model.User{
		Username: username,
		Password: password,
	}

	return auth.Login(db, redis, user)
}

// unc(ts TokenStoreI,
// 	rt *model.RefreshToken,
// 	user *model.User) (*model.AccessToken, er

var accessTokenResolver = func(params graphql.ResolveParams) (interface{}, error) {
	rootValue := params.Info.RootValue.(map[string]interface{})
	redis := rootValue["redis"].(auth.TokenStoreI)
	db := rootValue["authDB"].(auth.DBI)

	rtStr := params.Args["refreshToken"].(string)
	uid := params.Args["sub"].(string)

	parsedUID, err := uuid.FromString(uid)
	if err != nil {
		return nil, errors.New("Error parsing RefreshToken: Cannot parse Sub")
	}

	rt := &model.RefreshToken{
		Sub:   parsedUID,
		Token: rtStr,
	}

	user, err := db.UserByUUID(parsedUID)
	if err != nil {
		return nil, errors.Wrap(
			err,
			"Error parsing RefreshToken: Cannot get user with specified UUID",
		)
	}
	at, err := auth.RefreshAccessToken(redis, rt, user)
	if err != nil {
		err = errors.Wrap(err, "Error renewing AccessToken")
		return nil, err
	}
	return &model.AuthResponse{
		AccessToken: at,
	}, nil
}

package gqltype

import "github.com/graphql-go/graphql"

var AuthResponse = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "AuthResponse",
		Fields: graphql.Fields{
			"access_token": &graphql.Field{
				Type: graphql.String,
			},
			"refresh_token": &graphql.Field{
				Type: graphql.String,
			},
			"user_data": &graphql.Field{
				Type: TokenData,
			},
		},
	},
)

var TokenData = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "TokenData",
		Fields: graphql.Fields{
			"first_name": &graphql.Field{
				Type: graphql.String,
			},
			"last_name": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

package schema

import (
	"github.com/TerrexTech/go-authserver-query/gqltype"
	"github.com/graphql-go/graphql"
)

var RootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"login": &graphql.Field{
			Type:        gqltype.AuthResponse,
			Description: "User Authentication",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"username": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: loginResolver,
		},
		"accessToken": &graphql.Field{
			Type:        gqltype.AuthResponse,
			Description: "User Authentication",
			Args: graphql.FieldConfigArgument{
				"refreshToken": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"sub": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: accessTokenResolver,
		},
	},
})

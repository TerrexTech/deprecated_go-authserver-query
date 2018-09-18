package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TerrexTech/go-authserver-query/auth"
	"github.com/TerrexTech/go-authserver-query/schema"
	"github.com/go-redis/redis"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

var Schema graphql.Schema

func graphqlHandler(w http.ResponseWriter, r *http.Request) {
	redis, err := auth.NewRedis(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	if err != nil {
		err = errors.Wrap(err, "Error creating Redis client")
		log.Println(err)
		return
	}

	db, err := auth.EnsureAuthDB()
	if err != nil {
		err = errors.Wrap(err, "Error connecting to Auth-DB")
		log.Println(err)
		return
	}

	rootObject := map[string]interface{}{
		"redis":  redis,
		"authDB": db,
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		err = errors.Wrap(err, "Error reading request body")
		log.Println(err)
		return
	}

	result := graphql.Do(graphql.Params{
		Schema:        Schema,
		RequestString: string(reqBytes),
		RootObject:    rootObject,
	})
	if len(result.Errors) > 0 {
		log.Printf("gql error: %v", result.Errors)
	}

	json.NewEncoder(w).Encode(result)
}

func init() {
	s, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: schema.RootQuery,
	})
	if err != nil {
		log.Fatalf("Error creating GraphQL Schema: %v", err)
	}
	Schema = s
}

func main() {
	http.HandleFunc("/graphql", graphqlHandler)
	http.ListenAndServe(":8080", nil)
	log.Println("listening on port 8080")
}

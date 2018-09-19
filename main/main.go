package main

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"
	"github.com/TerrexTech/go-authserver-query/auth"
	"github.com/TerrexTech/go-authserver-query/kafka"
	"github.com/TerrexTech/go-commonutils/commonutil"
	esmodel "github.com/TerrexTech/go-eventstore-models/model"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func main() {
	// Load environment-file.
	// Env vars will be read directly from environment if this file fails loading
	err := godotenv.Load()
	if err != nil {
		err = errors.Wrap(err,
			".env file not found, env-vars will be read as set in environment",
		)
		log.Println(err)
	}

	missingVar, err := commonutil.ValidateEnv(
		"KAFKA_BROKERS",
		"KAFKA_CONSUMER_GROUP_LOGIN",
		"KAFKA_CONSUMER_TOPIC_LOGIN",
		"KAFKA_PRODUCER_TOPIC_LOGIN",
	)
	if err != nil {
		log.Fatalf(
			"Error: Environment variable %s is required but was not found", missingVar,
		)
	}

	// Init IO
	db, err := auth.EnsureAuthDB()
	if err != nil {
		err = errors.Wrap(err, "Error connecting to Auth-DB")
		log.Println(err)
		return
	}
	kio, err := initKafkaIOLogin()
	if err != nil {
		err = errors.Wrap(err, "Error creating KafkaIO")
		log.Println(err)
		return
	}

	// Listen on Error-Channels
	go func() {
		for err := range kio.ConsumerErrors() {
			err = errors.Wrap(err, "Consumer Error")
			log.Println(err)
		}
	}()
	go func() {
		for err := range kio.ProducerErrors() {
			parsedErr := errors.Wrap(err.Err, "Producer Error")
			log.Println(parsedErr)
			log.Println(err)
		}
	}()

	for msg := range kio.ConsumerMessages() {
		go handleRequest(db, kio, msg)
	}
}

func handleRequest(db auth.DBI, kio *kafka.IO, msg *sarama.ConsumerMessage) {
	// Unmarshal msg to KafkaResponse
	kr := &esmodel.KafkaResponse{}
	err := json.Unmarshal(msg.Value, kr)
	if err != nil {
		err = errors.Wrap(err, "Error unmarshalling message into KafkaResponse")
		log.Println(err)
		kio.MarkOffset() <- msg
		return
	}

	// Unmarshal KafkaResponse.Input to User
	user := &auth.User{}
	err = json.Unmarshal([]byte(kr.Input), user)
	if err != nil {
		err = errors.Wrap(err, "Error unmarshalling KafkaResponse input into User")
		log.Println(err)
		kio.MarkOffset() <- msg
		return
	}

	// Authenticate User from DB
	errStr := ""
	var userJSON string
	user, err = db.Login(user)
	if err != nil {
		err = errors.Wrap(err, "Eror authenticating user")
		log.Println(err)
		kio.MarkOffset() <- msg
		errStr = err.Error()
	} else {
		// Marshal User
		mUser, err := user.MarshalJSON()
		if err != nil {
			err = errors.Wrap(err, "Error marshalling login-user into JSON")
			log.Println(err)
			kio.MarkOffset() <- msg
			errStr = err.Error()
		} else {
			userJSON = string(mUser)
		}
	}

	kr = &esmodel.KafkaResponse{
		CorrelationID: kr.CorrelationID,
		Error:         errStr,
		Result:        string(userJSON),
	}
	kio.ProducerInput() <- kr
	kio.MarkOffset() <- msg
}

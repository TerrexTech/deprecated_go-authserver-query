package auth

import (
	"github.com/TerrexTech/go-mongoutils/mongo"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// DBI is the Database-interface for authentication.
// This fetches/writes data to/from database for auth-actions such as
// login, registeration etc.
type DBI interface {
	Collection() *mongo.Collection
	UserByUUID(uid uuid.UUID) (*User, error)
	Login(user *User) (*User, error)
}

// DB is the implementation for dbI.
// dbI is the Database-interface for authentication.
// It fetches/writes data to/from database for auth-actions such as
// login, registeration etc.
type DB struct {
	collection *mongo.Collection
}

// EnsureAuthDB exists ensures that the required Database and Collection exists before
// auth-operations can be done on them. It creates Database/Collection if they don't exist.
func EnsureAuthDB() (*DB, error) {
	// Would ideally set these config-params as environment vars
	config := mongo.ClientConfig{
		Hosts:               []string{"localhost:27017"},
		Username:            "root",
		Password:            "root",
		TimeoutMilliseconds: 3000,
	}

	client, err := mongo.NewClient(config)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}

	conn := &mongo.ConnectionConfig{
		Client:  client,
		Timeout: 5000,
	}

	indexConfigs := []mongo.IndexConfig{
		mongo.IndexConfig{
			ColumnConfig: []mongo.IndexColumnConfig{
				mongo.IndexColumnConfig{
					Name: "username",
				},
			},
			IsUnique: true,
			Name:     "username_index",
		},
	}

	// ====> Create New Collection
	collConfig := &mongo.Collection{
		Connection:   conn,
		Database:     "rns_projections",
		Name:         "user_auth",
		SchemaStruct: &User{},
		Indexes:      indexConfigs,
	}
	c, err := mongo.EnsureCollection(collConfig)
	if err != nil {
		err = errors.Wrap(err, "Error creating DB-client")
		return nil, err
	}
	return &DB{
		collection: c,
	}, nil
}

func (d *DB) UserByUUID(uid uuid.UUID) (*User, error) {
	user := &User{
		UUID: uid,
	}

	findResults, err := d.collection.Find(user)
	if err != nil {
		err = errors.Wrap(err, "UserByUUID: Error getting user from Database")
		return nil, err
	}
	if len(findResults) == 0 {
		return nil, errors.New("UserByUUID: User not found")
	}

	resultUser := findResults[0].(*User)
	return resultUser, nil
}

func (d *DB) Login(user *User) (*User, error) {
	authUser := &User{
		Email:    user.Email,
		Username: user.Username,
	}

	findResults, err := d.collection.Find(authUser)
	if err != nil {
		err = errors.Wrap(err, "Login: Error getting user from Database")
		return nil, err
	}
	if len(findResults) == 0 {
		return nil, errors.New("Login: Invalid Credentials")
	}

	newUser := findResults[0].(*User)
	passErr := bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte(user.Password))
	if passErr != nil {
		return nil, errors.New("Login: Invalid Credentials")
	}

	return newUser, nil
}

func (d *DB) Collection() *mongo.Collection {
	return d.collection
}

package model

import (
	"github.com/gofrs/uuid"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type User struct {
	ID        objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UUID      uuid.UUID         `bson:"uuid,omitempty" json:"uuid,omitempty"`
	Email     string            `bson:"email,omitempty" json:"email,omitempty"`
	FirstName string            `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string            `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Username  string            `bson:"username,omitempty" json:"username,omitempty"`
	Password  string            `bson:"password,omitempty" json:"password,omitempty"`
	Role      string            `bson:"role,omitempty" json:"role,omitempty"`
}

// marshalUser is a simple
type marshalUser struct {
	ID        objectid.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UUID      string            `bson:"uuid,omitempty" json:"uuid,omitempty"`
	Email     string            `bson:"email,omitempty" json:"email,omitempty"`
	FirstName string            `bson:"first_name,omitempty" json:"first_name,omitempty"`
	LastName  string            `bson:"last_name,omitempty" json:"last_name,omitempty"`
	Username  string            `bson:"username,omitempty" json:"username,omitempty"`
	Password  string            `bson:"password,omitempty" json:"password,omitempty"`
	Role      string            `bson:"role,omitempty" json:"role,omitempty"`
}

func (u *User) MarshalBSON() ([]byte, error) {
	mu := &marshalUser{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		Role:      u.Role,
	}

	mu.ID = u.ID
	if u.UUID.String() != (uuid.UUID{}).String() {
		mu.UUID = u.UUID.String()
	}

	return bson.Marshal(mu)
}

func (u *User) UnmarshalBSON(in []byte) error {
	m := make(map[string]interface{})
	err := bson.Unmarshal(in, m)
	if err != nil {
		return err
	}
	u.ID = m["_id"].(objectid.ObjectID)

	u.UUID, err = uuid.FromString(m["uuid"].(string))
	if err != nil {
		return err
	}
	u.Email = m["email"].(string)
	u.FirstName = m["first_name"].(string)
	u.LastName = m["last_name"].(string)
	u.Username = m["username"].(string)
	u.Password = m["password"].(string)
	u.Role = m["role"].(string)

	return nil
}

package models

import (
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MongoUUID embedding UUID and adds bson support
type MongoUUID struct {
	uuid.UUID
}

// NewMongoUUID creates new random UUID
func NewMongoUUID() MongoUUID {
	return MongoUUID{uuid.New()}
}

// MarshalBSONValue implements mongo Marshaler interface
func (u MongoUUID) MarshalBSONValue() (bsontype.Type, []byte, error) {
	b, err := u.MarshalBinary()
	if err != nil {
		return bsontype.Binary, nil, err
	}
	return bsontype.Binary, bsoncore.AppendBinary(nil, 4, b), nil
}

// UnmarshalBSONValue implements mongo UnMarshaler interface
func (u *MongoUUID) UnmarshalBSONValue(t bsontype.Type, raw []byte) error {
	if t != bsontype.Binary {
		return errors.New("invalid format on unmarshal bson value")
	}

	_, data, _, ok := bsoncore.ReadBinary(raw)
	if !ok {
		return errors.New("not enough bytes to unmarshal bson value")
	}

	u.UnmarshalBinary(data)
	return nil
}

package models

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// Timestamp embedding time.Time with BSON support
type Timestamp struct {
	time.Time
}

// MarshalBSONValue implements the bsoncodec.ValueMarshaler interface.
func (t Timestamp) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bsonx.Time(t.Time).MarshalBSONValue()
}

// UnmarshalBSONValue implements the bsoncodec.ValueUnmarshaler interface.
func (t *Timestamp) UnmarshalBSONValue(bt bsontype.Type, raw []byte) error {
	switch bt {
	case bsontype.DateTime:
		rv := bson.RawValue{Type: bsontype.DateTime, Value: raw}
		t.Time = rv.Time()
		return nil
	}

	return fmt.Errorf("BSON type '%s' is not supported when unmarshalling Timestamp", bt.String())
}

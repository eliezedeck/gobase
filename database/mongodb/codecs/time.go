package codecs

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type CSharpNullTimeDecoder struct{}

func (d CSharpNullTimeDecoder) DecodeValue(ctx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	if !val.CanSet() {
		return errors.New("can't set the field value")
	}
	var (
		vrType   = vr.Type()
		t        time.Time
		outIsPtr = val.Kind() == reflect.Ptr
		outIsNil = false
	)

	switch vrType {
	case bsontype.String:
		if str, err := vr.ReadString(); err != nil {
			return err
		} else {
			t, err = time.Parse("2006-01-02T15:04:05.999Z07:00", str) // this was taken directly from bson/bsoncodec/time_codec.go
			if err != nil {
				return err
			}
		}

	case bsontype.DateTime:
		if dt, err := vr.ReadDateTime(); err != nil {
			return err
		} else {
			t = time.Unix(dt/1000, dt%1000*1000000)
		}

	case bsontype.EmbeddedDocument:
		dr, err := vr.ReadDocument()
		if err != nil {
			return err
		}
		v, rdr, err := dr.ReadElement()
		if err != nil {
			return err
		}
		if err = rdr.Skip(); err != nil { // point to the next element now, we were only interested with the first element
			return err
		}

		// This is the reason for this whole decoder
		if v != "_csharpnull" {
			return fmt.Errorf("cannot decode %v (embedded) into a %s", val.Type())
		}
		outIsNil = true

	case bsontype.Null:
		if err := vr.ReadNull(); err != nil {
			return err
		}
		outIsNil = true

	default:
		return fmt.Errorf("cannot decode %v into a %s", vrType, val.Type())
	}

	if outIsPtr {
		if outIsNil {
			val.Set(reflect.Zero(val.Type())) // return as nil pointer
		} else {
			val.Set(reflect.ValueOf(&t)) // return as pointer to t
		}
	} else {
		val.Set(reflect.ValueOf(t)) // return as struct
	}
	return nil
}

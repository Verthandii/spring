package form

import (
	"net/url"
	"reflect"

	"github.com/go-playground/form/v4"

	"github.com/Verthandii/spring/encoding"
)

const (
	// Name is form codec name
	Name = "x-www-form-urlencoded"
	// Null value string
	nullStr = "null"
)

var (
	encoder = form.NewEncoder()
	decoder = form.NewDecoder()
)

func init() {
	decoder.SetTagName("json")
	encoder.SetTagName("json")
	encoding.RegisterCodec(codec{encoder: encoder, decoder: decoder})
}

type codec struct {
	encoder *form.Encoder
	decoder *form.Decoder
}

func (c codec) Marshal(v interface{}) ([]byte, error) {
	var vs url.Values
	var err error
	switch m := v.(type) {
	case url.Values:
		vs = m
	default:
		vs, err = c.encoder.Encode(v)
		if err != nil {
			return nil, err
		}
	}
	for k, v := range vs {
		if len(v) == 0 {
			delete(vs, k)
		}
	}
	return []byte(vs.Encode()), nil
}

func (c codec) Unmarshal(data []byte, v interface{}) error {
	vs, err := url.ParseQuery(string(data))
	if err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		rv = rv.Elem()
	}

	return c.decoder.Decode(v, vs)
}

func (codec) Name() string {
	return Name
}

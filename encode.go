package plist

import (
	"fmt"
	"io"
	"log"
	"reflect"
)

// Encoder ...
type Encoder struct {
	w io.Writer
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w}
}

// Encode ...
func (e *Encoder) Encode(v interface{}) error {
	pval, err := e.marshal(reflect.ValueOf(v))
	if err != nil {
		return err
	}

	if err := newXMLEncoder(e.w).generateDocument(pval); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (e *Encoder) marshal(v reflect.Value) (*plistValue, error) {
	switch v.Kind() {
	case reflect.String:
		return &plistValue{String, v.String()}, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &plistValue{Integer, signedInt{uint64(v.Int()), true}}, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return &plistValue{Integer, signedInt{uint64(v.Uint()), false}}, nil
	case reflect.Float32, reflect.Float64:
		return &plistValue{Real, sizedFloat{v.Float(), v.Type().Bits()}}, nil
	case reflect.Bool:
		return &plistValue{Boolean, v.Bool()}, nil
	case reflect.Slice, reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			bytes := []byte(nil)
			if v.CanAddr() {
				bytes = v.Bytes()
			} else {
				bytes = make([]byte, v.Len())
				reflect.Copy(reflect.ValueOf(bytes), v)
			}
			return &plistValue{Data, bytes}, nil
		}
		subvalues := make([]*plistValue, v.Len())
		for idx, length := 0, v.Len(); idx < length; idx++ {
			subpval, err := e.marshal(v.Index(idx))
			if err != nil {
				return nil, err
			}
			if subpval != nil {
				subvalues[idx] = subpval
			}
		}
		return &plistValue{Array, subvalues}, nil
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			//TODO return marshalerr
			panic(v.Type())
		}

		l := v.Len()
		dict := &dictionary{
			m: make(map[string]*plistValue, l),
		}
		for _, keyv := range v.MapKeys() {
			subpval, err := e.marshal(v.MapIndex(keyv))
			if err != nil {
				return nil, err
			}
			if subpval != nil {
				dict.m[keyv.String()] = subpval
			}
		}
		return &plistValue{Dictionary, dict}, nil
	default:
		fmt.Println(v.Kind())
		panic("not implemented")
	}
}

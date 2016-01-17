package plist

import (
	"os"
	"testing"
	"time"
)

func TestGenDoc(t *testing.T) {
	err := newXMLEncoder(os.Stdout).generateDocument(&plistValue{String, "foo"})
	if err != nil {
		t.Fatal(err)
	}
	err = newXMLEncoder(os.Stdout).generateDocument(&plistValue{Boolean, true})
	if err != nil {
		t.Fatal(err)
	}
	err = newXMLEncoder(os.Stdout).generateDocument(&plistValue{Integer, signedInt{4, false}})
	if err != nil {
		t.Fatal(err)
	}
	err = newXMLEncoder(os.Stdout).generateDocument(&plistValue{Date, time.Now()})
	if err != nil {
		t.Fatal(err)
	}
	err = newXMLEncoder(os.Stdout).generateDocument(&plistValue{Dictionary, &dictionary{
		m: map[string]*plistValue{
			"foo": &plistValue{Integer, signedInt{uint64(1), false}},
			"bar": &plistValue{String, "foobar"},
			"baz": &plistValue{Boolean, true},
		},
	}})
	if err != nil {
		t.Fatal(err)
	}
}

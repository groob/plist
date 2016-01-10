package plist

import (
	"os"
	"testing"
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
}

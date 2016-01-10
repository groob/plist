package plist

import (
	"log"
	"os"
	"testing"
)

func TestGenDoc(t *testing.T) {
	err := newXMLEncoder(os.Stdout).generateDocument(&plistValue{String, "foo"})
	if err != nil {
		log.Fatal(err)
	}
}

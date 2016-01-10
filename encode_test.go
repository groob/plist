package plist

import (
	"os"
	"testing"
)

func TestEncode(t *testing.T) {
	input := "test"
	err := NewEncoder(os.Stdout).Encode(input)
	if err != nil {
		t.Fatal(err)
	}
}

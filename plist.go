package plist

import (
	"fmt"
	"sort"
)

type plistKind uint

const (
	Invalid plistKind = iota
	Dictionary
	Array
	String
	Integer
	Real
	Boolean
	Data
	Date
)

var plistKindNames = map[plistKind]string{
	Invalid:    "invalid",
	Dictionary: "dictionary",
	Array:      "array",
	String:     "string",
	Integer:    "integer",
	Real:       "real",
	Boolean:    "boolean",
	Data:       "data",
	Date:       "date",
}

type plistValue struct {
	kind  plistKind
	value interface{}
}

type signedInt struct {
	value  uint64
	signed bool
}

func (sf signedInt) String() string {
	return fmt.Sprintf("%v", sf.value)
}

type sizedFloat struct {
	value float64
	bits  int
}

func (sf sizedFloat) String() string {
	return fmt.Sprintf("%v", sf.value)
}

type array []*plistValue

func (a array) String() string {
	return "array"
}

type dictionary struct {
	count  int
	m      map[string]*plistValue
	keys   sort.StringSlice
	values []*plistValue
}

func (d *dictionary) Len() int {
	return len(d.m)
}

func (d *dictionary) Less(i, j int) bool {
	return d.keys.Less(i, j)
}

func (d *dictionary) Swap(i, j int) {
	d.keys.Swap(i, j)
	d.values[i], d.values[j] = d.values[j], d.values[i]
}

func (d *dictionary) populateArrays() {
	d.keys = make([]string, len(d.m))
	d.values = make([]*plistValue, len(d.m))
	i := 0
	for k, v := range d.m {
		d.keys[i] = k
		d.values[i] = v
		i++
	}
	sort.Sort(d)
}

func (d *dictionary) String() string {
	return "dict"
}

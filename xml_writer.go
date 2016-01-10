package plist

import (
	"encoding/xml"
	"io"
	"log"
)

const xmlDOCTYPE = `<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">`

type xmlEncoder struct {
	writer io.Writer
	*xml.Encoder
}

func newXMLEncoder(w io.Writer) *xmlEncoder {
	return &xmlEncoder{w, xml.NewEncoder(w)}
}

func (e *xmlEncoder) generateDocument(pval *plistValue) error {
	_, err := e.writer.Write([]byte(xml.Header))
	if err != nil {
		log.Fatal(err)
	}
	_, err = e.writer.Write([]byte(xmlDOCTYPE))
	if err != nil {
		log.Fatal(err)
	}
	plistStartElement := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "plist",
		},
		Attr: []xml.Attr{{
			Name: xml.Name{
				Space: "",
				Local: "version"},
			Value: "1.0"},
		},
	}

	if err := e.EncodeToken(plistStartElement); err != nil {
		return err
	}
	// do stuff here
	if err := e.writePlistValue(pval); err != nil {
		return err
	}
	if err := e.EncodeToken(plistStartElement.End()); err != nil {
		return err
	}

	if err := e.Flush(); err != nil {
		return err
	}
	return nil
}
func (e *xmlEncoder) writePlistValue(pval *plistValue) error {
	var key string
	encodedValue := pval.value
	switch pval.kind {
	case String:
		key = "string"
	default:
		panic(pval.kind)

	}
	if key == "" {
		panic("nil key")
	}
	err := e.EncodeElement(encodedValue, xml.StartElement{Name: xml.Name{Local: key}})
	e.Flush()
	return err
}

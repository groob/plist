package plist

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"time"
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
	if err := e.Flush(); err != nil {
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
	switch pval.kind {
	case String:
		return e.writeStringValue(pval)
	case Boolean:
		return e.writeBoolValue(pval)
	case Integer:
		return e.writeIntegerValue(pval)
	case Dictionary:
		return e.writeDictionaryValue(pval)
	case Date:
		return e.writeDateValue(pval)
	default:
		panic(pval.kind)
	}
	return nil
}

func (e *xmlEncoder) writeDictionaryValue(pval *plistValue) error {
	encodedValue := pval.value
	dictStartElement := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "dict",
		}}
	if err := e.EncodeToken(dictStartElement); err != nil {
		return err
	}
	if err := e.Flush(); err != nil {
		return err
	}
	dict := encodedValue.(*dictionary)
	dict.populateArrays()
	for i, k := range dict.keys {
		if err := e.EncodeElement(k, xml.StartElement{Name: xml.Name{Local: "key"}}); err != nil {
			return err
		}
		if err := e.writePlistValue(dict.values[i]); err != nil {
			return err
		}
	}
	if err := e.EncodeToken(dictStartElement.End()); err != nil {
		return err
	}
	return e.Flush()
}
func (e *xmlEncoder) writeStringValue(pval *plistValue) error {
	return e.EncodeElement(pval.value, xml.StartElement{Name: xml.Name{Local: "string"}})
}

func (e *xmlEncoder) writeBoolValue(pval *plistValue) error {
	// EncodeElement results in <true></true> instead of <true/>
	// use writer to write self closing tags
	b := pval.value.(bool)
	_, err := e.writer.Write([]byte(fmt.Sprintf("<%t/>", b)))
	if err != nil {
		return err
	}
	return nil
}

func (e *xmlEncoder) writeIntegerValue(pval *plistValue) error {
	encodedValue := pval.value
	if pval.value.(signedInt).signed {
		encodedValue = int64(pval.value.(signedInt).value)
	} else {
		encodedValue = pval.value.(signedInt).value
	}
	return e.EncodeElement(encodedValue, xml.StartElement{Name: xml.Name{Local: "integer"}})
}

func (e *xmlEncoder) writeDateValue(pval *plistValue) error {
	encodedValue := pval.value.(time.Time).In(time.UTC).Format(time.RFC3339)
	return e.EncodeElement(encodedValue, xml.StartElement{Name: xml.Name{Local: "date"}})
}

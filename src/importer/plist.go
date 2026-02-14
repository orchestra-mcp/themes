package importer

import "encoding/xml"

// tmTheme plist XML structures.

type plistRoot struct {
	XMLName xml.Name  `xml:"plist"`
	Dict    plistDict `xml:"dict"`
}

type plistDict struct {
	Items []plistItem
}

type plistItem struct {
	Key   string
	Value plistValue
}

type plistValue struct {
	String string
	Array  []plistDict
	Dict   *plistDict
}

// UnmarshalXML parses plist <dict> elements into key-value pairs.
func (d *plistDict) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	d.Items = nil
	var currentKey string
	expectKey := true

	for {
		tok, err := dec.Token()
		if err != nil {
			return err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			switch t.Name.Local {
			case "key":
				var s string
				if err := dec.DecodeElement(&s, &t); err != nil {
					return err
				}
				currentKey = s
				expectKey = false
			case "string":
				var s string
				if err := dec.DecodeElement(&s, &t); err != nil {
					return err
				}
				d.Items = append(d.Items, plistItem{
					Key:   currentKey,
					Value: plistValue{String: s},
				})
				expectKey = true
			case "dict":
				var child plistDict
				if err := dec.DecodeElement(&child, &t); err != nil {
					return err
				}
				d.Items = append(d.Items, plistItem{
					Key:   currentKey,
					Value: plistValue{Dict: &child},
				})
				expectKey = true
			case "array":
				arr, err := decodeArray(dec)
				if err != nil {
					return err
				}
				d.Items = append(d.Items, plistItem{
					Key:   currentKey,
					Value: plistValue{Array: arr},
				})
				expectKey = true
			default:
				if err := dec.Skip(); err != nil {
					return err
				}
				_ = expectKey
			}
		case xml.EndElement:
			return nil
		}
	}
}

// decodeArray parses a plist <array> containing <dict> elements.
func decodeArray(dec *xml.Decoder) ([]plistDict, error) {
	var result []plistDict
	for {
		tok, err := dec.Token()
		if err != nil {
			return nil, err
		}

		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "dict" {
				var d plistDict
				if err := dec.DecodeElement(&d, &t); err != nil {
					return nil, err
				}
				result = append(result, d)
			} else {
				if err := dec.Skip(); err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			return result, nil
		}
	}
}

// dictGet retrieves a string value by key from a plistDict.
func dictGet(d *plistDict, key string) string {
	for _, item := range d.Items {
		if item.Key == key {
			return item.Value.String
		}
	}
	return ""
}

// dictGetDict retrieves a child dict by key.
func dictGetDict(d *plistDict, key string) *plistDict {
	for _, item := range d.Items {
		if item.Key == key {
			return item.Value.Dict
		}
	}
	return nil
}

// dictGetArray retrieves a child array by key.
func dictGetArray(d *plistDict, key string) []plistDict {
	for _, item := range d.Items {
		if item.Key == key {
			return item.Value.Array
		}
	}
	return nil
}

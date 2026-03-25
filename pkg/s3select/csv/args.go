/*
 * MinIO Cloud Storage, (C) 2019 MinIO, Inc.
 * PGG Obstor, (C) 2021-2026 PGG, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package csv

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

const (
	none   = "none"
	use    = "use"
	ignore = "ignore"

	defaultRecordDelimiter      = "\n"
	defaultFieldDelimiter       = ","
	defaultQuoteCharacter       = `"`
	defaultQuoteEscapeCharacter = `"`
	defaultCommentCharacter     = "#"

	asneeded = "asneeded"
)

// ReaderArgs - represents elements inside <InputSerialization><CSV> in request XML.
type ReaderArgs struct {
	FileHeaderInfo             string `xml:"FileHeaderInfo"`
	RecordDelimiter            string `xml:"RecordDelimiter"`
	FieldDelimiter             string `xml:"FieldDelimiter"`
	QuoteCharacter             string `xml:"QuoteCharacter"`
	QuoteEscapeCharacter       string `xml:"QuoteEscapeCharacter"`
	CommentCharacter           string `xml:"Comments"`
	AllowQuotedRecordDelimiter bool   `xml:"AllowQuotedRecordDelimiter"`
	unmarshaled                bool
}

// IsEmpty - returns whether reader args is empty or not.
func (args *ReaderArgs) IsEmpty() bool {
	return !args.unmarshaled
}

func requireSingleRune(s, fieldName string) error {
	if utf8.RuneCountInString(s) > 1 {
		return fmt.Errorf("%s must be a single character, got %q", fieldName, s)
	}
	return nil
}

func validFileHeaderInfo(s string) bool {
	return s == none || s == use || s == ignore
}

// UnmarshalXML - decodes XML data.
func (args *ReaderArgs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) (err error) {
	args.FileHeaderInfo = none
	args.RecordDelimiter = defaultRecordDelimiter
	args.FieldDelimiter = defaultFieldDelimiter
	args.QuoteCharacter = defaultQuoteCharacter
	args.QuoteEscapeCharacter = defaultQuoteEscapeCharacter
	args.CommentCharacter = defaultCommentCharacter
	args.AllowQuotedRecordDelimiter = false

	for {
		// Read tokens from the XML document in a stream.
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		se, ok := t.(xml.StartElement)
		if !ok {
			continue
		}

		tagName := se.Name.Local
		if tagName == "AllowQuotedRecordDelimiter" {
			var b bool
			if err = d.DecodeElement(&b, &se); err != nil {
				return err
			}
			args.AllowQuotedRecordDelimiter = b
			continue
		}

		var s string
		if err = d.DecodeElement(&s, &se); err != nil {
			return err
		}

		if err := args.setReaderField(tagName, s); err != nil {
			return err
		}
	}

	args.unmarshaled = true
	return nil
}

func (args *ReaderArgs) setReaderField(tagName, value string) error {
	switch tagName {
	case "FileHeaderInfo":
		v := strings.ToLower(value)
		if v != "" && !validFileHeaderInfo(v) {
			return fmt.Errorf("FileHeaderInfo %q is not recognized", value)
		}
		if v != "" {
			args.FileHeaderInfo = v
		}
	case "RecordDelimiter":
		if value != "" {
			args.RecordDelimiter = value
		}
	case "FieldDelimiter":
		if value != "" {
			args.FieldDelimiter = value
		}
	case "QuoteCharacter":
		if err := requireSingleRune(value, "QuoteCharacter"); err != nil {
			return err
		}
		args.QuoteCharacter = value
	case "QuoteEscapeCharacter":
		runeCount := utf8.RuneCountInString(value)
		if runeCount > 1 {
			return fmt.Errorf("QuoteEscapeCharacter must be a single character, got %q", value)
		}
		if runeCount == 0 {
			args.QuoteEscapeCharacter = defaultQuoteEscapeCharacter
		} else {
			args.QuoteEscapeCharacter = value
		}
	case "Comments":
		if value != "" {
			args.CommentCharacter = value
		}
	default:
		return fmt.Errorf("unknown CSV input element %q", tagName)
	}
	return nil
}

// WriterArgs - represents elements inside <OutputSerialization><CSV/> in request XML.
type WriterArgs struct {
	QuoteFields          string `xml:"QuoteFields"`
	RecordDelimiter      string `xml:"RecordDelimiter"`
	FieldDelimiter       string `xml:"FieldDelimiter"`
	QuoteCharacter       string `xml:"QuoteCharacter"`
	QuoteEscapeCharacter string `xml:"QuoteEscapeCharacter"`
	unmarshaled          bool
}

// IsEmpty - returns whether writer args is empty or not.
func (args *WriterArgs) IsEmpty() bool {
	return !args.unmarshaled
}

// UnmarshalXML - decodes XML data.
func (args *WriterArgs) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	args.QuoteFields = asneeded
	args.RecordDelimiter = defaultRecordDelimiter
	args.FieldDelimiter = defaultFieldDelimiter
	args.QuoteCharacter = defaultQuoteCharacter
	args.QuoteEscapeCharacter = defaultQuoteCharacter

	for {
		// Read tokens from the XML document in a stream.
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		se, ok := t.(xml.StartElement)
		if !ok {
			continue
		}

		var s string
		if err = d.DecodeElement(&s, &se); err != nil {
			return err
		}

		if err := args.setWriterField(se.Name.Local, s); err != nil {
			return err
		}
	}

	args.unmarshaled = true
	return nil
}

func (args *WriterArgs) setWriterField(tagName, value string) error {
	switch tagName {
	case "QuoteFields":
		args.QuoteFields = strings.ToLower(value)
	case "RecordDelimiter":
		args.RecordDelimiter = value
	case "FieldDelimiter":
		args.FieldDelimiter = value
	case "QuoteCharacter":
		runeCount := utf8.RuneCountInString(value)
		if runeCount > 1 {
			return fmt.Errorf("QuoteCharacter must be a single character, got %q", value)
		}
		if runeCount == 0 {
			args.QuoteCharacter = "\x00"
		} else {
			args.QuoteCharacter = value
		}
	case "QuoteEscapeCharacter":
		runeCount := utf8.RuneCountInString(value)
		if runeCount > 1 {
			return fmt.Errorf("QuoteEscapeCharacter must be a single character, got %q", value)
		}
		if runeCount == 0 {
			args.QuoteEscapeCharacter = defaultQuoteEscapeCharacter
		} else {
			args.QuoteEscapeCharacter = value
		}
	default:
		return fmt.Errorf("unknown CSV output element %q", tagName)
	}
	return nil
}

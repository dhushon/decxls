

// Copyright 2014 Jonathan Picques. All rights reserved.
// Use of this source code is governed by a MIT license
// The license can be found in the LICENSE file.

// The GoXLS package aims to provide easy XLS serialization and deserialization to the golang programming language

// sourced here under MIT license: https://github.com/gocarina/gocsv/blob/master/csv.go

package decxls

import (
	//"encoding/xls"
	//"fmt"
	"errors"
	"io"
	"reflect"
	"sync"
	"github.com/360EntSecGroup-Skylar/excelize"
)

// TagSeparator defines seperator string for multiple xls tags in struct fields
var TagSeparator = ","

// UseColumnHeader suggests that we try and pull the column name (row 0 in excelize) vs. row 1 in excellize = first visible row
var UseColumnHeader = true

// FailIfUnmatchedStructTags indicates whether it is considered an error when there is an unmatched
// struct tag.
var FailIfUnmatchedStructTags = false

// FailIfDoubleHeaderNames indicates whether it is considered an error when a header name is repeated
// in the xls header.
var FailIfDoubleHeaderNames = false

// ShouldAlignDuplicateHeadersWithStructFieldOrder indicates whether we should align duplicate XLS
// headers per their alignment in the struct definition.
var ShouldAlignDuplicateHeadersWithStructFieldOrder = false

// TagName defines key in the struct field's tag to scan
var TagName = "xls"

// Normalizer is a function that takes and returns a string. It is applied to
// struct and header field values before they are compared. It can be used to alter
// names for comparison. For instance, you could allow case insensitive matching
// or convert '-' to '_'.
type Normalizer func(string) string

//ErrorHandler - 
type ErrorHandler func(*ParseError) bool

// normalizeName function initially set to a nop Normalizer.
var normalizeName = DefaultNameNormalizer()

// DefaultNameNormalizer is a nop Normalizer.
func DefaultNameNormalizer() Normalizer { return func(s string) string { return s } }

// SetHeaderNormalizer sets the normalizer used to normalize struct and header field names.
func SetHeaderNormalizer(f Normalizer) {
	normalizeName = f
	// Need to clear the cache hen the header normalizer changes.
	structInfoCache = sync.Map{}
}

// --------------------------------------------------------------------------
// XLSReader used to parse XLS

var selfXLSReader = DefaultXLSReader

// DefaultXLSReader is the default XLS reader used to parse XLS (cf. xls.NewReader)
func DefaultXLSReader(in io.Reader) XLSReader {
	return nil //xls.NewReader(in)
}

// SetXLSReader sets the XLS reader used to parse XLS.
func SetXLSReader(xlsReader func(io.Reader) XLSReader) {
	selfXLSReader = xlsReader
}

func getXLSReader(in io.Reader) XLSReader {
	return selfXLSReader(in)
}

// --------------------------------------------------------------------------
// Unmarshal functions

// Unmarshal

// UnmarshalFile parses the XLS from the file in the interface.
func UnmarshalFile(filename string, sheetName string, out interface{}) error {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return err
	}
	return UnmarshalExcelize(f, sheetName, out)
}

// UnmarshalExcelize will take an excelize opened XML file and marshal to an interface
func UnmarshalExcelize(f *excelize.File, sheet string, out interface{}) error {
	//error check input
	if reflect.ValueOf(f).IsNil() {
		return errors.New("XLS file is not open, or pointer is nil")
	}
	sn, err := findSheet(f, sheet, out)
	if err != nil {
		return err
	}
	//log.Printf("found \"%s\" sheet in document\n", sn)
	r,err := f.Rows(sn)
	if err != nil{
		return err
	}
	err = readToExcelizeWithErrorHandler(r, nil, out )
	return nil
}

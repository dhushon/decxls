package decxls

import (
	"fmt"
	"errors"
)

// A ParseError is returned for parsing errors.
// Line numbers are 1-indexed and columns are 0-indexed.
type ParseError struct {
	StartLine int   // Line where the record starts
	Line      int   // Line where the error occurred
	Column    int   // Column (rune index) where the error occurred
	Err       error // The actual error
}

func (e *ParseError) Error() string {
	if e.Err == ErrFieldCount {
		return fmt.Sprintf("record on line %d: %v", e.Line, e.Err)
	}
	if e.StartLine != e.Line {
		return fmt.Sprintf("record on line %d; parse error on line %d, column %d: %v", e.StartLine, e.Line, e.Column, e.Err)
	}
	return fmt.Sprintf("parse error on line %d, column %d: %v", e.Line, e.Column, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

// These are the errors that can be returned in ParseError.Err.
var (
	ErrQuote         = errors.New("extraneous or missing \" in quoted-field")
	ErrFieldCount    = errors.New("wrong number of fields")
)

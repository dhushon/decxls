package decxls

// prototype here: https://github.com/gocarina/goxls/blob/master/decode.go
import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

//Decoder interface
type Decoder interface {
	getXLSRows() ([][]string, error)
}

// SimpleDecoder interface
type SimpleDecoder interface {
	getXLSRow() ([]string, error)
	getXLSRows() ([][]string, error)
}

// XLSReader interface
type XLSReader interface {
	Read() ([]string, error)
	ReadAll() ([][]string, error)
}

var (
	//ErrEmptyFile is an error type
	ErrEmptyFile = errors.New("empty xls file given")
	//ErrNoStructTags error type if no tags are found
	ErrNoStructTags = errors.New("no xls struct tags found")
)


func readToExcelizeWithErrorHandler(r *excelize.Rows, errHandler ErrorHandler, out interface{}) error {
	// reflect on the interface{} provided for the resultant
	// get its concrete reference and type
	outValue, outType := getConcreteReflectValueAndType(out) // Get the concrete type (not pointer) (Slice<?> or Array<?>)
	//fmt.Printf("OutType: %s\n", outType.String())
	if err := ensureOutType(outType); err != nil {
		return err
	}
	// go thru the struct and pull the types
	// outInnerIsPointer deals with the &[]*struc where struct is the outInner
	outInnerIsPointer, outInnerType := getConcreteContainerInnerType(outType) // Get the concrete inner type (not pointer) (Container<"?">)
	if err := ensureOutInnerType(outInnerType); err != nil {
		return err
	}
	
	outInnerStructInfo := getStructInfo(outInnerType) // Get the inner struct info to get XLS annotations
	if len(outInnerStructInfo.Fields) == 0 {
		return ErrEmptyFile
	}
	//fmt.Printf("struct xls: tags discovered %v\n", outInnerStructInfo.Fields)

	// start reading the input to extract headers
	if !r.Next() {
		return errors.New("expected Rows, none found")
	}
	headers, err := r.Columns()
	if err != nil {
		return err
	}

	xlsHeadersLabels := make(map[int]*fieldInfo, len(outInnerStructInfo.Fields)) // Used to store the correspondance header <-> position in XLS

	headerCount := map[string]int{}
	for i, xlsColumnHeader := range headers {
		curHeaderCount := headerCount[xlsColumnHeader]
		if fieldInfo := getXLSFieldPosition(xlsColumnHeader, outInnerStructInfo, curHeaderCount); fieldInfo != nil {
			xlsHeadersLabels[i] = fieldInfo
			if ShouldAlignDuplicateHeadersWithStructFieldOrder {
				curHeaderCount++
				headerCount[xlsColumnHeader] = curHeaderCount
			}
		}
	}

	if FailIfUnmatchedStructTags {
		if err := maybeMissingStructFields(outInnerStructInfo.Fields, headers); err != nil {
			return err
		}
	}
	if FailIfDoubleHeaderNames {
		if err := maybeDoubleHeaderNames(headers); err != nil {
			return err
		}
	}
	// read spreadsheet data rows into struct
	var withFieldsOK bool
	var fieldTypeUnmarshallerWithKeys TypeUnmarshalXLSWithFields
	var i int = -1 // skip header row

	for r.Next() {
		xlsRow, err := r.Columns()
		if err != nil {
			return err
		}
		i++ // can increment row at bottom
		objectIface := reflect.New(outInnerType).Interface()
		outInner := createNewOutInner(outInnerIsPointer, outInnerType)
		
		for j, xlsColumnContent := range xlsRow {
			if fieldInfo, ok := xlsHeadersLabels[j]; ok { // Position found accordingly to header name
				if outValue.CanInterface() { //make sure thatthe value is not Private
					//fmt.Printf("Parsing Row: %d Column %d : %s\n", i, j, fieldInfo.getFirstKey())
					fieldTypeUnmarshallerWithKeys, withFieldsOK = objectIface.(TypeUnmarshalXLSWithFields)
					if withFieldsOK { // if there is an field based interface driven unmarshaller, use it
						if err := fieldTypeUnmarshallerWithKeys.UnmarshalXLSWithFields(fieldInfo.getFirstKey(), xlsColumnContent); err != nil {
							// tell the developer as much as we can about the error
							parseError := ParseError{
								Line:   i + 2, //add 2 to account for the header & 0-indexing of arrays
								Column: j + 1,
								Err:    err,
							}
							return &parseError
						}
						continue
					}
				}

				if err := setInnerField(&outInner, outInnerIsPointer, fieldInfo.IndexChain, xlsColumnContent, fieldInfo.omitEmpty); err != nil { // Set field of struct
					// tell the developer as much as we can about the error
					parseError := ParseError{
						Line:   i + 2, //add 2 to account for the header & 0-indexing of arrays
						Column: j + 1,
						Err:    err,
					}
					if errHandler == nil || !errHandler(&parseError) {
						return &parseError
					}
				}
			}
		}

		if withFieldsOK {
			reflectedObject := reflect.ValueOf(objectIface)
			outValue = reflectedObject.Elem()
		}
		// append new row to main slice structure
		outValue.Set(reflect.Append(outValue, outInner))
	}
	return nil
}

func mismatchStructFields(structInfo []fieldInfo, headers []string) []string {
	missing := make([]string, 0)
	if len(structInfo) == 0 {
		return missing
	}

	headerMap := make(map[string]struct{}, len(headers))
	for idx := range headers {
		headerMap[headers[idx]] = struct{}{}
	}

	for _, info := range structInfo {
		found := false
		for _, key := range info.keys {
			if _, ok := headerMap[key]; ok {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, info.keys...)
		}
	}
	return missing
}

func mismatchHeaderFields(structInfo []fieldInfo, headers []string) []string {
	missing := make([]string, 0)
	if len(headers) == 0 {
		return missing
	}

	keyMap := make(map[string]struct{})
	for _, info := range structInfo {
		for _, key := range info.keys {
			keyMap[key] = struct{}{}
		}
	}

	for _, header := range headers {
		if _, ok := keyMap[header]; !ok {
			missing = append(missing, header)
		}
	}
	return missing
}

func maybeMissingStructFields(structInfo []fieldInfo, headers []string) error {
	missing := mismatchStructFields(structInfo, headers)
	if len(missing) != 0 {
		return fmt.Errorf("found unmatched struct field with tags %v", missing)
	}
	return nil
}

// Check that no header name is repeated twice
func maybeDoubleHeaderNames(headers []string) error {
	headerMap := make(map[string]bool, len(headers))
	for _, v := range headers {
		if _, ok := headerMap[v]; ok {
			return fmt.Errorf("repeated header name: %v", v)
		}
		headerMap[v] = true
	}
	return nil
}

// apply normalizer func to headers
func normalizeHeaders(headers []string) []string {
	out := make([]string, len(headers))
	for i, h := range headers {
		out[i] = normalizeName(h)
	}
	return out
}

//Try and reduce complexity in Sheet Names... 
func normalizeString(s string) string {
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(s, "")
}

func findSheet(f *excelize.File, sheet string, out interface{}) (string, error) {
	// see if we can find the sheet by name
	sh := f.GetSheetList()
	if len(sh) == 0 {
		return "", errors.New("XLS document may be empty or a non XLS document as there are no sheets")
	}

	var eSheet string = ""
	for _, s := range sh {
		if s == sheet {
			return s, nil
		}
	}

	// if not found try normalization to find a match
	sheet = normalizeString(sheet)
	for _, s := range sh {
		if normalizeString(s) == sheet {
			if eSheet != "" {
				return eSheet, errors.New("found a duplicate sheet because of normalization")
			}
			eSheet = s
		}
	}
	if eSheet == "" {
		return eSheet, fmt.Errorf("sheet not found only these %s", sh)
	}
	return eSheet, nil
}

// Check if the outType is an array or a slice
func ensureOutType(outType reflect.Type) error {
	switch outType.Kind() {
	case reflect.Slice:
		fallthrough
	case reflect.Chan:
		fallthrough
	case reflect.Array:
		return nil
	}
	return fmt.Errorf("cannot use " + outType.String() + ", only slice or array supported")
}

// Check if the outInnerType is of type struct
func ensureOutInnerType(outInnerType reflect.Type) error {
	switch outInnerType.Kind() {
	case reflect.Struct:
		return nil
	}
	return fmt.Errorf("cannot use " + outInnerType.String() + ", only struct supported")
}

func getXLSFieldPosition(key string, structInfo *structInfo, curHeaderCount int) *fieldInfo {
	matchedFieldCount := 0
	for _, field := range structInfo.Fields {
		if field.matchesKey(key) {
			if matchedFieldCount >= curHeaderCount {
				return &field
			}
			matchedFieldCount++
		}
	}
	return nil
}

func createNewOutInner(outInnerWasPointer bool, outInnerType reflect.Type) reflect.Value {
	if outInnerWasPointer {
		return reflect.New(outInnerType)
	}
	return reflect.New(outInnerType).Elem()
}

func setInnerField(outInner *reflect.Value, outInnerWasPointer bool, index []int, value string, omitEmpty bool) error {
	oi := *outInner
	if outInnerWasPointer {
		// initialize nil pointer
		if oi.IsNil() {
			setField(oi, "", omitEmpty)
		}
		oi = outInner.Elem()
	}
	// because pointers can be nil need to recurse one index at a time and perform nil check
	if len(index) > 1 {
		nextField := oi.Field(index[0])
		return setInnerField(&nextField, nextField.Kind() == reflect.Ptr, index[1:], value, omitEmpty)
	}
	return setField(oi.FieldByIndex(index), value, omitEmpty)
}


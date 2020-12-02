package decxls

import (
	"encoding"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"encoding/json"
)

// --------------------------------------------------------------------------
// Conversion interfaces

// TypeUnmarshaller is implemented by any value that has an UnmarshalXLS method
// This converter is used to convert a string to your value representation of that string
type TypeUnmarshaller interface {
	UnmarshalXLS(string) error
}

// TypeUnmarshalXLSWithFields can be implemented on whole structs to allow for whole structures to customized internal vs one off fields
type TypeUnmarshalXLSWithFields interface {
	UnmarshalXLSWithFields(key, value string) error
}

// NoUnmarshalFuncError is the custom error type to be raised in case there is no unmarshal function defined on type
type NoUnmarshalFuncError struct {
	msg string
}

func (e NoUnmarshalFuncError) Error() string {
	return e.msg
}

// --------------------------------------------------------------------------
// Conversion helpers

func toString(in interface{}) (string, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		return inValue.String(), nil
	case reflect.Bool:
		b := inValue.Bool()
		if b {
			return "true", nil
		}
		return "false", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%v", inValue.Uint()), nil
	case reflect.Float32:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 32), nil
	case reflect.Float64:
		return strconv.FormatFloat(inValue.Float(), byte('f'), -1, 64), nil
	}
	return "", fmt.Errorf("No known conversion from " + inValue.Type().String() + " to string")
}

func toBool(in interface{}) (bool, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := inValue.String()
		s = strings.TrimSpace(s)
		if strings.EqualFold(s, "yes") {
			return true, nil
		} else if strings.EqualFold(s, "no") || s == "" {
			return false, nil
		} else {
			return strconv.ParseBool(s)
		}
	case reflect.Bool:
		return inValue.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i := inValue.Int()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i := inValue.Uint()
		if i != 0 {
			return true, nil
		}
		return false, nil
	case reflect.Float32, reflect.Float64:
		f := inValue.Float()
		if f != 0 {
			return true, nil
		}
		return false, nil
	}
	return false, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to bool")
}

func toInt(in interface{}) (int64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		out := strings.SplitN(s, ".", 2)
		return strconv.ParseInt(out[0], 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return inValue.Int(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return int64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to int")
}

func toUint(in interface{}) (uint64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}

		// support the float input
		if strings.Contains(s, ".") {
			f, err := strconv.ParseFloat(s, 64)
			if err != nil {
				return 0, err
			}
			return uint64(f), nil
		}
		return strconv.ParseUint(s, 0, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return uint64(inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return inValue.Uint(), nil
	case reflect.Float32, reflect.Float64:
		return uint64(inValue.Float()), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to uint")
}

func toFloat(in interface{}) (float64, error) {
	inValue := reflect.ValueOf(in)

	switch inValue.Kind() {
	case reflect.String:
		s := strings.TrimSpace(inValue.String())
		if s == "" {
			return 0, nil
		}
		return strconv.ParseFloat(s, 64)
	case reflect.Bool:
		if inValue.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(inValue.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(inValue.Uint()), nil
	case reflect.Float32, reflect.Float64:
		return inValue.Float(), nil
	}
	return 0, fmt.Errorf("No known conversion from " + inValue.Type().String() + " to float")
}

func setField(field reflect.Value, value string, omitEmpty bool) error {
	if field.Kind() == reflect.Ptr {
		if omitEmpty && value == "" {
			return nil
		}
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		field = field.Elem()
	}
    //fmt.Printf("field.Interface().(type) %s and (kind) %s\n",field.Type(),field.Kind())
	switch field.Interface().(type) {
	case string:
		s, err := toString(value)
		if err != nil {
			return err
		}
		field.SetString(s)
	case bool:
		b, err := toBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
	case int, int8, int16, int32, int64:
		i, err := toInt(value)
		if err != nil {
			return err
		}
		field.SetInt(i)
	case uint, uint8, uint16, uint32, uint64:
		ui, err := toUint(value)
		if err != nil {
			return err
		}
		field.SetUint(ui)
	case float32, float64:
		f, err := toFloat(value)
		if err != nil {
			return err
		}
		field.SetFloat(f)
	default:
		// for nested structs... we might want to get the Unmarshal associated with the struct
		// may need to 
		// Not a native type, check for unmarshal method
		if err := unmarshall(field, value); err != nil {
			if _, ok := err.(NoUnmarshalFuncError); !ok {
				return err
			}
			// Could not unmarshal, check for kind, e.g. renamed type from basic type
			switch field.Kind() {
			case reflect.String:
				s, err := toString(value)
				if err != nil {
					return err
				}
				field.SetString(s)
			case reflect.Bool:
				b, err := toBool(value)
				if err != nil {
					return err
				}
				field.SetBool(b)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				i, err := toInt(value)
				if err != nil {
					return err
				}
				field.SetInt(i)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				ui, err := toUint(value)
				if err != nil {
					return err
				}
				field.SetUint(ui)
			case reflect.Float32, reflect.Float64:
				f, err := toFloat(value)
				if err != nil {
					return err
				}
				field.SetFloat(f)
			case reflect.Slice, reflect.Struct:
				err := json.Unmarshal([]byte(value), field.Addr().Interface())
				if err != nil {
					return err
				}
			default:
				return err
			}
		} else {
			return nil
		}
	}
	return nil
}

// --------------------------------------------------------------------------
// Un/serializations helpers

// --------------------------------------------------------------------------
// Un/serializations helpers

func canMarshal(t reflect.Type) bool {
	// unless it implements marshalXLS. Structs that implement this
	// should result in one value and not have their fields exposed
	_, canMarshalXLS := t.MethodByName("MarshalXLS")
	return canMarshalXLS
}

func unmarshall(field reflect.Value, value string) error {
	dupField := field
	//fmt.Printf("unmarshal.field.Kind: %s, .Type %s value %s\n",field.Kind(), field.Type(), value)
	unMarshallIt := func(finalField reflect.Value) error {
		if finalField.CanInterface() {
			fieldIface := finalField.Interface()

			fieldTypeUnmarshaller, ok := fieldIface.(TypeUnmarshaller)
			if ok {
				return fieldTypeUnmarshaller.UnmarshalXLS(value)
			}

			// Otherwise try to use TextUnmarshaler
			fieldTextUnmarshaler, ok := fieldIface.(encoding.TextUnmarshaler)
			if ok {
				return fieldTextUnmarshaler.UnmarshalText([]byte(value))
			}
		}

		return NoUnmarshalFuncError{"No known conversion from string to " + field.Type().String() + ", " + field.Type().String() + " does not implement TypeUnmarshaller"}
	}
	for dupField.Kind() == reflect.Interface || dupField.Kind() == reflect.Ptr {
		if dupField.IsNil() {
			dupField = reflect.New(field.Type().Elem())
			field.Set(dupField)
			return unMarshallIt(dupField)
		}
		dupField = dupField.Elem()
	}
	if dupField.CanAddr() {
		return unMarshallIt(dupField.Addr())
	}
	return NoUnmarshalFuncError{"No known conversion from string to " + field.Type().String() + ", " + field.Type().String() + " does not implement TypeUnmarshaller"}
}
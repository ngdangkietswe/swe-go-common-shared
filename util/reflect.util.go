package util

import (
	"errors"
	"fmt"
	"reflect"
)

// GetField is a utility function to get the value of a field in a struct.
func GetField(obj interface{}, field string) (interface{}, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil, errors.New("obj must be a struct or a pointer to a struct")
	}

	fieldVal := val.FieldByName(field)
	if !fieldVal.IsValid() {
		return nil, errors.New(fmt.Sprintf("field %s not found", field))
	}

	return fieldVal.Interface(), nil
}

// SetField is a utility function to set the value of a field in a struct.
func SetField(obj interface{}, fieldName string, value interface{}) error {
	val := reflect.ValueOf(obj)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("obj must be a struct or a pointer to a struct")
	}

	field := val.FieldByName(fieldName)
	if !field.IsValid() {
		return fmt.Errorf("field '%s' does not exist", fieldName)
	}

	if !field.CanSet() {
		return fmt.Errorf("field '%s' cannot be set", fieldName)
	}

	fieldVal := reflect.ValueOf(value)
	if field.Type() != fieldVal.Type() {
		return fmt.Errorf("type mismatch: field '%s' has type %s, value has type %s",
			fieldName, field.Type(), fieldVal.Type())
	}

	field.Set(fieldVal)
	return nil
}

// HasField is a utility function to check if a struct has a field.
func HasField(obj interface{}, field string) bool {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return false
	}

	fieldVal := val.FieldByName(field)
	return fieldVal.IsValid()
}

// IsZero is a utility function to check if a value is the zero value of its type.
func IsZero(obj interface{}) bool {
	val := reflect.ValueOf(obj)
	return reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface())
}

// IsNil is a utility function to check if a value is nil.
func IsNil(obj interface{}) bool {
	val := reflect.ValueOf(obj)

	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice:
		return val.IsNil()
	default:
		return false
	}
}

// IsEmpty is a utility function to check if a value is empty.
func IsEmpty(obj interface{}) bool {
	val := reflect.ValueOf(obj)
	switch val.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice, reflect.String:
		return val.Len() == 0
	default:
		return false
	}
}

// IsNilOrEmpty is a utility function to check if a value is nil or empty.
func IsNilOrEmpty(obj interface{}) bool {
	return IsNil(obj) || IsEmpty(obj)
}

// DeepCopy is a utility function to create a deep copy of a value.
func DeepCopy(obj interface{}) (interface{}, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	dst := reflect.New(val.Type()).Elem()

	if err := DeepCopyValue(val, dst); err != nil {
		return nil, fmt.Errorf("failed to deep copy value: %w", err)
	}

	return dst.Interface(), nil
}

// DeepCopyValue is a utility function to create a deep copy of a value.
func DeepCopyValue(src, dst reflect.Value) error {
	switch src.Kind() {
	case reflect.Ptr:
		if src.IsNil() {
			return nil
		}
		dst.Set(reflect.New(src.Type().Elem()))
		return DeepCopyValue(src.Elem(), dst.Elem())

	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			if err := DeepCopyValue(src.Field(i), dst.Field(i)); err != nil {
				return err
			}
		}

	case reflect.Slice:
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			if err := DeepCopyValue(src.Index(i), dst.Index(i)); err != nil {
				return err
			}
		}

	case reflect.Map:
		dst.Set(reflect.MakeMap(src.Type()))
		for _, key := range src.MapKeys() {
			val := src.MapIndex(key)
			newKey := reflect.New(key.Type()).Elem()
			newVal := reflect.New(val.Type()).Elem()
			if err := DeepCopyValue(key, newKey); err != nil {
				return err
			}
			if err := DeepCopyValue(val, newVal); err != nil {
				return err
			}
			dst.SetMapIndex(newKey, newVal)
		}

	default:
		dst.Set(src)
	}

	return nil
}

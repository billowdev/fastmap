package convert

import (
	"reflect"
	"strings"
)

type FastmapConvertOptions struct {
	MatchCase             bool
	IgnoreUnmatchedFields bool
	CustomMatchers        map[string]string
}

func ConvertStruct(src, dst interface{}, opts *FastmapConvertOptions) error {
	if opts == nil {
		opts = &FastmapConvertOptions{MatchCase: true}
	}
	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(dst)
	if srcVal.Kind() != reflect.Ptr || dstVal.Kind() != reflect.Ptr {
		return FastmapErrNotPointer
	}
	srcVal = srcVal.Elem()
	dstVal = dstVal.Elem()
	if srcVal.Kind() != reflect.Struct || dstVal.Kind() != reflect.Struct {
		return FastmapErrNotStruct
	}
	srcType := srcVal.Type()
	dstType := dstVal.Type()
	for i := 0; i < srcType.NumField(); i++ {
		srcField := srcType.Field(i)
		srcFieldVal := srcVal.Field(i)
		dstFieldName := getMatchingField(srcField.Name, dstType, opts)
		if dstFieldName == "" {
			if !opts.IgnoreUnmatchedFields {
				return ErrFieldNotFound
			}
			continue
		}
		dstFieldVal := dstVal.FieldByName(dstFieldName)
		if !dstFieldVal.IsValid() || !dstFieldVal.CanSet() {
			continue
		}
		if srcFieldVal.Type().AssignableTo(dstFieldVal.Type()) {
			dstFieldVal.Set(srcFieldVal)
		}
	}
	return nil
}

func getMatchingField(srcField string, dstType reflect.Type, opts *FastmapConvertOptions) string {
	if mapped, ok := opts.CustomMatchers[srcField]; ok {
		return mapped
	}

	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)
		if opts.MatchCase {
			if srcField == dstField.Name {
				return dstField.Name
			}
		} else {
			if strings.EqualFold(srcField, dstField.Name) {
				return dstField.Name
			}
		}
	}
	return ""
}

var (
	FastmapErrNotPointer = Error("source and destination must be pointers")
	FastmapErrNotStruct  = Error("source and destination must be structs")
	ErrFieldNotFound     = Error("required field not found in destination")
)

type Error string

func (e Error) Error() string { return string(e) }

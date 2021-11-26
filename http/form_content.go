package http

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
)

const (
	attrOmitEmpty = "omitempty"
	attrRequired  = "required"
	attrDefault   = "default"
	attrEnum      = "enum"

	queryTagName = "query"
	formTagName  = "form"
)

func encodeForm(m map[string]interface{}) string {
	if len(m) == 0 {
		return ""
	}
	query := url.Values{}
	for key, value := range m {
		query.Add(key, fmt.Sprintf("%v", value))
	}
	return query.Encode()
}

func formToMap(v interface{}, tagName string) (map[string]interface{}, error) {
	typ := reflect.TypeOf(v)
	isPtr := false
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		isPtr = true
	}
	switch typ.Kind() {
	case reflect.Map:
		return reflectFormFromMap(v, isPtr)
	case reflect.Struct:
		return reflectFormFromStruct(v, isPtr, tagName)
	default:
		return nil, fmt.Errorf("unsupported %s type: %s", tagName, typ.String())
	}
}

func reflectFormFromMap(query interface{}, isPtr bool) (map[string]interface{}, error) {
	m, ok := query.(map[string]interface{})
	if ok {
		return m, nil
	}
	m = make(map[string]interface{})
	var val reflect.Value
	if isPtr {
		val = reflect.ValueOf(query).Elem()
	} else {
		val = reflect.ValueOf(query)
	}
	for _, key := range val.MapKeys() {
		m[fmt.Sprintf("%v", fieldValue(key, false))] = val.MapIndex(key).Interface()
	}
	return m, nil
}

func reflectFormFromStruct(query interface{}, isPtr bool, tagName string) (map[string]interface{}, error) {
	var val reflect.Value
	if isPtr {
		val = reflect.ValueOf(query).Elem()
	} else {
		val = reflect.ValueOf(query)
	}
	var m = make(map[string]interface{})
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanInterface() {
			continue
		}
		structField := val.Type().Field(i)
		tag := structField.Tag
		label := tag.Get(tagName)
		labelElems := strings.Split(label, ",")
		key := labelElems[0]
		if key == "" {
			key = strings.ToLower(structField.Name)
		}
		tagAttrs := func(attrs []string) map[string]struct{} {
			var set = make(map[string]struct{})
			for j := range attrs {
				set[attrs[j]] = struct{}{}
			}
			return set
		}(labelElems[1:])
		var enum bool
		if len(labelElems) > 1 {
			isZero := field.IsZero()
			if _, omitempty := tagAttrs[attrOmitEmpty]; omitempty {
				if isZero {
					continue
				}
			}
			if _, required := tagAttrs[attrRequired]; required {
				if isZero {
					return nil, fmt.Errorf("required field `%s` is empty", structField.Name)
				}
			}
			if isZero {
				for _, attr := range labelElems[1:] {
					if strings.HasPrefix(attr, attrDefault) {
						defaultVal := strings.TrimPrefix(attr, attrDefault)
						m[key] = defaultVal[1:]
						break
					}
				}
				continue
			}
			_, enum = tagAttrs[attrEnum]
		}
		m[key] = fieldValue(field, enum)
	}
	return m, nil
}

func fieldValue(field reflect.Value, enum bool) interface{} {
	separator := ","
	if enum {
		separator = "|"
	}
	switch field.Type().Kind() {
	case reflect.Array:
		fallthrough
	case reflect.Slice:
		stringBuilder := strings.Builder{}
		for i := 0; i < field.Len(); i++ {
			stringBuilder.WriteString(fmt.Sprintf("%v%s", field.Index(i), separator))
		}
		return stringBuilder.String()[:stringBuilder.Len()-1]
	default:
		return field.Interface()
	}
}

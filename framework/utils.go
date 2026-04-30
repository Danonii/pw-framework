package framework

import (
	"reflect"
)

func convertStructToMap(a any) (map[string]any, error) {
	const tagKey = "structmap"
	var err error = nil
	m := make(map[string]any)
	ptr := reflect.ValueOf(a)
	v := reflect.Indirect(ptr)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		fv := v.Field(i)

		var key string

		if key == "" {
			key = f.Name
		}

		value := fv.Interface()
		if f.Type.Kind() == reflect.Struct {
			value, err = convertStructToMap(value)
		}

		m[key] = value
	}
	return m, err
}

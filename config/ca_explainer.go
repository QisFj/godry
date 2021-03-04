package config

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Explainer func(raw string) (interface{}, error)

func explain(explainer Explainer, raw string) (interface{}, error) {
	if explainer == nil {
		return nil, nil
	}
	return explainer(raw)
}

func NewCommonJSONExplainer(v interface{}) Explainer {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr || t.Kind() == reflect.Interface {
		t = t.Elem()
	}
	return func(raw string) (interface{}, error) {
		vv := reflect.New(t).Interface()
		if err := json.Unmarshal([]byte(raw), vv); err != nil {
			return nil, fmt.Errorf("json unmarshal error: %w", err)
		}
		return vv, nil
	}
}

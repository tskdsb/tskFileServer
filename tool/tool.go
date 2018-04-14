package tool

import (
  "encoding/json"
  "reflect"
)

func Json2object(jsonData []byte, object interface{}) error {
  return json.Unmarshal(jsonData, object)
}

func Object2json(object interface{}) ([]byte, error) {
  return json.MarshalIndent(object, "", "  ")
}

func Field(i interface{}, name string) interface{} {
  return reflect.ValueOf(i).Elem().FieldByName(name).Interface()
}

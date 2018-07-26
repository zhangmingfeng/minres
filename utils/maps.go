package utils

import (
	"encoding/json"
	"errors"
	"github.com/devfeel/mapper"
)

func Map2Struct(fromMap interface{}, toStruct interface{}) error {
	switch val := fromMap.(type) {
	case string:
		tmpMap := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(val), &tmpMap)
		if err != nil {
			return err
		}
		mapper.MapperMap(tmpMap, toStruct)
	case []byte:
		tmpMap := make(map[string]interface{}, 0)
		err := json.Unmarshal(val, &tmpMap)
		if err != nil {
			return err
		}
		mapper.MapperMap(tmpMap, toStruct)
	case map[string]interface{}:
		mapper.MapperMap(val, toStruct)
	default:
		return errors.New("类型不正确，无法转换")
	}
	return nil
}

func Json2Struct(fromJson interface{}, toStruct interface{}) error {
	switch val := fromJson.(type) {
	case string:
		return json.Unmarshal([]byte(val), toStruct)
	case []byte:
		return json.Unmarshal(val, toStruct)
	default:
		return errors.New("类型不正确，无法转换")
	}
	return nil
}

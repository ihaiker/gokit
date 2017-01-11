package json

import (
	"github.com/ihaiker/gokit/config"
	"errors"
	"reflect"

	"encoding/json"
)

func mergerMethod(src,dst []byte,path string) ([]byte,error){
	var srcJson, dstJson map[string]interface{}

	if err := json.Unmarshal(src, &srcJson); err != nil {
		return nil,err
	}
	if err := json.Unmarshal(dst, &dstJson); err != nil {
		return nil, err
	}
	outJson, err := merge(srcJson, dstJson, path)
	if err != nil {
		return nil,err
	}
	out, err := json.Marshal(outJson)
	return out,err
}

func merge(srcJson, dstJson map[string]interface{}, path string) (map[string]interface{}, error) {

	outJson := srcJson
	for key, value := range dstJson {
		current_path := path + "." + key
		srcValue := srcJson[key]
		typeof := reflect.TypeOf(value)

		switch typeof.Kind() {
		default:
			return nil, errors.New("unsupport " + typeof.String() + " in (" + current_path + ")")
		case reflect.Bool,
            reflect.Int,reflect.Int8,reflect.Int16,reflect.Int32,reflect.Int64,
            reflect.Uint,reflect.Uint8,reflect.Uint16,reflect.Uint32,reflect.Uint64,
            reflect.Float32,reflect.Float64,
            reflect.String:
            outJson[key] = value
		case reflect.Array        : fallthrough
		case reflect.Slice        : outJson[key] = value
		case reflect.Map          :
			if srcValue == nil {
				srcValue = make(map[string]interface{})
			}
			outValue, err := merge(srcValue.(map[string]interface{}), value.(map[string]interface{}), current_path)
			if err != nil {
				return nil, err
			}
			outJson[key] = outValue
		}
	}
	return outJson, nil
}

func Config(items... interface{}) (*config.Config, error) {
	cfg := config.NewConfig([]byte("{}"),json.Unmarshal, Merger())
	var err error
	if len(items) != 0 {
		err = cfg.Load(items...)
	}
	return cfg, err
}


func Merger() *config.Merger {
	return config.NewMerger(mergerMethod)
}
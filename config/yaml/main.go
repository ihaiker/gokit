package yaml

import (
	"gopkg.in/yaml.v2"
	"github.com/ihaiker/gokit/config"
	"reflect"
	"errors"
    "github.com/ihaiker/gokit/convert"
)
func merge(srcYAML, dstYAML map[interface{}]interface{}, path string) (map[interface{}]interface{}, error) {
	outYAML := srcYAML
	for key, value := range dstYAML {
		srcValue := srcYAML[key]
		current_path := path + "." + convertKit.SafeString(key)
		typeof := reflect.TypeOf(value)

		switch typeof.Kind() {
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.String       :
			outYAML[key] = value
		case reflect.Array        : fallthrough
		case reflect.Slice        : outYAML[key] = value
		case reflect.Map          :
			if srcValue == nil {
				srcValue = make(map[interface{}]interface{})
            }
            
			outValue, err := merge(srcValue.(map[interface{}]interface{}), value.(map[interface{}]interface{}), current_path)
			if err != nil {
				return nil, err
			}
			outYAML[key] = outValue
		default:
			return nil, errors.New("unsupport " + typeof.String() + " in (" + current_path + ")")
		}
	}
	return outYAML, nil
}

func mergerMethod(src, dst []byte, path string) ([]byte, error) {
	var srcYaml, dstYaml map[interface{}]interface{}

	if err := yaml.Unmarshal(src, &srcYaml); err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(dst, &dstYaml); err != nil {
		return nil, err
	}
	outYaml, err := merge(srcYaml, dstYaml, path)
	if err != nil {
		return nil, err
	}
	out, err := yaml.Marshal(outYaml)
	return out, err
}

func Config(items... interface{}) (*config.Config, error) {
	cfg := config.NewConfig([]byte("{}"), yaml.Unmarshal, Merger())
	var err error
	if len(items) != 0 {
		err = cfg.Load(items...)
	}
	return cfg, err
}

func Merger() *config.Merger {
	return config.NewMerger(mergerMethod)
}
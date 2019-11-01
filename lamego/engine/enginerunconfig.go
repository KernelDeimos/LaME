package engine

import (
	"encoding/json"
	"errors"
)

type EngineRunConfig struct {
	Name                     string                            `yaml:"name"`
	TargetLanguage           string                            `yaml:"target"`
	ModelSourceDirectory     string                            `yaml:"source"`
	GeneratorOutputDirectory string                            `yaml:"output"`
	PluginConfig             map[string]map[string]interface{} `yaml:"config"`
}

func (conf EngineRunConfig) GetConfig(name string) string {
	// Unfortunately the yaml unmarshaller makes a map[interface{}]interface{},
	// which is useless here. This API is supposed to return a JSON string that
	// that the plugin can choose how to process.

	// Recursively replace every map in the PluginConfig map
	thisConfig := conf.PluginConfig[name]
	var recur func(v interface{}) (interface{}, error)
	recur = func(v interface{}) (interface{}, error) {
		uselessMap, isMap := v.(map[interface{}]interface{})
		if !isMap {
			return v, nil
		}
		notUselessMap := map[string]interface{}{}
		for k, v := range uselessMap {
			kStr, ok := k.(string)
			if !ok {
				return nil, errors.New("unexpected non-string key in config")
			}
			v, err := recur(v)
			if err != nil {
				return nil, err
			}
			notUselessMap[kStr] = v
		}
		return notUselessMap, nil
	}
	for k, v := range thisConfig {
		var err error
		v, err = recur(v)
		if err != nil {
			panic(err)
		}
		thisConfig[k] = v
	}

	// Call json.Marshal, now that its input won't produce an error
	b, err := json.Marshal(conf.PluginConfig[name])
	if err != nil {
		panic(err)
	}
	return string(b)
}

package config

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func RemarshalYAML(from interface{}, to interface{}) error {
	bytes, err := yaml.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	return UnmarshalYAML(bytes, to)
}

func UnmarshalYAML(from []byte, to interface{}) error {
	err := yaml.UnmarshalStrict(from, to)
	if err != nil {
		return errors.Wrap(err, "unmarshalling")
	}
	//
	// defaultable, ok := to.(Defaultable)
	// if ok {
	// 	defaultable.ApplyDefaults()
	// }

	return nil
}

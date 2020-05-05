package configutil

import (
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

func RemarshalYAML(from interface{}, to interface{}) error {
	bytes, err := yaml.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "marshaling")
	}

	return UnmarshalYAMLStrict(bytes, to)
}

func UnmarshalYAMLStrict(from []byte, to interface{}) error {
	err := yaml.UnmarshalStrict(from, to)
	if err != nil {
		return errors.Wrap(err, "unmarshalling")
	}

	return nil
}

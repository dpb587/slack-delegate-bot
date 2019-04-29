package args

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type LogLevel logrus.Level

func (ll *LogLevel) UnmarshalFlag(data string) error {
	parsed, err := logrus.ParseLevel(data)
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}

	*ll = LogLevel(parsed)

	return nil
}

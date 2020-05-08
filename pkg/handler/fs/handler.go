package fs

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/dpb587/slack-delegate-bot/pkg/handler"
	"github.com/dpb587/slack-delegate-bot/pkg/handler/yaml"
	"github.com/pkg/errors"
)

func BuildHandler(parser *yaml.Parser, paths ...string) (handler.Handler, error) {
	var handlers []handler.Handler

	paths, err := squashPaths(paths)
	if err != nil {
		return nil, errors.Wrap(err, "squashing paths")
	}

	for _, path := range paths {
		pathBytes, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, errors.Wrapf(err, "reading %s", path)
		}

		h, err := parser.Parse(pathBytes)
		if err != nil {
			return nil, errors.Wrapf(err, "parsing %s", path)
		}

		handlers = append(handlers, h)
	}

	if len(handlers) == 1 {
		return handlers[0], nil
	}

	return handler.NewCoalesceHandler(handlers...), nil
}

func squashPaths(paths []string) ([]string, error) {
	var squashed []string

	for _, path := range paths {
		globbed, err := filepath.Glob(path)
		if err != nil {
			return nil, errors.Wrapf(err, "globbing %s", path)
		}

		sort.Strings(globbed)

		squashed = append(squashed, globbed...)
	}

	return squashed, nil
}

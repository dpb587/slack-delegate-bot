package fs

import (
	"io/ioutil"
	"path/filepath"
	"sort"

	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/coalesce"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/provider/yaml"
	"github.com/pkg/errors"
)

func BuildDelegator(parser *yaml.Parser, paths ...string) (delegate.Delegator, error) {
	var delegators []delegate.Delegator

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

		delegators = append(delegators, h)
	}

	if len(delegators) == 1 {
		return delegators[0], nil
	}

	return coalesce.Delegator{Delegators: delegators}, nil
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

package opts

import (
	"fmt"
	"os"
	"strings"

	"github.com/dpb587/slack-delegate-bot/cmd/delegatebot/args"
	conditionsfactory "github.com/dpb587/slack-delegate-bot/pkg/condition/conditions/defaultfactory"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/coalesce"
	interruptsfactory "github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/defaultfactory"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/provider/db"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/provider/fs"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/provider/yaml"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	// include potential database adapters
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Root struct {
	Configs   []string `long:"config" description:"Path to configuration files" env:"CONFIG"`
	delegator delegate.Delegator

	LogLevel args.LogLevel `long:"log-level" description:"Show additional levels of log messages" env:"LOG_LEVEL" default:"INFO"`
	logger   logrus.FieldLogger

	parser *yaml.Parser
}

func (r *Root) GetLogger() logrus.FieldLogger {
	if r.logger == nil {
		var logger = logrus.New()
		logger.Out = os.Stderr
		logger.Formatter = &logrus.JSONFormatter{}
		logger.Level = logrus.Level(r.LogLevel)

		r.logger = logger
	}

	return r.logger
}

func (r *Root) GetParser() *yaml.Parser {
	if r.parser == nil {
		conditions := conditionsfactory.New()
		interrupts := interruptsfactory.New(conditions)

		r.parser = yaml.NewParser(interrupts, conditions)
	}

	return r.parser
}

func (r *Root) GetDelegator() (delegate.Delegator, error) {
	var delegators []delegate.Delegator
	var filePaths []string

	for _, uri := range r.Configs {
		uriSplit := strings.SplitN(uri, "://", 2)

		if len(uriSplit) == 1 {
			filePaths = append(filePaths, uriSplit[0])

			continue
		}

		switch uriSplit[0] {
		case "mysql", "sqlite3":
			dbh, err := db.OpenDB(uriSplit[0], uriSplit[1])
			if err != nil {
				return nil, errors.Wrapf(err, "opening db")
			}

			delegators = append(delegators, db.NewDelegator(dbh, r.GetParser()))
		default:
			return nil, fmt.Errorf("unsupported handler uri: %s", uri)
		}
	}

	if len(filePaths) > 0 {
		// collected for later to be able to squash paths
		h, err := fs.BuildDelegator(r.GetParser(), r.Configs...)
		if err != nil {
			return nil, errors.Wrap(err, "loading configs")
		}

		delegators = append(delegators, h)
	}

	if len(delegators) == 1 {
		return delegators[0], nil
	}

	res := coalesce.Delegator{
		Delegators: delegators,
	}

	return res, nil
}

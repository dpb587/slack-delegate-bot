package factory

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dpb587/go-pairist/api"
	"github.com/dpb587/slack-delegate-bot/pkg/config"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates"
	"github.com/dpb587/slack-delegate-bot/pkg/delegate/delegates/pairist"
	"github.com/pkg/errors"
)

type factory struct{}

type Options struct {
	Team     string `yaml:"team"`
	Password string `yaml:"password"`

	Role  string `yaml:"role"`
	Track string `yaml:"track"`
}

func New() delegates.Factory {
	return &factory{}
}

func (f factory) Create(name string, options interface{}) (delegate.Delegator, error) {
	if name != "pairist" {
		return nil, fmt.Errorf("unsupported delegate: %s", name)
	}

	parsed := Options{}

	err := config.RemarshalYAML(options, &parsed)
	if err != nil {
		return nil, errors.Wrap(err, "remarshalling")
	}

	if parsed.Role != "" && parsed.Track != "" {
		return nil, errors.New("only one of the following may be set: role, track")
	}

	var clientAuth *api.Auth

	if parsed.Password != "" {
		if strings.HasPrefix(parsed.Password, "$") && len(parsed.Password) > 1 {
			parsed.Password = os.Getenv(parsed.Password[1:])
		}

		clientAuth = &api.Auth{
			APIKey:   os.Getenv("PAIRIST_API_KEY"),
			Team:     parsed.Team,
			Password: parsed.Password,
		}
	}

	return &pairist.Delegator{
		Client: api.NewClient(http.DefaultClient, api.DefaultDatabaseURL, clientAuth),
		Team:   parsed.Team,
		Role:   parsed.Role,
		Track:  parsed.Track,
	}, nil
}

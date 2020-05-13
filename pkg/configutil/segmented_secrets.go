package configutil

import (
	"crypto/sha1"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var segmentedSecretPaths = []string{
	"/pairist/password",
	"/pagerduty/api_key",
}

func DesanitizeConfig(in string, secrets map[string]interface{}) (string, map[string]interface{}, error) {
	var intermediate interface{}

	err := yaml.Unmarshal([]byte(in), &intermediate)
	if err != nil {
		return "", nil, errors.Wrap(err, "unmarshalling")
	}

	intermediate, usedSecrets, err := desanitizeConfig(intermediate, secrets, nil)
	if err != nil {
		return "", nil, err
	}

	out, err := yaml.Marshal(intermediate)
	if err != nil {
		return "", nil, errors.Wrap(err, "marshalling")
	}

	return string(out), usedSecrets, nil
}

func desanitizeConfig(inT interface{}, secrets map[string]interface{}, chain []string) (interface{}, map[string]interface{}, error) {
	var err error
	usedSecrets := map[string]interface{}{}
	var localUsedSecrets map[string]interface{}

	switch in := inT.(type) {
	case []interface{}:
		for key, val := range in {
			in[key], localUsedSecrets, err = desanitizeConfig(val, secrets, append(chain, fmt.Sprintf("%v", key)))
			if err != nil {
				return nil, nil, errors.Wrapf(err, "sanitizing node %s/%d", strings.Join(chain, "/"), key)
			}

			for k, v := range localUsedSecrets {
				usedSecrets[k] = v
			}
		}
	case map[interface{}]interface{}:
		for key, val := range in {
			in[key], localUsedSecrets, err = desanitizeConfig(val, secrets, append(chain, fmt.Sprintf("%v", key)))
			if err != nil {
				return nil, nil, errors.Wrapf(err, "sanitizing node %s/%v", strings.Join(chain, "/"), key)
			}

			for k, v := range localUsedSecrets {
				usedSecrets[k] = v
			}
		}
	case string:
		if !strings.HasPrefix(in, "@secret:") {
			break
		}

		secretID := strings.TrimPrefix(in, "@secret:")

		secret, found := secrets[secretID]
		if !found {
			// fail (un?)safe
			break
		}

		inT = secret
		usedSecrets[secretID] = secret
	}

	return inT, usedSecrets, nil
}

func SanitizeConfig(in string) (string, map[string]interface{}, error) {
	var intermediate interface{}

	err := yaml.Unmarshal([]byte(in), &intermediate)
	if err != nil {
		return "", nil, errors.Wrap(err, "unmarshalling")
	}

	intermediate, secrets, err := sanitizeConfig(intermediate, map[string]interface{}{}, nil)
	if err != nil {
		return "", nil, err
	}

	out, err := yaml.Marshal(intermediate)
	if err != nil {
		return "", nil, errors.Wrap(err, "marshalling")
	}

	return string(out), secrets, nil
}

func sanitizeConfig(inT interface{}, secrets map[string]interface{}, chain []string) (interface{}, map[string]interface{}, error) {
	var err error

	switch in := inT.(type) {
	case []interface{}:
		for key, val := range in {
			in[key], secrets, err = sanitizeConfig(val, secrets, append(chain, fmt.Sprintf("%v", key)))
			if err != nil {
				return nil, nil, errors.Wrapf(err, "sanitizing node %s/%d", strings.Join(chain, "/"), key)
			}
		}
	case map[interface{}]interface{}:
		for key, val := range in {
			in[key], secrets, err = sanitizeConfig(val, secrets, append(chain, fmt.Sprintf("%v", key)))
			if err != nil {
				return nil, nil, errors.Wrapf(err, "sanitizing node %s/%v", strings.Join(chain, "/"), key)
			}
		}
	default:
		secretID, secret := checkSecretPath(chain)
		if !secret {
			break
		}

		secrets[secretID] = inT
		inT = fmt.Sprintf("@secret:%s", secretID)
	}

	return inT, secrets, nil
}

func checkSecretPath(chain []string) (string, bool) {
	chainJoin := strings.Join(chain, "/")

	for _, pathset := range segmentedSecretPaths {
		if strings.HasSuffix(chainJoin, pathset) {
			h := sha1.New()
			h.Write([]byte(chainJoin))
			h.Write([]byte(fmt.Sprintf("%d", time.Now().Unix())))

			return fmt.Sprintf("%x", h.Sum(nil))[0:12], true
		}
	}

	return "", false
}

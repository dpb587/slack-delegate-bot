package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	identitytoolkit "google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
)

type Auth struct {
	APIKey   string
	Team     string
	Password string

	idToken      string
	refreshAfter time.Time
}

func (a *Auth) IDToken() (string, error) {
	if a.idToken == "" || a.refreshAfter.Before(time.Now()) {
		ctx := context.Background()

		client, err := identitytoolkit.NewService(ctx, option.WithAPIKey(a.APIKey))
		if err != nil {
			return "", err
		}

		resp, err := client.Relyingparty.VerifyPassword(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
			Email:             fmt.Sprintf("%s@pair.ist", a.Team),
			Password:          a.Password,
			ReturnSecureToken: true,
		}).Context(ctx).Do()
		if err != nil {
			return "", err
		}

		a.idToken = resp.IdToken

		idTokenSplit := strings.SplitN(a.idToken, ".", 3)
		if len(idTokenSplit) != 3 {
			return "", fmt.Errorf("parsing auth token: expected 3 parts, found %d", len(idTokenSplit))
		}

		claimsData, err := base64.StdEncoding.DecodeString(idTokenSplit[1])
		if err != nil {
			return "", errors.Wrap(err, "decoding auth token")
		}

		var claims struct {
			Exp int64 `json:"exp"`
		}

		err = json.Unmarshal(claimsData, &claims)
		if err != nil {
			return "", errors.Wrap(err, "unmarshalling auth token claims")
		}

		a.refreshAfter = time.Unix(claims.Exp, 0).Add(-5 * time.Minute)
	}

	return a.idToken, nil
}

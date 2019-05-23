package api

import (
	"context"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
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

		var claims jwt.StandardClaims
		parser := &jwt.Parser{}

		_, _, err = parser.ParseUnverified(a.idToken, &claims)
		if err != nil {
			return "", errors.Wrap(err, "parsing auth token")
		}

		a.refreshAfter = time.Unix(claims.ExpiresAt, 0).Add(-5 * time.Minute)
	}

	return a.idToken, nil
}

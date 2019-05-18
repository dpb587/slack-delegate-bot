package api

import (
	"context"
	"fmt"

	identitytoolkit "google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
)

type Auth struct {
	APIKey   string
	Team     string
	Password string

	idToken string
}

func (a *Auth) IDToken() (string, error) {
	if a.idToken == "" {
		ctx := context.Background()

		client, err := identitytoolkit.NewService(ctx, option.WithAPIKey(a.APIKey))
		if err != nil {
			return "", err
		}

		// seems like these tokens do not expire, so no need to refresh?
		resp, err := client.Relyingparty.VerifyPassword(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
			Email:             fmt.Sprintf("%s@pair.ist", a.Team),
			Password:          a.Password,
			ReturnSecureToken: true,
		}).Context(ctx).Do()
		if err != nil {
			return "", err
		}

		a.idToken = resp.IdToken
	}

	return a.idToken, nil
}

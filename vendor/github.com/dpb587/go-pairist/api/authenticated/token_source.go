package authenticated

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	identitytoolkit "google.golang.org/api/identitytoolkit/v3"
	"google.golang.org/api/option"
)

type TokenSource struct {
	ctx                    context.Context
	client                 *identitytoolkit.Service
	latest                 *identitytoolkit.VerifyPasswordResponse
	teamName, teamPassword string
}

func CreateTokenSource(ctx context.Context, apiKey, teamName, teamPassword string) (*TokenSource, error) {
	client, err := identitytoolkit.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		panic(err)
	}

	return &TokenSource{
		ctx:          ctx,
		client:       client,
		teamName:     teamName,
		teamPassword: teamPassword,
	}, nil
}

func (ts *TokenSource) Token() (*oauth2.Token, error) {
	if ts.latest == nil {
		// seems like these tokens do not expire, so no need to refresh?
		resp, err := ts.client.Relyingparty.VerifyPassword(&identitytoolkit.IdentitytoolkitRelyingpartyVerifyPasswordRequest{
			Email:    fmt.Sprintf("%s@pair.ist", ts.teamName),
			Password: ts.teamPassword,
		}).Context(ts.ctx).Do()
		if err != nil {
			return nil, err
		}

		ts.latest = resp
	}

	return &oauth2.Token{
		TokenType:    "Bearer",
		AccessToken:  ts.latest.OauthAccessToken,
		RefreshToken: ts.latest.RefreshToken,
	}, nil
}

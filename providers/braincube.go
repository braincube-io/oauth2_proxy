package providers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/bitly/oauth2_proxy/api"
)

type BraincubeProvider struct {
	*ProviderData
}

func NewBraincubeProvider(p *ProviderData) *BraincubeProvider {
	p.ProviderName = "braincube"
	fmt.Println("Braincube provider")
	if p.LoginURL == nil || p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{
			Scheme: "https",
			Host:   "mybraincube.com",
			Path:   "/sso-server/ws/oauth2/authorize.jsp",
		}
	}
	if p.RedeemURL == nil || p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{
			Scheme: "https",
			Host:   "mybraincube.com",
			Path:   "/sso-server/ws/oauth2/token",
		}
	}
	if p.ProfileURL == nil || p.ProfileURL.String() == "" {
		p.ProfileURL = &url.URL{
			Scheme: "https",
			Host:   "mybraincube.com",
			Path:   "/sso-server/ws/oauth2/me",
		}
	}
	if p.Scope == "" {
		p.Scope = "BASE"
	}
	return &BraincubeProvider{ProviderData: p}
}

func getBraincubeHeader(AccessToken string) http.Header {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", AccessToken))
	return header
}

func (p *BraincubeProvider) GetEmailAddress(s *SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ProfileURL.String(), nil)
	if err != nil {
		return "", err
	}
	req.Header = getBraincubeHeader(s.AccessToken)

	json, err := api.Request(req)
	if err != nil {
		return "", err
	}

	email, err := json.Get("userEmail").String()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return "", err
	}
	return email, nil
}

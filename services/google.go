package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"shirini-backend/config"
)

type GoogleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
}

func ExchangeCode(code string, cfg *config.Config) (*GoogleUser, error) {

	data := url.Values{}

	data.Set("code", code)
	data.Set("client_id", cfg.Google.ClientID)
	data.Set("client_secret", cfg.Google.ClientSecret)
	data.Set("redirect_uri", cfg.Google.RedirectURI)
	data.Set("grant_type", "authorization_code")


	resp, err := http.PostForm(
		"https://oauth2.googleapis.com/token",
		data,
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()


	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google token exchange failed: %s", resp.Status)
	}


	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
	}

	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)

	if err != nil {
		return nil, err
	}


	userResp, err := http.Get(
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token=" +
			tokenResponse.AccessToken,
	)

	if err != nil {
		return nil, err
	}

	defer userResp.Body.Close()


	var user GoogleUser

	err = json.NewDecoder(userResp.Body).Decode(&user)

	if err != nil {
		return nil, err
	}


	return &user, nil
}
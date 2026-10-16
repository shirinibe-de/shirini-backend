package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const stateCookieName = "oauth_state"
const stateMaxAge = 10 * time.Minute

type oauthState struct {
	Nonce       string `json:"n"`
	RedirectURI string `json:"r"`
}

func signState(s *oauthState, secret string) (string, error) {
	payload, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	sig := mac.Sum(nil)

	return base64.RawURLEncoding.EncodeToString(payload) + "." +
		base64.RawURLEncoding.EncodeToString(sig), nil
}

func verifyState(encoded string, secret string) (*oauthState, error) {
	parts := strings.Split(encoded, ".")
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed state")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, fmt.Errorf("malformed state payload")
	}

	sig, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("malformed state signature")
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)

	if !hmac.Equal(sig, mac.Sum(nil)) {
		return nil, fmt.Errorf("invalid state signature")
	}

	var s oauthState
	if err := json.Unmarshal(payload, &s); err != nil {
		return nil, fmt.Errorf("malformed state payload")
	}

	return &s, nil
}

func setStateCookie(c *gin.Context, nonce string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     stateCookieName,
		Value:    nonce,
		Path:     "/api/v1/auth/google/callback",
		MaxAge:   int(stateMaxAge.Seconds()),
		HttpOnly: true,
		Secure:   false, // localhost dev; set true behind TLS
		SameSite: http.SameSiteLaxMode,
	})
}

func clearStateCookie(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     stateCookieName,
		Value:    "",
		Path:     "/api/v1/auth/google/callback",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

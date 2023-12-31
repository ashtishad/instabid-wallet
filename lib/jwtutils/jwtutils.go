package jwtutils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var (
	HMACSecret        = []byte(os.Getenv("HMACSecret"))
	ErrUnauthorized   = errors.New("unauthorized")
	ErrEmptyToken     = errors.New("token cannot be empty")
	ErrEmptyEnvVars   = errors.New("API_HOST, AUTH_API_PORT, or API_SCHEME cannot be empty")
	ErrInvalidURL     = errors.New("could not build a valid URL")
	ErrEmptyRouteName = errors.New("routeName cannot be empty")
)

// ParseAndValidateToken parses a JWT token string and validates its signature.
// The function returns the parsed token if it's valid, and an error otherwise.
// nolint:wrapcheck
func ParseAndValidateToken(tokenStr string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}

		return HMACSecret, nil
	})

	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, jwt.ErrTokenMalformed
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, jwt.ErrTokenSignatureInvalid
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, err
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, err
		default:
			return nil, err
		}
	}

	return token, nil
}

// VerifyTokenWithAuthAPI sends a request to the Auth APIs verify endpoint to validate a JWT token.
// The function takes a JWT token string as input and returns its claims if the token is valid.
// If the token is invalid or an error occurs (e.g., failed to build the URL, HTTP request failure, etc.),
// the function will return an error.
func VerifyTokenWithAuthAPI(tokenStr string, routeName string, pathUserID string) (jwt.MapClaims, error) {
	verifyURL, err := buildVerifyURL(tokenStr, routeName, pathUserID)
	if err != nil {
		return nil, fmt.Errorf("error building URL: %w", err)
	}

	var resp *http.Response

	// nolint:gosec // using env variables to build secure verify url
	resp, err = http.Get(verifyURL)
	if err != nil {
		return nil, fmt.Errorf("unable to get response from verify url:%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ErrUnauthorized
	}

	var response struct {
		Claims jwt.MapClaims `json:"claims"`
	}

	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("unable to decode json:%w", err)
	}

	return response.Claims, nil
}

// buildVerifyURL constructs a URL for the verify endpoint of the Auth API.
// The function takes a JWT token string as an argument, and returns a formatted URL.
// It uses environment variables "API_HOST" and "AUTH_API_PORT" to determine the APIs location.
// It returns an error if it fails to construct a valid URL.
// e.g: http://127.0.0.1:8001/verify?token=JWT_TOKEN&routeName=?&userId=?
func buildVerifyURL(tokenStr string, routeName string, pathUserID string) (string, error) {
	if tokenStr == "" {
		return "", ErrEmptyToken
	}

	if routeName == "" {
		return "", ErrEmptyRouteName
	}

	apiHost := os.Getenv("API_HOST")
	authAPIPort := os.Getenv("AUTH_API_PORT")
	apiScheme := os.Getenv("API_SCHEME")

	if apiHost == "" || authAPIPort == "" || apiScheme == "" {
		return "", ErrEmptyEnvVars
	}

	u := &url.URL{
		Scheme: apiScheme,
		Host:   fmt.Sprintf("%s:%s", apiHost, authAPIPort),
		Path:   "/verify",
	}

	q := u.Query()
	q.Add("token", tokenStr)
	q.Add("routeName", routeName)
	q.Add("userId", pathUserID)
	u.RawQuery = q.Encode()

	_, err := url.ParseRequestURI(u.String())
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	return u.String(), nil
}

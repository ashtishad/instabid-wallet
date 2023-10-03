package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"

	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthHeader = "Authorization"
	Bearer     = "Bearer"
)

var (
	ErrAuthHeaderNotFound  = errors.New("authorization header not found")
	ErrBearerTokenNotFound = errors.New("bearer token not found in auth header")
	ErrTypeAssertionFailed = errors.New("type assertion failed for one or more token claims")
	ErrUnauthorized        = errors.New("unauthorized")
)

// validateJWTMiddleware is a Gin middleware function that authorises incoming HTTP requests
// by validating JWT tokens found in the "Authorization" header.
// If the token is valid, it extracts the claims and sets them in the Gin context.
// It builds verify url and sends a get request to auth-api using verifyTokenWithAuthAPI,
// Otherwise, it responds with a 401 Unauthorized status and aborts the request.
func validateJWTMiddleware(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := extractToken(c)
		if err != nil {
			l.Error("unable to extract token", "err", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()

			return
		}

		claims, err := verifyTokenWithAuthAPI(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()

			return
		}

		var user *domain.AuthorizedUser
		user, err = getAuthorizedUserFromClaims(claims)

		if err != nil {
			l.Error("unable to get authorized user from claims", "err", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()

			return
		}

		c.Set("authorizedUserRequest", user)
		c.Next()

		c.Next()
	}
}

// extractToken retrieves the JWT token from the "Authorization" header.
// Returns an error if the header is missing or improperly formatted.
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader(AuthHeader)
	if authHeader == "" {
		return "", ErrAuthHeaderNotFound
	}

	var tokenStr string
	_, err := fmt.Sscanf(authHeader, Bearer+" %s", &tokenStr)

	if err != nil {
		return "", ErrBearerTokenNotFound
	}

	return tokenStr, nil
}

// getAuthorizedUserFromClaims extracts a User object from a set of JWT claims.
// It returns the User object if all required claims are present and valid,
// otherwise returns an error.
func getAuthorizedUserFromClaims(claims jwt.MapClaims) (*domain.AuthorizedUser, error) {
	userID, ok1 := claims["UserID"].(string)
	userName, ok2 := claims["Username"].(string)
	email, ok3 := claims["Email"].(string)
	role, ok4 := claims["Role"].(string)
	status, ok5 := claims["Status"].(string)

	if ok1 && ok2 && ok3 && ok4 && ok5 {
		user := &domain.AuthorizedUser{
			UserID:   userID,
			UserName: userName,
			Email:    email,
			Role:     role,
			Status:   status,
		}

		return user, nil
	}

	return nil, ErrTypeAssertionFailed
}

// buildVerifyURL constructs a URL for the verify endpoint of the Auth API.
// The function takes a JWT token string as an argument, and returns a formatted URL.
// It uses environment variables "API_HOST" and "AUTH_API_PORT" to determine the APIs location.
// It returns an error if it fails to construct a valid URL.
// e.g: http://127.0.0.1:8001/verify?token=JWT_TOKEN
func buildVerifyURL(tokenStr string) (string, error) {
	apiHost := os.Getenv("API_HOST")
	authAPIPort := os.Getenv("AUTH_API_PORT")

	rawURL := fmt.Sprintf("http://%s:%s/verify?token=%s", apiHost, authAPIPort, url.QueryEscape(tokenStr))

	_, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return "", fmt.Errorf("could not build a valid URL: %w", err)
	}

	return rawURL, nil
}

// verifyTokenWithAuthAPI sends a request to the Auth APIs verify endpoint to validate a JWT token.
// The function takes a JWT token string as input and returns its claims if the token is valid.
// If the token is invalid or an error occurs (e.g., failed to build the URL, HTTP request failure, etc.),
// the function will return an error.
func verifyTokenWithAuthAPI(tokenStr string) (jwt.MapClaims, error) {
	verifyURL, err := buildVerifyURL(tokenStr)
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

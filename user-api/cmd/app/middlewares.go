package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/ashtishad/instabid-wallet/user-api/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AuthHeader = "Authorization"
	Bearer     = "Bearer"
)

var HMACSecret = []byte(os.Getenv("HMACSecret"))

var (
	ErrAuthHeaderNotFound    = errors.New("authorization header not found")
	ErrBearerTokenNotFound   = errors.New("bearer token not found in auth header")
	ErrTokenValidationFailed = errors.New("token validation failed")
)

// validateJWT is a middleware function that validates the JWT token in the Authorization header.
// If the token is valid, it sets the user ID in the Gin context for further use.
// If the token is invalid or missing, it returns a 401 Unauthorized error.
func validateJWT(l *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := extractToken(c)
		if err != nil {
			l.Error("unable to extract token", "err", err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()

			return
		}

		var token *jwt.Token
		token, err = parseAndValidateToken(tokenStr)

		if err != nil {
			l.Error("token validation failed:", "err", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()

			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			l.Error(ErrTokenValidationFailed.Error())
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

// parseAndValidateToken parses a JWT token string and validates its signature.
// The function returns the parsed token if it's valid, and an error otherwise.
// nolint:wrapcheck
func parseAndValidateToken(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return HMACSecret, nil
	})
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

	return nil, fmt.Errorf("type assertion failed for one or more token claims")
}

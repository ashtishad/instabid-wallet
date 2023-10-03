package app

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/ashtishad/instabid-wallet/lib/jwtutils"
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

		claims, err := jwtutils.VerifyTokenWithAuthAPI(tokenStr)
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

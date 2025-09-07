package middleware

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"github.com/SteelyBretty/consultant-time-tracker/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func BasicAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := strings.TrimSpace(c.GetHeader("Authorization"))
		if auth == "" {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization required",
			})
			return
		}

		if !strings.HasPrefix(auth, "Basic ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			return
		}

		payload, err := base64.StdEncoding.DecodeString(strings.TrimSpace(auth[6:]))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization encoding",
			})
			return
		}

		pair := strings.SplitN(string(payload), ":", 2)
		if len(pair) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization format",
			})
			return
		}

		username := pair[0]
		password := pair[1]

		authService := services.NewAuthService()
		user, err := authService.ValidateCredentials(username, password)
		if err != nil {
			c.Header("WWW-Authenticate", `Basic realm="Restricted"`)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials",
			})
			return
		}

		// store UUID directly
		c.Set("user_id", user.ID)
		c.Set("user", user)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, errors.New("user not present in context")
	}

	switch v := userIDVal.(type) {
	case uuid.UUID:
		return v, nil
	case string:
		u, err := uuid.Parse(v)
		if err != nil {
			return uuid.Nil, err
		}
		return u, nil
	default:
		return uuid.Nil, errors.New("invalid user_id type in context")
	}
}

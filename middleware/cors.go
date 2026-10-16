package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"shirini-backend/config"
)

// CORSMiddleware allows the configured frontend origins to call the API
// from the browser. Origins come from the CORS_ORIGINS env var,
// comma-separated. Defaults to local development origins.
func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	origins := cfg.CORS.Origins
	if len(origins) == 0 {
		origins = []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		}
	}

	allowed := make(map[string]struct{}, len(origins))
	for _, o := range origins {
		allowed[o] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if _, ok := allowed[origin]; ok {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

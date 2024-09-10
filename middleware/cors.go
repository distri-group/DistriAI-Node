package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors is a Gin middleware function that handles Cross-Origin Resource Sharing (CORS) headers.
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// Receive origin sent by the client
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			// Methods for all cross-domain requests supported by the server
			c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
			// c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			// Allow cross-domain Settings to return other subsegments, and you can customize fields
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, token, session, X-Requested-With")
			// A header that allows the browser (client) to parse
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			// Setting the cache time
			c.Header("Access-Control-Max-Age", "172800")
			// Allows clients to pass validation information such as cookies
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken, X-CSRF-Token, Authorization, token, X-Requested-With") // custom Header
			c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

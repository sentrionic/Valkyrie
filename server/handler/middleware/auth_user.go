package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
)

// AuthUser checks if the request contains a valid session
// and saves the session's userId in the context
func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("userId")

		if id == nil {
			e := apperrors.NewAuthorization(apperrors.InvalidSession)
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			c.Abort()
			return
		}

		userId := id.(string)

		c.Set("userId", userId)

		// Recreate session to extend its lifetime
		session.Set("userId", id)
		if err := session.Save(); err != nil {
			log.Printf("Failed recreate the session: %v\n", err.Error())
		}

		c.Next()
	}
}

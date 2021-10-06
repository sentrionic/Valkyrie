package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
)

func getAuthenticatedTestRouter(uid string) *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions(model.CookieName, store))

	router.Use(func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("userId", uid)
		c.Set("userId", uid)
	})

	return router
}

func getTestRouter() *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions(model.CookieName, store))
	return router
}

func getTestFieldErrorResponse(field, message string) gin.H {
	return gin.H{
		"errors": []model.FieldError{
			{
				Field:   field,
				Message: message,
			},
		},
	}
}

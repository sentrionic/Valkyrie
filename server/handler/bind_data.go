package handler

import (
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Request contains the validate function which validates the request with bindData
type Request interface {
	validate() error
}

// bindData is helper function, returns false if data is not bound
func bindData(c *gin.Context, req Request) bool {
	// Bind incoming json to struct and check for validation errors
	if err := c.ShouldBind(req); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return false
	}

	if err := req.validate(); err != nil {
		errors := strings.Split(err.Error(), ";")
		fErrors := make([]model.FieldError, 0)

		for _, e := range errors {
			split := strings.Split(e, ":")
			er := model.FieldError{
				Field:   strings.TrimSpace(split[0]),
				Message: strings.TrimSpace(split[1]),
			}
			fErrors = append(fErrors, er)
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"errors": fErrors,
		})
		return false
	}
	return true
}

package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

/*
 * AccountHandler contains all routes related to account actions (/api/account)
 * that the authenticated user can do
 */

// GetCurrent handler calls services for getting
// a user's details
// GetCurrent godoc
// @Tags Account
// @Summary Get Current User
// @Produce  json
// @Success 200 {object} model.User
// @Failure 404 {object} model.ErrorResponse
// @Router /account [get]
func (h *Handler) GetCurrent(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	user, err := h.userService.Get(userId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

type editReq struct {
	// Min 3, max 30 characters.
	Username string `form:"username"`
	// Must be unique
	Email string `form:"email"`
	// image/png or image/jpeg
	Image *multipart.FileHeader `form:"image" swaggertype:"string" format:"binary"`
} //@name EditUser

func (r editReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
	)
}

func (r *editReq) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
}

// Edit handler edits the users account details
// Edit godoc
// @Tags Account
// @Summary Update Current User
// @Accept mpfd
// @Produce  json
// @Param account body editReq true "Update Account"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /account [put]
func (h *Handler) Edit(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, h.MaxBodyBytes)

	var req editReq

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	authUser, err := h.userService.Get(userId)

	if err != nil {
		e := apperrors.NewAuthorization(apperrors.InvalidSession)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	authUser.Username = req.Username

	// New email, check if it's unique
	if authUser.Email != req.Email {
		inUse := h.userService.IsEmailAlreadyInUse(req.Email)

		if inUse {
			toFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}
		authUser.Email = req.Email
	}

	if req.Image != nil {

		// Validate image mime-type is allowable
		mimeType := req.Image.Header.Get("Content-Type")

		if valid := isAllowedImageType(mimeType); !valid {
			toFieldErrorResponse(c, "Image", apperrors.InvalidImageType)
			return
		}

		directory := fmt.Sprintf("valkyrie/users/%s", authUser.ID)
		url, err := h.userService.ChangeAvatar(req.Image, directory)

		if err != nil {
			e := apperrors.NewInternal()
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		_ = h.userService.DeleteImage(authUser.Image)

		authUser.Image = url
	}

	err = h.userService.UpdateAccount(authUser)

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, authUser)
}

type changeRequest struct {
	CurrentPassword string `json:"currentPassword"`
	// Min 6, max 150 characters.
	NewPassword string `json:"newPassword"`
	// Must be the same as the newPassword value.
	ConfirmNewPassword string `json:"confirmNewPassword"`
} //@name ChangePasswordRequest

func (r changeRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.CurrentPassword, validation.Required, validation.Length(6, 150)),
		validation.Field(&r.NewPassword, validation.Required, validation.Length(6, 150)),
		validation.Field(&r.ConfirmNewPassword, validation.Required, validation.Length(6, 150)),
	)
}

func (r *changeRequest) sanitize() {
	r.CurrentPassword = strings.TrimSpace(r.CurrentPassword)
	r.NewPassword = strings.TrimSpace(r.NewPassword)
	r.ConfirmNewPassword = strings.TrimSpace(r.ConfirmNewPassword)
}

// ChangePassword handler changes the user's password
// ChangePassword godoc
// @Tags Account
// @Summary Change Current User's Password
// @Accept json
// @Produce  json
// @Param request body changeRequest true "Change Password"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /account/change-password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	var req changeRequest

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	// Check if passwords are equal
	if req.NewPassword != req.ConfirmNewPassword {
		toFieldErrorResponse(c, "password", apperrors.PasswordsDoNotMatch)
		return
	}

	authUser, err := h.userService.Get(userId)

	if err != nil {
		e := apperrors.NewAuthorization(apperrors.InvalidSession)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	err = h.userService.ChangePassword(req.CurrentPassword, req.NewPassword, authUser)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

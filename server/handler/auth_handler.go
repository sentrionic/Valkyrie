package handler

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
	"strings"
)

/*
 * AuthHandler contains all routes related to account actions (/api/account)
 */

type registerReq struct {
	// Must be unique
	Email string `json:"email"`
	// Min 3, max 30 characters.
	Username string `json:"username"`
	// Min 6, max 150 characters.
	Password string `json:"password"`
} //@name RegisterRequest

func (r registerReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Username, validation.Required, validation.Length(3, 30)),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *registerReq) sanitize() {
	r.Username = strings.TrimSpace(r.Username)
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

// Register handler creates a new user
// Register godoc
// @Tags Account
// @Summary Create an Account
// @Accept  json
// @Produce  json
// @Param account body registerReq true "Create account"
// @Success 201 {object} model.User
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /account/register [post]
func (h *Handler) Register(c *gin.Context) {
	var req registerReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	initial := &model.User{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	}

	user, err := h.userService.Register(initial)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.DuplicateEmail).Error() {
			toFieldErrorResponse(c, "Email", apperrors.DuplicateEmail)
			return
		}
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusCreated, user)
}

type loginReq struct {
	// Must be unique
	Email string `json:"email"`
	// Min 6, max 150 characters.
	Password string `json:"password"`
} //@name LoginRequest

func (r loginReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
	)
}

func (r *loginReq) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
	r.Password = strings.TrimSpace(r.Password)
}

// Login used to authenticate existent user
// Login godoc
// @Tags Account
// @Summary User Login
// @Accept  json
// @Produce  json
// @Param account body loginReq true "Login account"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /account/login [post]
func (h *Handler) Login(c *gin.Context) {
	var req loginReq

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	user, err := h.userService.Login(req.Email, req.Password)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusOK, user)
}

// Logout handler removes the current session
// Logout godoc
// @Tags Account
// @Summary User Logout
// @Accept  json
// @Produce  json
// @Param account body loginReq true "Login account"
// @Success 200 {object} model.Success
// @Router /account/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	c.Set("user", nil)

	session := sessions.Default(c)
	session.Set("userId", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	err := session.Save()

	if err != nil {
		log.Printf("error clearing session: %v\n", err.Error())
	}

	c.JSON(http.StatusOK, true)
}

type forgotRequest struct {
	Email string `json:"email"`
} //@name ForgotPasswordRequest

func (r forgotRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.EmailFormat),
	)
}

func (r *forgotRequest) sanitize() {
	r.Email = strings.TrimSpace(r.Email)
	r.Email = strings.ToLower(r.Email)
}

// ForgotPassword sends a password reset email to the requested email
// ForgotPassword godoc
// @Tags Account
// @Summary Forgot Password Request
// @Accept  json
// @Produce  json
// @Param email body forgotRequest true "Forgot Password"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.ErrorsResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /account/forgot-password [post]
func (h *Handler) ForgotPassword(c *gin.Context) {
	var req forgotRequest
	if valid := bindData(c, &req); !valid {
		return
	}

	req.sanitize()

	user, err := h.userService.GetByEmail(req.Email)

	if err != nil {
		// No user with the email found
		if err.Error() == apperrors.NewNotFound("email", req.Email).Error() {
			c.JSON(http.StatusOK, true)
			return
		}

		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	ctx := c.Request.Context()
	err = h.userService.ForgotPassword(ctx, user)

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

type resetRequest struct {
	// The token the user got from the email.
	Token string `json:"token"`
	// Min 6, max 150 characters.
	Password string `json:"newPassword"`
	// Must be the same as the password value.
	ConfirmPassword string `json:"confirmNewPassword"`
} //@name ResetPasswordRequest

func (r resetRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token, validation.Required),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 150)),
		validation.Field(&r.ConfirmPassword, validation.Required, validation.Length(6, 150)),
	)
}

func (r *resetRequest) sanitize() {
	r.Token = strings.TrimSpace(r.Token)
	r.Password = strings.TrimSpace(r.Password)
	r.ConfirmPassword = strings.TrimSpace(r.ConfirmPassword)
}

// ResetPassword resets the user's password with the provided token
// ResetPassword godoc
// @Tags Account
// @Summary Reset Password
// @Accept  json
// @Produce  json
// @Param request body resetRequest true "Reset Password"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorsResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /account/reset-password [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var req resetRequest

	if valid := bindData(c, &req); !valid {
		return
	}

	req.sanitize()

	// Check if passwords match
	if req.Password != req.ConfirmPassword {
		toFieldErrorResponse(c, "Password", apperrors.PasswordsDoNotMatch)
		return
	}

	ctx := c.Request.Context()
	user, err := h.userService.ResetPassword(ctx, req.Password, req.Token)

	if err != nil {
		if err.Error() == apperrors.NewBadRequest(apperrors.InvalidResetToken).Error() {
			toFieldErrorResponse(c, "Token", apperrors.InvalidResetToken)
			return
		}
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	setUserSession(c, user.ID)

	c.JSON(http.StatusOK, user)
}

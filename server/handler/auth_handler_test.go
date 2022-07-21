package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/sentrionic/valkyrie/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_Register(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	user := fixture.GetMockUser()
	reqUser := &model.User{
		Email:    user.Email,
		Password: user.Password,
		Username: user.Username,
	}

	t.Run("Email, Username and Password Required", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "",
			"username": "",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Invalid email", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob",
			"username": "bobby",
			"password": "supersecret1234",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Signup")
	})

	t.Run("Username too short", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"username": "bo",
			"password": "superpassword",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Password too short", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"username": "bobby",
			"password": "supe",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Username too long", func(t *testing.T) {
		// We just want this to show that it's not called in this case
		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", mock.AnythingOfType("*model.User")).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    "bob@bob.com",
			"username": "kjhasiudaiusdiuadiuagszuidgaiszugdziasgdiazgsdiazugdipas",
			"password": "superpassword",
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 400, rr.Code)
		mockUserService.AssertNotCalled(t, "Register")
	})

	t.Run("Error returned from UserService", func(t *testing.T) {
		u := &model.User{
			Email:    reqUser.Email,
			Username: reqUser.Username,
			Password: reqUser.Password,
		}

		mockUserService := new(mocks.UserService)
		mockUserService.On("Register", u).Return(nil, apperrors.NewConflict("User Already Exists", u.Email))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// don't need a middleware as we don't yet have authorized user
		router := gin.Default()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"username": u.Username,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, 409, rr.Code)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Successful Creation", func(t *testing.T) {
		u := &model.User{
			Email:    reqUser.Email,
			Username: reqUser.Username,
			Password: reqUser.Password,
		}

		mockUserService := new(mocks.UserService)

		mockUserService.
			On("Register", u).
			Return(reqUser, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		// create a request body with empty email and password
		reqBody, err := json.Marshal(gin.H{
			"email":    u.Email,
			"username": u.Username,
			"password": u.Password,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(u)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertExpectations(t)
	})
}

func TestHandler_Login(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	// setup mock services, gin engine/router, handler layer
	mockUserService := new(mocks.UserService)

	router := getTestRouter()

	NewHandler(&Config{
		R:           router,
		UserService: mockUserService,
	})

	t.Run("Bad request data", func(t *testing.T) {
		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with invalid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    "notanemail",
			"password": "short",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		mockUserService.AssertNotCalled(t, "Login")
	})

	t.Run("Error Returned from UserService.Login", func(t *testing.T) {
		email := "bob@bob.com"
		password := "pwdoesnotmatch123"

		mockUSArgs := mock.Arguments{
			email, password,
		}

		// so we can check for a known status code
		mockError := apperrors.NewAuthorization("invalid email/password combo")

		mockUserService.On("Login", mockUSArgs...).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    email,
			"password": password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockUserService.AssertCalled(t, "Login", mockUSArgs...)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Successful Login", func(t *testing.T) {
		user := fixture.GetMockUser()

		mockUSArgs := mock.Arguments{
			user.Email,
			user.Password,
		}

		mockUserService.On("Login", mockUSArgs...).Return(user, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// create a request body with valid fields
		reqBody, err := json.Marshal(gin.H{
			"email":    user.Email,
			"password": user.Password,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(user)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockUserService.AssertCalled(t, "Login", mockUSArgs...)
	})
}

func TestHandler_Logout(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		uid := service.GenerateId()

		rr := httptest.NewRecorder()

		// creates a test context for setting a user
		router := getAuthenticatedTestRouter(uid)

		NewHandler(&Config{
			R: router,
		})

		request, _ := http.NewRequest(http.MethodPost, "/api/account/logout", nil)
		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(true)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		router.Use(func(c *gin.Context) {
			contextUserId, exists := c.Get("userId")
			assert.Equal(t, exists, false)
			assert.Nil(t, contextUserId)

			session := sessions.Default(c)
			id := session.Get("userId")
			assert.Nil(t, id)
		})
	})
}

func TestHandler_ForgotPassword(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()

	t.Run("ForgotPassword success", func(t *testing.T) {
		router := gin.Default()

		mockUserService := new(mocks.UserService)
		mockUserService.On("GetByEmail", mockUser.Email).Return(mockUser, nil)

		ForgotPasswordArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockUserService.
			On("ForgotPassword", ForgotPasswordArgs...).
			Return(nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email": mockUser.Email,
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPost, "/api/account/forgot-password", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(true)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "ForgotPassword", ForgotPasswordArgs...)
	})

	t.Run("ForgotPassword Failure", func(t *testing.T) {
		router := gin.Default()

		mockUserService := new(mocks.UserService)
		mockUserService.On("GetByEmail", mockUser.Email).Return(mockUser, nil)

		ForgotPasswordArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockError := apperrors.NewInternal()
		mockUserService.
			On("ForgotPassword", ForgotPasswordArgs...).
			Return(mockError)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email": mockUser.Email,
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPost, "/api/account/forgot-password", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "ForgotPassword", ForgotPasswordArgs...)
	})

	t.Run("No user found", func(t *testing.T) {
		router := gin.Default()

		mockUserService := new(mocks.UserService)

		mockError := apperrors.NewNotFound("email", mockUser.Email)
		mockUserService.On("GetByEmail", mockUser.Email).Return(&model.User{}, mockError)

		ForgotPasswordArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser,
		}

		mockUserService.
			On("ForgotPassword", ForgotPasswordArgs...).
			Return(nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"email": mockUser.Email,
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPost, "/api/account/forgot-password", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(true)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertNotCalled(t, "ForgotPassword", ForgotPasswordArgs...)
	})
}

func TestHandler_ForgotPassword_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()

	router := gin.Default()

	mockUserService := new(mocks.UserService)
	mockUserService.On("GetByEmail", mockUser.Email).Return(mockUser, nil)

	NewHandler(&Config{
		R:            router,
		UserService:  mockUserService,
		MaxBodyBytes: 4 * 1024 * 1024,
	})

	testCases := []struct {
		name string
		body gin.H
	}{
		{
			name: "Email required",
			body: gin.H{},
		},
		{
			name: "Invalid Email",
			body: gin.H{
				"email": "invalidemail",
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/api/account/forgot-password", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockUserService.AssertNotCalled(t, "ForgotPassword")
		})
	}
}

func TestHandler_ResetPassword(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	token := fixture.RandStringRunes(18)

	t.Run("ResetPassword success", func(t *testing.T) {
		mockUserService := new(mocks.UserService)

		router := getTestRouter()

		ResetPasswordArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.Password,
			token,
		}

		mockUserService.
			On("ResetPassword", ResetPasswordArgs...).
			Return(mockUser, nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"token":              token,
			"newPassword":        mockUser.Password,
			"confirmNewPassword": mockUser.Password,
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPost, "/api/account/reset-password", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(mockUser)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "ResetPassword", ResetPasswordArgs...)
	})

	t.Run("ResetPassword Failure", func(t *testing.T) {
		mockUserService := new(mocks.UserService)

		router := getTestRouter()

		ResetPasswordArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockUser.Password,
			token,
		}

		mockError := apperrors.NewInternal()
		mockUserService.
			On("ResetPassword", ResetPasswordArgs...).
			Return(nil, mockError)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"token":              token,
			"newPassword":        mockUser.Password,
			"confirmNewPassword": mockUser.Password,
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPost, "/api/account/reset-password", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "ResetPassword", ResetPasswordArgs...)
	})
}

func TestHandler_ResetPassword_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	router := gin.Default()

	mockUserService := new(mocks.UserService)

	NewHandler(&Config{
		R:           router,
		UserService: mockUserService,
	})

	password := fixture.RandStringRunes(6)
	confirmPassword := password

	testCases := []struct {
		name string
		body gin.H
	}{
		{
			name: "Token required",
			body: gin.H{
				"newPassword":        password,
				"confirmNewPassword": password,
			},
		},
		{
			name: "Password required",
			body: gin.H{
				"token":              fixture.RandStringRunes(18),
				"confirmNewPassword": password,
			},
		},
		{
			name: "Password too short",
			body: gin.H{
				"token":              fixture.RandStringRunes(18),
				"newPassword":        fixture.RandStringRunes(5),
				"confirmNewPassword": confirmPassword,
			},
		},
		{
			name: "NewPassword too long",
			body: gin.H{
				"token":              fixture.RandStringRunes(16),
				"newPassword":        fixture.RandStringRunes(151),
				"confirmNewPassword": confirmPassword,
			},
		},
		{
			name: "ConfirmNewPassword too short",
			body: gin.H{
				"token":              fixture.RandStringRunes(6),
				"newPassword":        fixture.RandStringRunes(6),
				"confirmNewPassword": fixture.RandStringRunes(5),
			},
		},
		{
			name: "ConfirmNewPassword too long",
			body: gin.H{
				"currentPassword":    fixture.RandStringRunes(6),
				"newPassword":        fixture.RandStringRunes(6),
				"confirmNewPassword": fixture.RandStringRunes(151),
			},
		},
		{
			name: "ConfirmNewPassword required",
			body: gin.H{
				"token":       fixture.RandStringRunes(6),
				"newPassword": fixture.RandStringRunes(6),
			},
		},
		{
			name: "NewPassword and ConfirmNewPassword not equal",
			body: gin.H{
				"token":              fixture.RandStringRunes(6),
				"newPassword":        fixture.RandStringRunes(6),
				"confirmNewPassword": fixture.RandStringRunes(6),
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/api/account/reset-password", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockUserService.AssertNotCalled(t, "ResetPassword")
		})
	}
}

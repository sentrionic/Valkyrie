package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/sentrionic/valkyrie/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestHandler_GetCurrent(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	t.Run("Success", func(t *testing.T) {
		uid := service.GenerateId()

		mockUserResp := fixture.GetMockUser()
		mockUserResp.ID = uid

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUserResp, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(uid)

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockUserResp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		uid := service.GenerateId()
		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(nil, fmt.Errorf("some error down call chain"))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(uid)

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respErr := apperrors.NewNotFound("user", uid)

		respBody, err := json.Marshal(gin.H{
			"error": respErr,
		})
		assert.NoError(t, err)

		assert.Equal(t, respErr.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertExpectations(t) // assert that UserService.Get was called
	})

	t.Run("Unauthorized", func(t *testing.T) {
		uid := service.GenerateId()
		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(nil, nil)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertNotCalled(t, "Get", uid)
	})
}

func TestHandler_EditAccount(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid := service.GenerateId()
	mockUser := fixture.GetMockUser()
	mockUser.ID = uid

	t.Run("Unauthorized", func(t *testing.T) {
		router := getTestRouter()
		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		rr := httptest.NewRecorder()

		newName := fixture.Username()
		newEmail := fixture.Email()

		form := url.Values{}
		form.Add("username", newName)
		form.Add("email", newEmail)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertNotCalled(t, "UpdateAccount")
	})

	t.Run("UpdateAccount success", func(t *testing.T) {
		router := getAuthenticatedTestRouter(uid)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		newName := fixture.Username()
		newEmail := fixture.Email()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_ = writer.WriteField("username", newName)
		_ = writer.WriteField("email", newEmail)

		_ = writer.Close()

		request, _ := http.NewRequest(http.MethodPut, "/api/account", body)
		request.Header.Set("Content-Type", writer.FormDataContentType())

		mockUser.Username = newName
		mockUser.Email = newEmail

		UpdateAccountArgs := mock.Arguments{
			mockUser,
		}

		dbImageURL := "https://website.com/696292a38f493a4283d1a308e4a11732/84d81/Profile.jpg"

		mockUserService.
			On("UpdateAccount", UpdateAccountArgs...).
			Run(func(args mock.Arguments) {
				userArg := args.Get(0).(*model.User)
				userArg.Image = dbImageURL
			}).
			Return(nil)

		router.ServeHTTP(rr, request)

		mockUser.Image = dbImageURL
		respBody, _ := json.Marshal(mockUser)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "UpdateAccount", UpdateAccountArgs...)
	})

	t.Run("UpdateAccount Failure", func(t *testing.T) {
		router := getAuthenticatedTestRouter(uid)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		form := url.Values{}
		form.Add("username", mockUser.Username)
		form.Add("email", mockUser.Email)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		mockError := apperrors.NewInternal()

		mockUserService.
			On("UpdateAccount", mockUser).
			Return(mockError)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "UpdateAccount", mockUser)
	})

	t.Run("Disallowed mimetype", func(t *testing.T) {
		router := getAuthenticatedTestRouter(uid)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		multipartImageFixture := fixture.NewMultipartImage("image.txt", "mage/svg+xml")
		defer multipartImageFixture.Close()

		request, _ := http.NewRequest(http.MethodPut, "/api/account", multipartImageFixture.MultipartBody)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockUserService.AssertNotCalled(t, "ChangeAvatar")
	})

	t.Run("Email already in use", func(t *testing.T) {
		router := getAuthenticatedTestRouter(uid)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		form := url.Values{}

		duplicateEmail := "duplicate@example.com"
		form.Add("username", mockUser.Username)
		form.Add("email", duplicateEmail)

		request, _ := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))
		request.Form = form

		mockUserService.
			On("IsEmailAlreadyInUse", duplicateEmail).
			Return(true)

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(getTestFieldErrorResponse("Email", apperrors.DuplicateEmail))

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertNotCalled(t, "UpdateAccount", mockUser)
	})
}

func TestHandler_ChangePassword(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid := service.GenerateId()
	mockUser := fixture.GetMockUser()
	mockUser.ID = uid

	t.Run("Unauthorized", func(t *testing.T) {
		router := getTestRouter()
		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"currentPassword":    "password",
			"newPassword":        "password!",
			"confirmNewPassword": "password!",
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPut, "/api/account/change-password", bytes.NewBuffer(reqBody))

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertNotCalled(t, "ChangePassword")
	})

	t.Run("ChangePassword success", func(t *testing.T) {
		router := getAuthenticatedTestRouter(uid)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		currentPassword := mockUser.Password
		newPassword := "password!"

		ChangePasswordArgs := mock.Arguments{
			currentPassword,
			newPassword,
			mockUser,
		}

		mockUserService.
			On("ChangePassword", ChangePasswordArgs...).
			Return(nil)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"currentPassword":    currentPassword,
			"newPassword":        newPassword,
			"confirmNewPassword": newPassword,
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPut, "/api/account/change-password", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(true)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "ChangePassword", ChangePasswordArgs...)
	})

	t.Run("ChangePassword Failure", func(t *testing.T) {
		router := getAuthenticatedTestRouter(uid)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", uid).Return(mockUser, nil)

		currentPassword := mockUser.Password
		newPassword := "password!"

		ChangePasswordArgs := mock.Arguments{
			currentPassword,
			newPassword,
			mockUser,
		}

		mockError := apperrors.NewInternal()
		mockUserService.
			On("ChangePassword", ChangePasswordArgs...).
			Return(mockError)

		NewHandler(&Config{
			R:            router,
			UserService:  mockUserService,
			MaxBodyBytes: 4 * 1024 * 1024,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"currentPassword":    currentPassword,
			"newPassword":        newPassword,
			"confirmNewPassword": newPassword,
		})
		assert.NoError(t, err)

		request, _ := http.NewRequest(http.MethodPut, "/api/account/change-password", bytes.NewBuffer(reqBody))
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockUserService.AssertCalled(t, "ChangePassword", ChangePasswordArgs...)
	})
}

func TestHandler_ChangePassword_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	uid := service.GenerateId()
	mockUser := fixture.GetMockUser()
	mockUser.ID = uid

	router := getAuthenticatedTestRouter(uid)

	mockUserService := new(mocks.UserService)
	mockUserService.On("Get", uid).Return(mockUser, nil)

	NewHandler(&Config{
		R:            router,
		UserService:  mockUserService,
		MaxBodyBytes: 4 * 1024 * 1024,
	})

	password := fixture.RandStringRunes(6)
	confirmPassword := password

	testCases := []struct {
		name string
		body gin.H
	}{
		{
			name: "CurrentPassword required",
			body: gin.H{
				"newPassword":        password,
				"confirmNewPassword": confirmPassword,
			},
		},
		{
			name: "NewPassword too short",
			body: gin.H{
				"currentPassword":    fixture.RandStringRunes(6),
				"newPassword":        fixture.RandStringRunes(5),
				"confirmNewPassword": confirmPassword,
			},
		},
		{
			name: "NewPassword too long",
			body: gin.H{
				"currentPassword":    fixture.RandStringRunes(6),
				"newPassword":        fixture.RandStringRunes(151),
				"confirmNewPassword": confirmPassword,
			},
		},
		{
			name: "NewPassword required",
			body: gin.H{
				"currentPassword":    fixture.RandStringRunes(6),
				"confirmNewPassword": fixture.RandStringRunes(6),
			},
		},
		{
			name: "ConfirmNewPassword too short",
			body: gin.H{
				"currentPassword":    fixture.RandStringRunes(6),
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
				"currentPassword": fixture.RandStringRunes(6),
				"newPassword":     fixture.RandStringRunes(6),
			},
		},
		{
			name: "NewPassword and ConfirmNewPassword not equal",
			body: gin.H{
				"currentPassword":    fixture.RandStringRunes(6),
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

			request, err := http.NewRequest(http.MethodPut, "/api/account/change-password", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockUserService.AssertNotCalled(t, "ChangePassword")
		})
	}
}

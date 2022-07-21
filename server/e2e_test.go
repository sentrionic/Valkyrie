package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sentrionic/valkyrie/config"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func setupTest(t *testing.T) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	ctx := context.Background()

	// Load config for local testing, ignore for CI
	_ = godotenv.Load()

	cfg, err := config.LoadConfig(ctx)
	assert.NoError(t, err)

	ds, err := initDS(ctx, cfg)
	assert.NoError(t, err)

	router, err := inject(ds, cfg)
	assert.NoError(t, err)

	return router
}

func TestMain_AccountE2E(t *testing.T) {
	router := setupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	mockUsername := fixture.Username()
	newPassword := fixture.RandStr(10)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": authUser.Username,
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")
				authUser.ID = respBody.ID
			},
		},
		{
			name: "Login Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, authUser.ID, respBody.ID)
				assert.NotNil(t, respBody.Image)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")

				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Get Account",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/account", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, authUser.ID, respBody.ID)
				assert.NotNil(t, respBody.Image)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
			},
		},
		{
			name: "Edit Account",
			setupRequest: func() (*http.Request, error) {
				form := url.Values{}
				form.Add("username", mockUsername)
				form.Add("email", authUser.Email)

				request, err := http.NewRequest(http.MethodPut, "/api/account", strings.NewReader(form.Encode()))

				if err != nil {
					return nil, err
				}

				request.Form = form

				return request, nil
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockUsername, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, authUser.ID, respBody.ID)
				assert.NotNil(t, respBody.Image)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
			},
		},
		{
			name: "Change Password",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"currentPassword":    authUser.Password,
					"newPassword":        newPassword,
					"confirmNewPassword": newPassword,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPut, "/api/account/change-password", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Login to verify changes",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"password": newPassword,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/login", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockUsername, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, authUser.ID, respBody.ID)
				assert.NotNil(t, respBody.Image)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
			},
		},
		{
			name: "Logout",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodPost, "/api/account/logout", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "Check session is invalid now",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/account", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}

func TestMain_FriendsE2E(t *testing.T) {
	router := setupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	mockUser := fixture.GetMockUser()
	mockUserCookie := ""

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": authUser.Username,
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")
				authUser.ID = respBody.ID
				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Register Friend",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": mockUser.Username,
					"password": mockUser.Password,
					"email":    mockUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockUser.Username, respBody.Username)
				assert.Equal(t, mockUser.Email, respBody.Email)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, mockUser.Image, respBody.Image)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				mockUser.ID = respBody.ID
				mockUserCookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Send friend request",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
				return http.NewRequest(http.MethodPost, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Get friend requests",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/account/me/pending", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.FriendRequest{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				requests := *respBody
				assert.Equal(t, 1, len(requests))
				request := requests[0]

				assert.Equal(t, authUser.Username, request.Username)
				assert.Equal(t, authUser.ID, request.Id)
				assert.Equal(t, authUser.Image, request.Image)
				assert.Equal(t, model.Incoming, request.Type)
			},
		},
		{
			name: "Accept friend request",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/account/%s/friend/accept", authUser.ID)
				return http.NewRequest(http.MethodPost, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Check that requests are empty now",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/account/me/pending", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.FriendRequest{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				requests := *respBody
				assert.Equal(t, 0, len(requests))
			},
		},
		{
			name: "Get friends",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/account/me/friends", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.Friend{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				friends := *respBody
				assert.Equal(t, 1, len(friends))
				friend := friends[0]

				assert.Equal(t, authUser.Username, friend.Username)
				assert.Equal(t, authUser.ID, friend.Id)
				assert.Equal(t, authUser.Image, friend.Image)
				assert.True(t, friend.IsOnline)
			},
		},
		{
			name: "Remove friend",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/account/%s/friend", authUser.ID)
				return http.NewRequest(http.MethodDelete, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Confirm friends list is empty",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/account/me/friends", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.Friend{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				friends := *respBody
				assert.Equal(t, 0, len(friends))
			},
		},
		{
			name: "Send another friend request",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
				return http.NewRequest(http.MethodPost, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Get authUser's friend requests",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/account/me/pending", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.FriendRequest{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				requests := *respBody
				assert.Equal(t, 1, len(requests))
				request := requests[0]

				assert.Equal(t, mockUser.Username, request.Username)
				assert.Equal(t, mockUser.ID, request.Id)
				assert.Equal(t, mockUser.Image, request.Image)
				assert.Equal(t, model.Outgoing, request.Type)
			},
		},
		{
			name: "Cancel friend request",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/account/%s/friend/accept", mockUser.ID)
				return http.NewRequest(http.MethodPost, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}

func TestMain_GuildsE2E(t *testing.T) {
	router := setupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	mockUser := fixture.GetMockUser()
	mockUserCookie := ""

	mockGuild := fixture.GetMockGuild("")
	inviteLink := ""

	mockName := fixture.Username()

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": authUser.Username,
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")
				authUser.ID = respBody.ID
				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Register Member",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": mockUser.Username,
					"password": mockUser.Password,
					"email":    mockUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockUser.Username, respBody.Username)
				assert.Equal(t, mockUser.Email, respBody.Email)
				assert.True(t, respBody.IsOnline)
				assert.Equal(t, mockUser.Image, respBody.Image)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				mockUser.ID = respBody.ID
				mockUserCookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Create Guild",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"name": mockGuild.Name,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockGuild.Name, respBody.Name)
				assert.Equal(t, authUser.ID, respBody.OwnerId)
				assert.NotNil(t, respBody.Id)
				assert.Nil(t, respBody.Icon)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.False(t, respBody.HasNotification)
				assert.NotNil(t, respBody.DefaultChannelId)

				mockGuild.ID = respBody.Id
				mockGuild.OwnerId = authUser.ID
			},
		},
		{
			name: "Get authUser's guilds",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Equal(t, 1, len(guilds))
				guild := guilds[0]

				assert.Equal(t, mockGuild.Name, guild.Name)
				assert.Equal(t, mockGuild.ID, guild.Id)
				assert.Equal(t, mockGuild.OwnerId, guild.OwnerId)
				assert.Nil(t, guild.Icon)
				assert.NotNil(t, guild.CreatedAt)
				assert.NotNil(t, guild.UpdatedAt)
				assert.False(t, guild.HasNotification)
				assert.NotNil(t, guild.DefaultChannelId)
			},
		},
		{
			name: "Edit Guild",
			setupRequest: func() (*http.Request, error) {
				form := url.Values{}
				form.Add("name", mockName)

				request, err := http.NewRequest(http.MethodPut, "/api/guilds/"+mockGuild.ID, strings.NewReader(form.Encode()))

				if err != nil {
					return nil, err
				}

				request.Form = form

				return request, nil
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)
			},
		},
		{
			name: "Get guild invite",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
				return http.NewRequest(http.MethodGet, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				err := json.Unmarshal(recorder.Body.Bytes(), &inviteLink)
				assert.NoError(t, err)
				assert.NotNil(t, inviteLink)
			},
		},
		{
			name: "Join Guild",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"link": inviteLink,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockName, respBody.Name)
				assert.Equal(t, authUser.ID, respBody.OwnerId)
				assert.NotNil(t, respBody.Id)
				assert.Nil(t, respBody.Icon)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.False(t, respBody.HasNotification)
				assert.NotNil(t, respBody.DefaultChannelId)
			},
		},
		{
			name: "Check mockUser is in the guild",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Equal(t, 1, len(guilds))
				guild := guilds[0]

				assert.Equal(t, mockName, guild.Name)
				assert.Equal(t, mockGuild.ID, guild.Id)
				assert.Equal(t, mockGuild.OwnerId, guild.OwnerId)
				assert.Nil(t, guild.Icon)
				assert.NotNil(t, guild.CreatedAt)
				assert.NotNil(t, guild.UpdatedAt)
				assert.False(t, guild.HasNotification)
				assert.NotNil(t, guild.DefaultChannelId)
			},
		},
		{
			name: "Leave Guild",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodDelete, "/api/guilds/"+mockGuild.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Check mockUser is in no guilds anymore",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Equal(t, 0, len(guilds))
			},
		},
		{
			name: "Get guild members",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/guilds/%s/members", mockGuild.ID)
				return http.NewRequest(http.MethodGet, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.MemberResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				members := *respBody
				assert.Equal(t, 1, len(members))
				member := members[0]

				assert.Equal(t, authUser.Username, member.Username)
				assert.Equal(t, authUser.ID, member.Id)
				assert.True(t, member.IsOnline)
				assert.NotNil(t, member.Image)
				assert.NotNil(t, member.CreatedAt)
				assert.NotNil(t, member.UpdatedAt)
				assert.Nil(t, member.Nickname)
				assert.Nil(t, member.Color)
			},
		},
		{
			name: "Delete Guild",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/guilds/%s/delete", mockGuild.ID)
				return http.NewRequest(http.MethodDelete, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Verify that there are no more guilds",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Len(t, guilds, 0)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}

func TestMain_MembersE2E(t *testing.T) {
	router := setupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	mockUser := fixture.GetMockUser()
	mockUserCookie := ""

	mockGuild := fixture.GetMockGuild("")
	inviteLink := ""

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": authUser.Username,
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")
				authUser.ID = respBody.ID
				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Register Member",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": mockUser.Username,
					"password": mockUser.Password,
					"email":    mockUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockUser.Username, respBody.Username)
				assert.Equal(t, mockUser.Email, respBody.Email)
				assert.Equal(t, mockUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.Image)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				mockUser.ID = respBody.ID
				mockUserCookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Create Guild",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"name": mockGuild.Name,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockGuild.Name, respBody.Name)
				assert.Equal(t, authUser.ID, respBody.OwnerId)
				assert.NotNil(t, respBody.Id)
				assert.Nil(t, respBody.Icon)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.False(t, respBody.HasNotification)
				assert.NotNil(t, respBody.DefaultChannelId)

				mockGuild.ID = respBody.Id
				mockGuild.OwnerId = authUser.ID
			},
		},
		{
			name: "Edit authUser's member settings",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"nickname": authUser.Username,
					"color":    "#fff",
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
				return http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)
			},
		},
		{
			name: "Get authUser's member settings",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
				return http.NewRequest(http.MethodGet, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &model.MemberSettings{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, *respBody.Nickname)
				assert.Equal(t, "#fff", *respBody.Color)
			},
		},
		{
			name: "Get guild invite",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/guilds/%s/invite?isPermanent=true", mockGuild.ID)
				return http.NewRequest(http.MethodGet, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				err := json.Unmarshal(recorder.Body.Bytes(), &inviteLink)
				assert.NoError(t, err)
				assert.NotNil(t, inviteLink)
			},
		},
		{
			name: "Join Guild",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"link": inviteLink,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockGuild.Name, respBody.Name)
				assert.Equal(t, authUser.ID, respBody.OwnerId)
				assert.NotNil(t, respBody.Id)
				assert.Nil(t, respBody.Icon)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.False(t, respBody.HasNotification)
				assert.NotNil(t, respBody.DefaultChannelId)
			},
		},
		{
			name: "Check mockUser is in the guild",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Equal(t, 1, len(guilds))
				guild := guilds[0]

				assert.Equal(t, mockGuild.Name, guild.Name)
				assert.Equal(t, mockGuild.ID, guild.Id)
				assert.Equal(t, mockGuild.OwnerId, guild.OwnerId)
				assert.Nil(t, guild.Icon)
				assert.NotNil(t, guild.CreatedAt)
				assert.NotNil(t, guild.UpdatedAt)
				assert.False(t, guild.HasNotification)
				assert.NotNil(t, guild.DefaultChannelId)
			},
		},
		{
			name: "Kick member",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"memberId": mockUser.ID,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
				return http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)
			},
		},
		{
			name: "Check mockUser is in no guilds anymore",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Equal(t, 0, len(guilds))
			},
		},
		{
			name: "Rejoin Guild",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"link": inviteLink,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockGuild.Name, respBody.Name)
				assert.Equal(t, authUser.ID, respBody.OwnerId)
				assert.NotNil(t, respBody.Id)
				assert.Nil(t, respBody.Icon)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.False(t, respBody.HasNotification)
				assert.NotNil(t, respBody.DefaultChannelId)
			},
		},
		{
			name: "Check mockUser is in the guild",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Equal(t, 1, len(guilds))
				guild := guilds[0]

				assert.Equal(t, mockGuild.Name, guild.Name)
				assert.Equal(t, mockGuild.ID, guild.Id)
				assert.Equal(t, mockGuild.OwnerId, guild.OwnerId)
				assert.Nil(t, guild.Icon)
				assert.NotNil(t, guild.CreatedAt)
				assert.NotNil(t, guild.UpdatedAt)
				assert.False(t, guild.HasNotification)
				assert.NotNil(t, guild.DefaultChannelId)
			},
		},
		{
			name: "Ban member",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"memberId": mockUser.ID,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
				return http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)
			},
		},
		{
			name: "Check mockUser is in no guilds anymore",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/guilds", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", mockUserCookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				guilds := *respBody
				assert.Equal(t, 0, len(guilds))
			},
		},
		{
			name: "Get guild ban list",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
				return http.NewRequest(http.MethodGet, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.BanResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				bans := *respBody
				assert.Equal(t, 1, len(bans))
				ban := bans[0]

				assert.Equal(t, mockUser.Username, ban.Username)
				assert.Equal(t, mockUser.ID, ban.Id)
				assert.NotNil(t, ban.Image)
			},
		},
		{
			name: "Unban member",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"memberId": mockUser.ID,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
				return http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)
			},
		},
		{
			name: "Verify ban list is empty",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
				return http.NewRequest(http.MethodGet, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.BanResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				bans := *respBody
				assert.Equal(t, 0, len(bans))
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}

func TestMain_ChannelsE2E(t *testing.T) {
	router := setupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	mockGuild := fixture.GetMockGuild("")

	mockChannel := fixture.GetMockChannel("")

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": authUser.Username,
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")
				authUser.ID = respBody.ID
				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Create Guild",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"name": mockGuild.Name,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockGuild.Name, respBody.Name)
				assert.Equal(t, authUser.ID, respBody.OwnerId)
				assert.NotNil(t, respBody.Id)
				assert.Nil(t, respBody.Icon)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.False(t, respBody.HasNotification)
				assert.NotNil(t, respBody.DefaultChannelId)

				mockGuild.ID = respBody.Id
				mockGuild.OwnerId = authUser.ID
			},
		},
		{
			name: "Create Channel",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"name": mockChannel.Name,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/channels/"+mockGuild.ID, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.ChannelResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockChannel.Name, respBody.Name)
				assert.NotNil(t, respBody.Id)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.True(t, respBody.IsPublic)
				assert.False(t, respBody.HasNotification)

				mockChannel.ID = respBody.Id
				mockChannel.GuildID = &mockGuild.ID
			},
		},
		{
			name: "Get guild channels",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/channels/"+mockGuild.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.ChannelResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				channels := *respBody
				assert.Equal(t, 2, len(channels))
				channel := channels[1]

				assert.Equal(t, mockChannel.Name, channel.Name)
				assert.Equal(t, mockChannel.ID, channel.Id)
				assert.NotNil(t, channel.CreatedAt)
				assert.NotNil(t, channel.UpdatedAt)
				assert.True(t, channel.IsPublic)
				assert.False(t, channel.HasNotification)
			},
		},
		{
			name: "Edit Channel",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"name":     mockChannel.Name,
					"isPublic": false,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPut, "/api/channels/"+mockChannel.ID, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Check that the channel is private now",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/channels/"+mockGuild.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.ChannelResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				channels := *respBody
				assert.Equal(t, 2, len(channels))
				channel := channels[1]

				assert.Equal(t, mockChannel.Name, channel.Name)
				assert.Equal(t, mockChannel.ID, channel.Id)
				assert.NotNil(t, channel.CreatedAt)
				assert.NotNil(t, channel.UpdatedAt)
				assert.False(t, channel.IsPublic)
				assert.False(t, channel.HasNotification)
			},
		},
		{
			name: "Get private channel members",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/channels/%s/members", mockChannel.ID)
				return http.NewRequest(http.MethodGet, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				var respBody []string
				err := json.Unmarshal(recorder.Body.Bytes(), &respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				assert.Equal(t, 1, len(respBody))
				id := respBody[0]

				assert.Equal(t, authUser.ID, id)
			},
		},
		{
			name: "Delete Channel",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodDelete, "/api/channels/"+mockChannel.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Check that the channel is got deleted",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/channels/"+mockGuild.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.ChannelResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				channels := *respBody
				assert.Equal(t, 1, len(channels))
				channel := channels[0]

				assert.Equal(t, "general", channel.Name)
				assert.NotNil(t, channel.Id)
				assert.NotNil(t, channel.CreatedAt)
				assert.NotNil(t, channel.UpdatedAt)
				assert.True(t, channel.IsPublic)
				assert.False(t, channel.HasNotification)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}

func TestMain_DMsE2E(t *testing.T) {
	router := setupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	mockUser := fixture.GetMockUser()

	dmId := ""

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": authUser.Username,
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")
				authUser.ID = respBody.ID
				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Register Member",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": mockUser.Username,
					"password": mockUser.Password,
					"email":    mockUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockUser.Username, respBody.Username)
				assert.Equal(t, mockUser.Email, respBody.Email)
				assert.Equal(t, mockUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.Image)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				mockUser.ID = respBody.ID
			},
		},
		{
			name: "Start a DM",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/channels/%s/dm", mockUser.ID)
				return http.NewRequest(http.MethodPost, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &model.DirectMessage{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.NotNil(t, respBody.User)
				assert.NotNil(t, respBody.Id)
				assert.Equal(t, mockUser.Username, respBody.User.Username)
				assert.Equal(t, mockUser.ID, respBody.User.Id)
				assert.Equal(t, mockUser.Image, respBody.User.Image)
				assert.True(t, respBody.User.IsOnline)
				assert.False(t, respBody.User.IsFriend)

				dmId = respBody.Id
			},
		},
		{
			name: "Get authUser's DMs",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/channels/me/dm", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.DirectMessage{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				dms := *respBody
				assert.Len(t, dms, 1)
				dm := dms[0]

				assert.NotNil(t, dm.User)
				assert.NotNil(t, dm.Id)
				assert.Equal(t, mockUser.Username, dm.User.Username)
				assert.Equal(t, mockUser.ID, dm.User.Id)
				assert.Equal(t, mockUser.Image, dm.User.Image)
				assert.True(t, dm.User.IsOnline)
				assert.False(t, dm.User.IsFriend)
			},
		},
		{
			name: "Close DM",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/channels/%s/dm", dmId)
				return http.NewRequest(http.MethodDelete, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Verify the user does not have any open DMs",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/channels/me/dm", nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.DirectMessage{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				dms := *respBody
				assert.Len(t, dms, 0)
			},
		},
		{
			name: "Get the already existing DM",
			setupRequest: func() (*http.Request, error) {
				reqUrl := fmt.Sprintf("/api/channels/%s/dm", mockUser.ID)
				return http.NewRequest(http.MethodPost, reqUrl, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
				request.Header.Add("Cookie", cookie)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &model.DirectMessage{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.NotNil(t, respBody.User)
				assert.NotNil(t, respBody.Id)
				assert.Equal(t, mockUser.Username, respBody.User.Username)
				assert.Equal(t, mockUser.ID, respBody.User.Id)
				assert.Equal(t, mockUser.Image, respBody.User.Image)
				assert.True(t, respBody.User.IsOnline)
				assert.False(t, respBody.User.IsFriend)
				assert.Equal(t, dmId, respBody.Id)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}

func TestMain_MessagesE2E(t *testing.T) {
	router := setupTest(t)

	authUser := fixture.GetMockUser()
	cookie := ""

	mockGuild := fixture.GetMockGuild("")
	mockChannel := fixture.GetMockChannel("")

	mockMessage := fixture.GetMockMessage("", "")
	mockText := fixture.RandStringRunes(10)

	testCases := []struct {
		name          string
		setupRequest  func() (*http.Request, error)
		setupHeaders  func(t *testing.T, request *http.Request)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "Register Account",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"username": authUser.Username,
					"password": authUser.Password,
					"email":    authUser.Email,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/account/register", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.User{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, authUser.Username, respBody.Username)
				assert.Equal(t, authUser.Email, respBody.Email)
				assert.Equal(t, authUser.Image, respBody.Image)
				assert.True(t, respBody.IsOnline)
				assert.NotNil(t, respBody.ID)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)

				assert.Contains(t, recorder.Header(), "Set-Cookie")
				authUser.ID = respBody.ID
				cookie = recorder.Header().Get("Set-Cookie")
			},
		},
		{
			name: "Create Guild",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"name": mockGuild.Name,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.GuildResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockGuild.Name, respBody.Name)
				assert.Equal(t, authUser.ID, respBody.OwnerId)
				assert.NotNil(t, respBody.Id)
				assert.Nil(t, respBody.Icon)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.False(t, respBody.HasNotification)
				assert.NotNil(t, respBody.DefaultChannelId)

				mockGuild.ID = respBody.Id
				mockGuild.OwnerId = authUser.ID
			},
		},
		{
			name: "Create Channel",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"name": mockChannel.Name,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/channels/"+mockGuild.ID, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody := &model.ChannelResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)

				assert.Equal(t, mockChannel.Name, respBody.Name)
				assert.NotNil(t, respBody.Id)
				assert.NotNil(t, respBody.CreatedAt)
				assert.NotNil(t, respBody.UpdatedAt)
				assert.True(t, respBody.IsPublic)
				assert.False(t, respBody.HasNotification)

				mockChannel.ID = respBody.Id
				mockChannel.GuildID = &mockGuild.ID
			},
		},
		{
			name: "Create Message",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"text": *mockMessage.Text,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPost, "/api/messages/"+mockChannel.ID, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Get channel messages",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/messages/"+mockChannel.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.MessageResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				messages := *respBody
				assert.Len(t, messages, 1)
				message := messages[0]

				assert.Equal(t, message.Text, mockMessage.Text)
				assert.NotNil(t, message.Id)
				assert.NotNil(t, message.CreatedAt)
				assert.NotNil(t, message.UpdatedAt)
				assert.NotNil(t, message.User)
				assert.Nil(t, message.Attachment)

				author := message.User
				assert.Equal(t, authUser.Username, author.Username)
				assert.Equal(t, authUser.ID, author.Id)
				assert.True(t, author.IsOnline)
				assert.NotNil(t, author.Image)
				assert.NotNil(t, author.CreatedAt)
				assert.NotNil(t, author.UpdatedAt)
				assert.Nil(t, author.Nickname)
				assert.Nil(t, author.Color)

				mockMessage.ID = message.Id
			},
		},
		{
			name: "Edit Message",
			setupRequest: func() (*http.Request, error) {
				data := gin.H{
					"text": mockText,
				}

				reqBody, err := json.Marshal(data)
				assert.NoError(t, err)

				return http.NewRequest(http.MethodPut, "/api/messages/"+mockMessage.ID, bytes.NewBuffer(reqBody))
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Verify message got edited",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/messages/"+mockChannel.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.MessageResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				messages := *respBody
				assert.Len(t, messages, 1)
				message := messages[0]

				assert.Equal(t, *message.Text, mockText)
				assert.NotNil(t, message.Id)
				assert.NotNil(t, message.CreatedAt)
				assert.NotNil(t, message.UpdatedAt)
				assert.NotNil(t, message.User)
				assert.Nil(t, message.Attachment)

				author := message.User
				assert.Equal(t, authUser.Username, author.Username)
				assert.Equal(t, authUser.ID, author.Id)
				assert.True(t, author.IsOnline)
				assert.NotNil(t, author.Image)
				assert.NotNil(t, author.CreatedAt)
				assert.NotNil(t, author.UpdatedAt)
				assert.Nil(t, author.Nickname)
				assert.Nil(t, author.Color)
			},
		},
		{
			name: "Delete Message",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody, err := json.Marshal(true)
				assert.NoError(t, err)

				assert.Equal(t, recorder.Body.Bytes(), respBody)
			},
		},
		{
			name: "Verify message got deleted",
			setupRequest: func() (*http.Request, error) {
				return http.NewRequest(http.MethodGet, "/api/messages/"+mockChannel.ID, nil)
			},
			setupHeaders: func(t *testing.T, request *http.Request) {
				request.Header.Add("Cookie", cookie)
				request.Header.Set("Content-Type", "application/json")
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, recorder.Code)

				respBody := &[]model.MessageResponse{}
				err := json.Unmarshal(recorder.Body.Bytes(), respBody)
				assert.NoError(t, err)
				assert.NotNil(t, respBody)

				messages := *respBody
				assert.Len(t, messages, 0)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			rr := httptest.NewRecorder()
			request, err := tc.setupRequest()
			tc.setupHeaders(t, request)
			assert.NoError(t, err)
			router.ServeHTTP(rr, request)
			tc.checkResponse(rr)
		})
	}
}

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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetMemberSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()
	mockGuild := fixture.GetMockGuild("")

	settings := &model.MemberSettings{
		Nickname: nil,
		Color:    nil,
	}

	t.Run("Successfully fetched settings", func(t *testing.T) {

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockArgs := mock.Arguments{
			authUser.ID,
			mockGuild.ID,
		}
		mockGuildService.On("GetMemberSettings", mockArgs...).Return(settings, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(settings)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "GetMemberSettings")
	})

	t.Run("Guild not found", func(t *testing.T) {
		id := fixture.RandID()

		mockGuildService := new(mocks.GuildService)

		mockError := apperrors.NewNotFound("guild", id)
		mockGuildService.On("GetGuild", id).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", id)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", id)
		mockGuildService.AssertNotCalled(t, "GetMemberSettings")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockError := apperrors.NewNotFound("user", authUser.ID)
		mockArgs := mock.Arguments{
			authUser.ID,
			mockGuild.ID,
		}
		mockGuildService.On("GetMemberSettings", mockArgs...).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})
}

func TestHandler_EditMemberSettings(t *testing.T) {
	gin.SetMode(gin.TestMode)

	authUser := fixture.GetMockUser()
	nickname := fixture.Username()
	color := "#fff"

	settings := &model.MemberSettings{
		Nickname: &nickname,
		Color:    &color,
	}

	t.Run("Successfully edited settings", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockArgs := mock.Arguments{
			settings,
			authUser.ID,
			mockGuild.ID,
		}
		mockGuildService.On("UpdateMemberSettings", mockArgs...).Return(nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"nickname": nickname,
			"color":    color,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
	})

	t.Run("Successfully reset member settings", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockArgs := mock.Arguments{
			&model.MemberSettings{},
			authUser.ID,
			mockGuild.ID,
		}
		mockGuildService.On("UpdateMemberSettings", mockArgs...).Return(nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"nickname": nil,
			"color":    nil,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
	})

	t.Run("Not a member of the server", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"nickname": nickname,
			"color":    color,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewNotFound("guild", mockGuild.ID)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockGuildService.AssertNotCalled(t, "UpdateMemberSettings")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuildService := new(mocks.GuildService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPut, reqUrl, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "UpdateMemberSettings")
	})

	t.Run("Guild not found", func(t *testing.T) {
		id := fixture.RandID()

		mockGuildService := new(mocks.GuildService)

		mockError := apperrors.NewNotFound("guild", id)
		mockGuildService.On("GetGuild", id).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"nickname": nickname,
			"color":    color,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", id)
		request, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", id)
		mockGuildService.AssertNotCalled(t, "UpdateMemberSettings")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockError := apperrors.NewInternal()
		mockArgs := mock.Arguments{
			settings,
			authUser.ID,
			mockGuild.ID,
		}
		mockGuildService.On("UpdateMemberSettings", mockArgs...).Return(mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"nickname": nickname,
			"color":    color,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/member", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})
}

func TestHandler_EditMemberSettings_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	router := getAuthenticatedTestRouter(mockUser.ID)

	mockGuildService := new(mocks.GuildService)

	NewHandler(&Config{
		R:            router,
		GuildService: mockGuildService,
	})

	testCases := []struct {
		name string
		body gin.H
	}{
		{
			name: "Nickname too short",
			body: gin.H{
				"nickname": fixture.RandStringRunes(2),
			},
		},
		{
			name: "Nickname too long",
			body: gin.H{
				"nickname": fixture.RandStringRunes(32),
			},
		},
		{
			name: "Color not a hex color",
			body: gin.H{
				"color": fixture.RandStringRunes(6),
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			reqUrl := fmt.Sprintf("/api/guilds/%s/member", fixture.RandID())
			request, err := http.NewRequest(http.MethodPut, reqUrl, bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockGuildService.AssertNotCalled(t, "UpdateMemberSettings")
		})
	}
}

func TestHandler_GetBanList(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful Fetch", func(t *testing.T) {
		response := make([]model.BanResponse, 0)
		mockGuild := fixture.GetMockGuild(authUser.ID)

		for i := 0; i < 5; i++ {
			mockUser := fixture.GetMockUser()
			response = append(response, model.BanResponse{
				Id:       mockUser.ID,
				Username: mockUser.Username,
				Image:    mockUser.Image,
			})
		}

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("GetBanList", mockGuild.ID).Return(&response, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Not the owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
		mockGuildService.AssertNotCalled(t, "GetBanList")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "GetBanList")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockError := apperrors.NewInternal()
		mockGuildService.On("GetBanList", mockGuild.ID).Return(nil, mockError)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})
}

func TestHandler_BanMember(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful Ban", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("GetUser", mockMember.ID).Return(mockMember, nil)
		mockGuildService.On("UpdateGuild", mockGuild).Return(nil)

		args := mock.Arguments{
			mockMember.ID,
			mockGuild.ID,
		}
		mockGuildService.On("RemoveMember", args...).Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitRemoveMember", mockGuild.ID, mockMember.ID)
		mockSocketService.On("EmitRemoveFromGuild", mockMember.ID, mockGuild.ID)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Not the owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("GetUser", mockMember.ID).Return(mockMember, nil)
		mockGuildService.On("UpdateGuild", mockGuild).Return(nil)

		mockError := apperrors.NewInternal()
		args := mock.Arguments{
			mockMember.ID,
			mockGuild.ID,
		}
		mockGuildService.On("RemoveMember", args...).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)

		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Guild not found", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockError := apperrors.NewNotFound("guild", mockGuild.ID)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Member not found", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockError := apperrors.NewNotFound("user", mockMember.ID)
		mockGuildService.On("GetUser", mockMember.ID).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertCalled(t, "GetUser", mockMember.ID)
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("MemberId required", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("MemberId and AuthUserId are equal", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": authUser.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.BanYourselfError)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertCalled(t, "GetUser", authUser.ID)
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})
}

func TestHandler_KickMember(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful Kick", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("GetUser", mockMember.ID).Return(mockMember, nil)

		args := mock.Arguments{
			mockMember.ID,
			mockGuild.ID,
		}
		mockGuildService.On("RemoveMember", args...).Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitRemoveMember", mockGuild.ID, mockMember.ID)
		mockSocketService.On("EmitRemoveFromGuild", mockMember.ID, mockGuild.ID)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Not the owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("GetUser", mockMember.ID).Return(mockMember, nil)

		mockError := apperrors.NewInternal()
		args := mock.Arguments{
			mockMember.ID,
			mockGuild.ID,
		}
		mockGuildService.On("RemoveMember", args...).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)

		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Guild not found", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockError := apperrors.NewNotFound("guild", mockGuild.ID)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("Member not found", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockError := apperrors.NewNotFound("user", mockMember.ID)
		mockGuildService.On("GetUser", mockMember.ID).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertCalled(t, "GetUser", mockMember.ID)
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("MemberId required", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})

	t.Run("MemberId and AuthUserId are equal", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": authUser.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/kick", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.KickYourselfError)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertCalled(t, "GetUser", authUser.ID)
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFromGuild")
	})
}

func TestHandler_UnbanMember(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful Unban", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		args := mock.Arguments{
			mockMember.ID,
			mockGuild.ID,
		}
		mockGuildService.On("UnbanMember", args...).Return(nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Not the owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UnbanMember")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "UnbanMember")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockError := apperrors.NewInternal()
		args := mock.Arguments{
			mockMember.ID,
			mockGuild.ID,
		}
		mockGuildService.On("UnbanMember", args...).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Guild not found", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockMember := fixture.GetMockUser()

		mockGuildService := new(mocks.GuildService)
		mockError := apperrors.NewNotFound("guild", mockGuild.ID)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(nil, mockError)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": mockMember.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UnbanMember")
	})

	t.Run("MemberId required", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "UnbanMember")
	})

	t.Run("MemberId and AuthUserId are equal", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqBody, err := json.Marshal(gin.H{
			"memberId": authUser.ID,
		})
		assert.NoError(t, err)

		reqUrl := fmt.Sprintf("/api/guilds/%s/bans", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.UnbanYourselfError)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UnbanMember")
	})
}

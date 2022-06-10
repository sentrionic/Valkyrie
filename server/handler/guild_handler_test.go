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
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestHandler_GetUserGuilds(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	response := make([]model.GuildResponse, 0)
	for i := 0; i < 5; i++ {
		mockGuild := fixture.GetMockGuild("")
		response = append(response, mockGuild.SerializeGuild(fixture.RandID()))
	}

	t.Run("Successful Fetch", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUserGuilds", authUser.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/guilds", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetUserGuilds", authUser.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/guilds", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		mockFriendService.AssertNotCalled(t, "GetUserGuilds", authUser.ID, "")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUserGuilds", authUser.ID).Return(nil, fmt.Errorf("some error down call chain"))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/guilds", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewNotFound("user", authUser.ID)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})
}

func TestHandler_GetGuildMembers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()
	guild := fixture.GetMockGuild("")
	guild.Members = append(guild.Members, *authUser)

	response := make([]model.MemberResponse, 0)
	for i := 0; i < 5; i++ {
		mockUser := fixture.GetMockUser()
		response = append(response, model.MemberResponse{
			Id:        mockUser.ID,
			Username:  mockUser.Username,
			Image:     mockUser.Image,
			IsOnline:  mockUser.IsOnline,
			CreatedAt: mockUser.CreatedAt,
			UpdatedAt: mockUser.UpdatedAt,
			IsFriend:  false,
		})
	}

	t.Run("Successful Fetch", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", guild.ID).Return(guild, nil)
		mockGuildService.On("GetGuildMembers", authUser.ID, guild.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/members", guild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", guild.ID).Return(guild, nil)
		mockGuildService.On("GetGuildMembers", authUser.ID, guild.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/members", guild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		mockGuildService.AssertNotCalled(t, "GetGuild", guild.ID)
		mockGuildService.AssertNotCalled(t, "GetGuildMembers", authUser.ID, guild.ID)
	})

	t.Run("Not a member of the guild", func(t *testing.T) {

		invalidGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", invalidGuild.ID).Return(invalidGuild, nil)
		mockGuildService.On("GetGuildMembers", authUser.ID, invalidGuild.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/members", invalidGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		mockError := apperrors.NewAuthorization(apperrors.NotAMember)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", invalidGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetGuildMembers", authUser.ID, invalidGuild.ID)
	})

	t.Run("Guild not found", func(t *testing.T) {

		invalidGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)

		mockError := apperrors.NewNotFound("guild", invalidGuild.ID)
		mockGuildService.On("GetGuild", invalidGuild.ID).Return(nil, mockError)
		mockGuildService.On("GetGuildMembers", authUser.ID, invalidGuild.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/members", invalidGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", invalidGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetGuildMembers", authUser.ID, invalidGuild.ID)
	})

	t.Run("Error", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", guild.ID).Return(guild, nil)

		mockError := apperrors.NewNotFound("user", authUser.ID)
		mockGuildService.On("GetGuildMembers", authUser.ID, guild.ID).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/members", guild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", guild.ID)
		mockGuildService.AssertCalled(t, "GetGuildMembers", authUser.ID, guild.ID)
	})
}

func TestHandler_GetVCMembers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()
	guild := fixture.GetMockGuild("")
	guild.Members = append(guild.Members, *authUser)

	response := make([]model.VCMemberResponse, 0)

	for i := 0; i < 5; i++ {
		mockUser := fixture.GetMockUser()

		guild.VCMembers = append(guild.VCMembers, *authUser)
		response = append(response, model.VCMemberResponse{
			Id:         mockUser.ID,
			Username:   mockUser.Username,
			Image:      mockUser.Image,
			IsMuted:    false,
			IsDeafened: false,
		})
	}

	t.Run("Successful Fetch", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", guild.ID).Return(guild, nil)
		mockGuildService.On("GetVCMembers", guild.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/vcmembers", guild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", guild.ID).Return(guild, nil)
		mockGuildService.On("GetVCMembers", guild.ID).Return(response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/vcmembers", guild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		mockGuildService.AssertNotCalled(t, "GetGuild", guild.ID)
		mockGuildService.AssertNotCalled(t, "GetVCMembers", guild.ID)
	})

	t.Run("Not a member of the guild", func(t *testing.T) {
		invalidGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", invalidGuild.ID).Return(invalidGuild, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/vcmembers", invalidGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		mockError := apperrors.NewAuthorization(apperrors.NotAMember)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", invalidGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetVCMembers", guild.ID)
	})

	t.Run("Guild not found", func(t *testing.T) {

		invalidGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)

		mockError := apperrors.NewNotFound("guild", invalidGuild.ID)
		mockGuildService.On("GetGuild", invalidGuild.ID).Return(nil, mockError)
		mockGuildService.On("GetVCMembers", guild.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/vcmembers", invalidGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", invalidGuild.ID)
		mockGuildService.AssertNotCalled(t, "GetVCMembers", guild.ID)
	})

	t.Run("Error", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)

		mockGuildService.On("GetGuild", guild.ID).Return(guild, nil)
		mockError := apperrors.NewNotFound("guild", guild.ID)
		mockGuildService.On("GetVCMembers", guild.ID).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/vcmembers", guild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", guild.ID)
		mockGuildService.AssertCalled(t, "GetVCMembers", guild.ID)
	})
}

func TestHandler_CreateGuild(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	authUser := fixture.GetMockUser()

	t.Run("Successful guild creation", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		router := getAuthenticatedTestRouter(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		reqBody, err := json.Marshal(gin.H{
			"name": mockGuild.Name,
		})
		assert.NoError(t, err)

		guildParams := &model.Guild{
			Name:    mockGuild.Name,
			OwnerId: authUser.ID,
		}
		guildParams.Members = append(guildParams.Members, *authUser)

		mockGuildService.On("CreateGuild", guildParams).Return(mockGuild, nil)

		defaultChannel := fixture.GetMockChannel(mockGuild.ID)
		defaultChannel.Name = "general"
		mockGuild.Channels = append(mockGuild.Channels, *defaultChannel)

		mockChannelService := new(mocks.ChannelService)
		channelParams := &model.Channel{
			GuildID:  &mockGuild.ID,
			Name:     defaultChannel.Name,
			IsPublic: true,
		}

		mockChannelService.On("CreateChannel", channelParams).Return(defaultChannel, nil)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(mockGuild.SerializeGuild(defaultChannel.ID))
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetUser", authUser.ID)
		mockGuildService.AssertCalled(t, "CreateGuild", guildParams)
		mockChannelService.AssertCalled(t, "CreateChannel", channelParams)
	})

	t.Run("Error Returned from GuildService.CreateGuild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		router := getAuthenticatedTestRouter(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		reqBody, err := json.Marshal(gin.H{
			"name": mockGuild.Name,
		})
		assert.NoError(t, err)

		guildParams := &model.Guild{
			Name:    mockGuild.Name,
			OwnerId: authUser.ID,
		}
		guildParams.Members = append(guildParams.Members, *authUser)

		mockError := apperrors.NewInternal()
		mockGuildService.On("CreateGuild", guildParams).Return(nil, mockError)

		mockChannelService := new(mocks.ChannelService)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetUser", authUser.ID)
		mockGuildService.AssertCalled(t, "CreateGuild", guildParams)
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("User already is in the maximum number of guilds", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		for i := 0; i < model.MaximumGuilds; i++ {
			guild := fixture.GetMockGuild("")
			authUser.Guilds = append(authUser.Guilds, *guild)
		}

		router := getAuthenticatedTestRouter(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		reqBody, err := json.Marshal(gin.H{
			"name": mockGuild.Name,
		})
		assert.NoError(t, err)

		mockError := apperrors.NewBadRequest(apperrors.GuildLimitReached)
		mockChannelService := new(mocks.ChannelService)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetUser", authUser.ID)
		mockGuildService.AssertNotCalled(t, "CreateGuild")
		mockChannelService.AssertNotCalled(t, "CreateChannel")

	})

	t.Run("Unauthorized", func(t *testing.T) {
		router := getTestRouter()

		mockGuildService := new(mocks.GuildService)
		mockChannelService := new(mocks.ChannelService)

		reqBody, err := json.Marshal(gin.H{
			"name": fixture.RandStringRunes(6),
		})
		assert.NoError(t, err)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "CreateGuild")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})
}

func TestHandler_CreateGuild_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	router := getAuthenticatedTestRouter(mockUser.ID)

	mockGuildService := new(mocks.GuildService)
	mockGuildService.On("GetUser", mockUser.ID).Return(mockUser, nil)

	NewHandler(&Config{
		R:            router,
		GuildService: mockGuildService,
	})

	testCases := []struct {
		name string
		body gin.H
	}{
		{
			name: "Name required",
			body: gin.H{},
		},
		{
			name: "Name too short",
			body: gin.H{
				"name": fixture.RandStr(2),
			},
		},
		{
			name: "Name too long",
			body: gin.H{
				"name": fixture.RandStr(31),
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			reqBody, err := json.Marshal(tc.body)
			assert.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, "/api/guilds/create", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockGuildService.AssertNotCalled(t, "CreateGuild")
		})
	}
}

func TestHandler_UpdateGuild(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	authUser := fixture.GetMockUser()

	t.Run("Successfully updated guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		router := getAuthenticatedTestRouter(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		name := fixture.RandStringRunes(8)
		form := url.Values{}
		form.Add("name", name)

		mockGuildService.On("UpdateGuild", mockGuild).
			Run(func(args mock.Arguments) {
				mockGuild.Name = name
			}).
			Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitEditGuild", mockGuild).Return()

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPut, "/api/guilds/"+mockGuild.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertCalled(t, "UpdateGuild", mockGuild)
		mockSocketService.AssertCalled(t, "EmitEditGuild", mockGuild)
	})

	t.Run("Error Returned from GuildService.UpdateGuild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		router := getAuthenticatedTestRouter(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		name := fixture.RandStringRunes(8)
		form := url.Values{}
		form.Add("name", name)

		mockError := apperrors.NewInternal()
		mockGuildService.On("UpdateGuild", mockGuild).Return(mockError)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitEditGuild", mockGuild).Return()

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPut, "/api/guilds/"+mockGuild.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertCalled(t, "UpdateGuild", mockGuild)
		mockSocketService.AssertNotCalled(t, "EmitEditGuild", mockGuild)
	})

	t.Run("Not the owner of the guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		router := getAuthenticatedTestRouter(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		name := fixture.RandStringRunes(8)
		form := url.Values{}
		form.Add("name", name)

		mockSocketService := new(mocks.SocketService)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPut, "/api/guilds/"+mockGuild.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		mockError := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UpdateGuild", mockGuild)
		mockSocketService.AssertNotCalled(t, "EmitEditGuild", mockGuild)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		router := getTestRouter()

		mockGuildService := new(mocks.GuildService)

		name := fixture.RandStringRunes(8)
		form := url.Values{}
		form.Add("name", name)

		mockSocketService := new(mocks.SocketService)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPut, "/api/guilds/"+mockGuild.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UpdateGuild", mockGuild)
		mockSocketService.AssertNotCalled(t, "EmitEditGuild", mockGuild)
	})

	t.Run("Guild not found", func(t *testing.T) {
		id := fixture.RandID()
		router := getAuthenticatedTestRouter(authUser.ID)

		mockError := apperrors.NewNotFound("guild", id)
		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", id).Return(nil, mockError)

		name := fixture.RandStringRunes(8)
		form := url.Values{}
		form.Add("name", name)

		mockSocketService := new(mocks.SocketService)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodPut, "/api/guilds/"+id, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", id)
		mockGuildService.AssertNotCalled(t, "UpdateGuild", mock.AnythingOfType("*model.Guild"))
		mockSocketService.AssertNotCalled(t, "EmitEditGuild", mock.AnythingOfType("*model.Guild"))
	})
}

func TestHandler_UpdateGuild_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	router := getAuthenticatedTestRouter(mockUser.ID)

	mockGuildService := new(mocks.GuildService)
	mockGuildService.On("GetUser", mockUser.ID).Return(mockUser, nil)

	NewHandler(&Config{
		R:            router,
		GuildService: mockGuildService,
	})

	testCases := []struct {
		name string
		body url.Values
	}{
		{
			name: "Name required",
			body: map[string][]string{},
		},
		{
			name: "Name too short",
			body: map[string][]string{
				"name": {fixture.RandStr(2)},
			},
		},
		{
			name: "Name too long",
			body: map[string][]string{
				"name": {fixture.RandStr(31)},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			reqUrl := fmt.Sprintf("/api/guilds/%s", fixture.RandID())
			form := tc.body
			request, _ := http.NewRequest(http.MethodPut, reqUrl, strings.NewReader(form.Encode()))
			request.Form = form

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockGuildService.AssertNotCalled(t, "GetGuild")
		})
	}
}

func TestHandler_GetInvite(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	origin := "http://localhost:3000"

	err := os.Setenv("CORS_ORIGIN", origin)
	assert.NoError(t, err)

	t.Run("Successful Fetch", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		getInviteArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockGuild.ID,
			false,
		}

		link := fixture.RandID()

		mockGuildService.On("GenerateInviteLink", getInviteArgs...).Return(link, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(fmt.Sprintf("%s/%s", origin, link))
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
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

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})

	t.Run("Not a member of the guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.MustBeMemberInvite)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "GenerateInviteLink")
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

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", id)
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
		mockGuildService.AssertNotCalled(t, "GenerateInviteLink")
	})

	t.Run("Invalid isPermanent value", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite?isPermanent=yes", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.IsPermanentError)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "GenerateInviteLink")
	})

	t.Run("Invite isPermanent success", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)

		link := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("UpdateGuild", mockGuild).
			Run(func(args mock.Arguments) {
				mockGuild.InviteLinks = append(mockGuild.InviteLinks, link)
			}).
			Return(nil)

		getInviteArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockGuild.ID,
			true,
		}

		mockGuildService.On("GenerateInviteLink", getInviteArgs...).Return(link, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite?isPermanent=true", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(fmt.Sprintf("%s/%s", origin, link))
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		getInviteArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockGuild.ID,
			false,
		}

		mockError := apperrors.NewInternal()
		mockGuildService.On("GenerateInviteLink", getInviteArgs...).Return("", mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
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

func TestHandler_DeleteGuildInvites(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully deleted", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockGuild,
		}

		mockGuildService.On("InvalidateInvites", mockArgs...).Return()
		mockGuildService.On("UpdateGuild", mockGuild).Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
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

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			mockGuild,
		}

		mockGuildService.On("InvalidateInvites", mockArgs...).Return()

		mockError := apperrors.NewInternal()
		mockGuildService.On("UpdateGuild", mockGuild).Return(mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
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

	t.Run("Not the server owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidateInvitesError)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "InvalidateInvites")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
	})

	t.Run("Guild not found", func(t *testing.T) {
		id := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockError := apperrors.NewNotFound("guild", id)
		mockGuildService.On("GetGuild", id).Return(nil, mockError)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:            router,
			GuildService: mockGuildService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/invite", id)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", id)
		mockGuildService.AssertNotCalled(t, "InvalidateInvites")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
	})
}

func TestHandler_JoinGuild(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully joined", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		link := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			link,
		}

		mockGuildService.On("GetGuildIdFromInvite", mockArgs...).Return(mockGuild.ID, nil)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("UpdateGuild", mockGuild).Return(nil)

		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockGuildService.On("GetDefaultChannel", mockGuild.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitAddMember", mockGuild.ID, authUser).Return()

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"link": link,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(mockGuild.SerializeGuild(mockChannel.ID))
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Link required", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"link": "",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "GetGuildIdFromInvite")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "GetDefaultChannel")
		mockSocketService.AssertNotCalled(t, "EmitAddMember")
	})

	t.Run("Unauthorized", func(t *testing.T) {
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
			"link": "link",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetUser")
		mockGuildService.AssertNotCalled(t, "GetGuildIdFromInvite")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "GetDefaultChannel")
		mockSocketService.AssertNotCalled(t, "EmitAddMember")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		link := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			link,
		}

		mockGuildService.On("GetGuildIdFromInvite", mockArgs...).Return(mockGuild.ID, nil)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockError := apperrors.NewInternal()
		mockGuildService.On("UpdateGuild", mockGuild).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"link": link,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
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
		mockGuildService.AssertNotCalled(t, "GetDefaultChannel")
		mockSocketService.AssertNotCalled(t, "EmitAddMember")
	})

	t.Run("User is already in the maximum number of guilds", func(t *testing.T) {
		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		newAuthUser := fixture.GetMockUser()
		for i := 0; i < model.MaximumGuilds; i++ {
			mockGuild := fixture.GetMockGuild("")
			newAuthUser.Guilds = append(newAuthUser.Guilds, *mockGuild)
		}

		mockGuildService.On("GetUser", newAuthUser.ID).Return(newAuthUser, nil)

		router := getAuthenticatedTestRouter(newAuthUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"link": "link",
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.GuildLimitReached)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetUser", newAuthUser.ID)
		mockGuildService.AssertNotCalled(t, "GetGuildIdFromInvite")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockGuildService.AssertNotCalled(t, "GetDefaultChannel")
		mockSocketService.AssertNotCalled(t, "EmitAddMember")
	})

	t.Run("User is banned from the guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Bans = append(mockGuild.Bans, *authUser)
		link := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			link,
		}

		mockGuildService.On("GetGuildIdFromInvite", mockArgs...).Return(mockGuild.ID, nil)
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
			"link": link,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.BannedFromServer)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockGuildService.AssertNotCalled(t, "GetDefaultChannel")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockSocketService.AssertNotCalled(t, "EmitAddMember")
	})

	t.Run("Invalid Invite", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		link := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			link,
		}

		mockError := apperrors.NewBadRequest(apperrors.InvalidInviteError)
		mockGuildService.On("GetGuildIdFromInvite", mockArgs...).Return(mockGuild.ID, mockError)

		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqBody, err := json.Marshal(gin.H{
			"link": link,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
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
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "GetDefaultChannel")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockSocketService.AssertNotCalled(t, "EmitAddMember")
	})

	t.Run("Already a member", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *authUser)
		link := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetUser", authUser.ID).Return(authUser, nil)

		mockArgs := mock.Arguments{
			mock.AnythingOfType("*context.emptyCtx"),
			link,
		}

		mockGuildService.On("GetGuildIdFromInvite", mockArgs...).Return(mockGuild.ID, nil)
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
			"link": link,
		})
		assert.NoError(t, err)

		request, err := http.NewRequest(http.MethodPost, "/api/guilds/join", bytes.NewBuffer(reqBody))
		assert.NoError(t, err)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.AlreadyMember)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockGuildService.AssertNotCalled(t, "GetDefaultChannel")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockSocketService.AssertNotCalled(t, "EmitAddMember")
	})
}

func TestHandler_LeaveGuild(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully left the guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("RemoveMember", authUser.ID, mockGuild.ID).Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitRemoveMember", mockGuild.ID, authUser.ID).Return()

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("AuthUser is the owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.OwnerCantLeave)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
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
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
	})

	t.Run("Guild not found", func(t *testing.T) {
		id := fixture.RandID()

		mockGuildService := new(mocks.GuildService)

		mockError := apperrors.NewNotFound("guild", id)
		mockGuildService.On("GetGuild", id).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s", id)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", id)
		mockGuildService.AssertNotCalled(t, "RemoveMember")
		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockError := apperrors.NewInternal()
		mockGuildService.On("RemoveMember", authUser.ID, mockGuild.ID).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)

		mockSocketService.AssertNotCalled(t, "EmitRemoveMember")
	})
}

func TestHandler_DeleteGuild(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully deleted", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("DeleteGuild", mockGuild.ID).Return(nil)

		mockSocketService := new(mocks.SocketService)

		members := make([]string, 0)
		mockSocketService.On("EmitDeleteGuild", mockGuild.ID, members).Return()

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/delete", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Not the guild owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)

		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/delete", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.DeleteGuildError)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "DeleteGuild")
		mockSocketService.AssertNotCalled(t, "EmitDeleteGuild")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/delete", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
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
		mockGuildService.AssertNotCalled(t, "DeleteGuild")
		mockSocketService.AssertNotCalled(t, "EmitDeleteGuild")
	})

	t.Run("Guild not found", func(t *testing.T) {
		id := fixture.RandID()

		mockGuildService := new(mocks.GuildService)

		mockError := apperrors.NewNotFound("guild", id)
		mockGuildService.On("GetGuild", id).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/delete", id)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", id)
		mockGuildService.AssertNotCalled(t, "DeleteGuild")
		mockSocketService.AssertNotCalled(t, "EmitDeleteGuild")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockError := apperrors.NewInternal()
		mockGuildService.On("DeleteGuild", mockGuild.ID).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			GuildService:  mockGuildService,
			SocketService: mockSocketService,
		})

		reqUrl := fmt.Sprintf("/api/guilds/%s/delete", mockGuild.ID)
		request, err := http.NewRequest(http.MethodDelete, reqUrl, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)

		mockSocketService.AssertNotCalled(t, "EmitDeleteGuild")
	})
}

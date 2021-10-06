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

func TestHandler_GuildChannels(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful Fetch", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockGuild.Members = append(mockGuild.Members, *authUser)

		response := make([]model.ChannelResponse, 0)
		for i := 0; i < 5; i++ {
			mockChannel := fixture.GetMockChannel(mockGuild.ID)
			response = append(response, mockChannel.SerializeChannel())
		}

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("GetChannels", authUser.ID, mockGuild.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Guild not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("guild", id)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", id).Return(nil, mockError)

		mockChannelService := new(mocks.ChannelService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s", id)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", id)
		mockChannelService.AssertNotCalled(t, "GetChannels")
	})

	t.Run("Not a member of the guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockError := apperrors.NewNotFound("guild", mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockChannelService.AssertNotCalled(t, "GetChannels")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		mockGuildService := new(mocks.GuildService)
		mockChannelService := new(mocks.ChannelService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockChannelService.AssertNotCalled(t, "GetChannels")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockGuild.Members = append(mockGuild.Members, *authUser)
		mockError := apperrors.NewNotFound("channels", mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("GetChannels", authUser.ID, mockGuild.ID).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
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

func TestHandler_CreateChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful channel creation", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("UpdateGuild", mockGuild).Return(nil)

		mockChannelService := new(mocks.ChannelService)

		channelParams := &model.Channel{
			Name:     mockChannel.Name,
			IsPublic: true,
			GuildID:  &mockGuild.ID,
		}
		mockChannelService.On("CreateChannel", channelParams).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)
		response := mockChannel.SerializeChannel()
		mockSocketService.On("EmitNewChannel", mockGuild.ID, &response)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(response)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Guild not found", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)

		mockError := apperrors.NewNotFound("guild", mockGuild.ID)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(nil, mockError)

		mockChannelService := new(mocks.ChannelService)
		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
		mockSocketService.AssertNotCalled(t, "EmitNewChannel")
	})

	t.Run("Guild already has the maximum number of channels", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		for i := 0; i < model.MaximumChannels; i++ {
			channel := fixture.GetMockChannel(mockGuild.ID)
			mockGuild.Channels = append(mockGuild.Channels, *channel)
		}

		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)

		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewBadRequest(apperrors.ChannelLimitError)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
		mockSocketService.AssertNotCalled(t, "EmitNewChannel")
	})

	t.Run("Not the guild owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
		mockSocketService.AssertNotCalled(t, "EmitNewChannel")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockChannelService := new(mocks.ChannelService)
		mockSocketService := new(mocks.SocketService)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
		mockSocketService.AssertNotCalled(t, "EmitNewChannel")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)

		channelParams := &model.Channel{
			Name:     mockChannel.Name,
			IsPublic: true,
			GuildID:  &mockGuild.ID,
		}
		mockError := apperrors.NewInternal()
		mockChannelService.On("CreateChannel", channelParams).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockGuildService.AssertNotCalled(t, "UpdateGuild")
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertNotCalled(t, "EmitNewChannel")
	})

	t.Run("Successful private channel creation", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		members := make([]model.User, 0)
		members = append(members, *authUser)

		reqMembers := []string{authUser.ID}

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)
		mockGuildService.On("UpdateGuild", mockGuild).Return(nil)
		mockGuildService.On("FindUsersByIds", reqMembers, mockGuild.ID).Return(&members, nil)

		mockChannelService := new(mocks.ChannelService)

		channelParams := &model.Channel{
			Name:     mockChannel.Name,
			IsPublic: false,
			GuildID:  &mockGuild.ID,
		}

		channelParams.PCMembers = append(channelParams.PCMembers, *authUser)
		mockChannelService.On("CreateChannel", channelParams).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)
		response := mockChannel.SerializeChannel()
		mockSocketService.On("EmitNewChannel", mockGuild.ID, &response)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name":     mockChannel.Name,
			"isPublic": false,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockGuild.ID)
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(response)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})
}

func TestHandler_CreateChannel_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	router := getAuthenticatedTestRouter(mockUser.ID)

	mockGuildService := new(mocks.GuildService)
	mockChannelService := new(mocks.ChannelService)

	NewHandler(&Config{
		R:              router,
		GuildService:   mockGuildService,
		ChannelService: mockChannelService,
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

			url := fmt.Sprintf("/api/channels/%s", fixture.RandID())
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockGuildService.AssertNotCalled(t, "GetGuild")
			mockChannelService.AssertNotCalled(t, "CreateChannel")
		})
	}
}

func TestHandler_PrivateChannelMembers(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful Fetch", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockChannel.IsPublic = false

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		response := make([]string, 0)
		for i := 0; i < 5; i++ {
			mockUser := fixture.GetMockUser()
			response = append(response, mockUser.ID)
		}

		mockChannelService.On("GetPrivateChannelMembers", mockChannel.ID).Return(&response, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s/members", mockChannel.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockGuildService.AssertExpectations(t)
	})

	t.Run("Guild not found", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel("")
		mockChannel.IsPublic = false
		mockError := apperrors.NewNotFound("channel", mockChannel.ID)

		mockGuildService := new(mocks.GuildService)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s/members", mockChannel.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockChannelService.AssertNotCalled(t, "GetPrivateChannelMembers")
	})

	t.Run("Channel not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("channel", id)

		mockGuildService := new(mocks.GuildService)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", id).Return(nil, mockError)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s/members", id)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", id)
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockChannelService.AssertNotCalled(t, "GetPrivateChannelMembers")
	})

	t.Run("Not the guild owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockChannel.IsPublic = false

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s/members", mockChannel.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, err := json.Marshal(gin.H{
			"error": e,
		})
		assert.NoError(t, err)
		router.ServeHTTP(rr, request)

		assert.Equal(t, e.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockChannelService.AssertNotCalled(t, "GetPrivateChannelMembers")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockChannel.IsPublic = false

		mockGuildService := new(mocks.GuildService)
		mockChannelService := new(mocks.ChannelService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s/members", mockChannel.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertNotCalled(t, "Get")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockChannelService.AssertNotCalled(t, "GetPrivateChannelMembers")
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockChannel.IsPublic = false

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockError := apperrors.NewNotFound("members", mockChannel.ID)
		mockChannelService.On("GetPrivateChannelMembers", mockChannel.ID).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
		})

		url := fmt.Sprintf("/api/channels/%s/members", mockChannel.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
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

func TestHandler_DirectMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Success", func(t *testing.T) {
		mockChannelService := new(mocks.ChannelService)

		response := make([]model.DirectMessage, 0)
		for i := 0; i < 5; i++ {
			user := fixture.GetMockUser()
			response = append(response, toDMChannel(user, fixture.RandID(), authUser.ID))
		}

		mockChannelService.On("GetDirectMessages", authUser.ID).Return(&response, nil)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/channels/me/dm", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(response)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockChannelService := new(mocks.ChannelService)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/channels/me/dm", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertNotCalled(t, "GetDirectMessages")
	})

	t.Run("Error", func(t *testing.T) {
		mockChannelService := new(mocks.ChannelService)

		mockError := apperrors.NewNotFound("dms", authUser.ID)
		mockChannelService.On("GetDirectMessages", authUser.ID).Return(nil, mockError)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		request, err := http.NewRequest(http.MethodGet, "/api/channels/me/dm", nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
	})
}

func TestHandler_GetOrCreateDM(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully returned an already existing DM", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		dmId := fixture.RandID()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("GetDirectMessageChannel", authUser.ID, mockUser.ID).Return(&dmId, nil)
		mockChannelService.On("SetDirectMessageStatus", dmId, authUser.ID, true).Return(nil)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			FriendService:  mockFriendService,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(toDMChannel(mockUser, dmId, authUser.ID))
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("Successfully returned a new DM", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("GetDirectMessageChannel", authUser.ID, mockUser.ID).Return(nil, nil)

		mockDM := fixture.GetMockDMChannel()
		channelParams := &model.Channel{
			Name:     fmt.Sprintf("%s-%s", authUser.ID, mockUser.ID),
			IsPublic: false,
			IsDM:     true,
		}
		mockChannelService.On("CreateChannel", channelParams).Return(mockDM, nil)

		ids := []string{authUser.ID, mockUser.ID}
		mockArgs := mock.Arguments{
			ids,
			mockDM.ID,
			authUser.ID,
		}
		mockChannelService.On("AddDMChannelMembers", mockArgs...).Return(nil)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			FriendService:  mockFriendService,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(toDMChannel(mockUser, mockDM.ID, authUser.ID))
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
	})

	t.Run("Member not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("member", id)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", id).Return(nil, mockError)

		mockChannelService := new(mocks.ChannelService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			FriendService:  mockFriendService,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", id)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertCalled(t, "GetMemberById", id)
		mockChannelService.AssertNotCalled(t, "GetDirectMessageChannel")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("Member and AuthUser are the same", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockChannelService := new(mocks.ChannelService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			FriendService:  mockFriendService,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", authUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		e := apperrors.NewBadRequest(apperrors.DMYourselfError)
		respBody, _ := json.Marshal(gin.H{
			"error": e,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, e.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertNotCalled(t, "GetMemberById")
		mockChannelService.AssertNotCalled(t, "GetDirectMessageChannel")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockChannelService := new(mocks.ChannelService)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			FriendService:  mockFriendService,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", fixture.RandID())
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertNotCalled(t, "GetMemberById")
		mockChannelService.AssertNotCalled(t, "GetDirectMessageChannel")
		mockChannelService.AssertNotCalled(t, "CreateChannel")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("GetDirectMessageChannel", authUser.ID, mockUser.ID).Return(nil, nil)

		mockDM := fixture.GetMockDMChannel()
		channelParams := &model.Channel{
			Name:     fmt.Sprintf("%s-%s", authUser.ID, mockUser.ID),
			IsPublic: false,
			IsDM:     true,
		}
		mockChannelService.On("CreateChannel", channelParams).Return(mockDM, nil)

		ids := []string{authUser.ID, mockUser.ID}
		mockArgs := mock.Arguments{
			ids,
			mockDM.ID,
			authUser.ID,
		}
		mockError := apperrors.NewInternal()
		mockChannelService.On("AddDMChannelMembers", mockArgs...).Return(mockError)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			FriendService:  mockFriendService,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
	})
}

func TestHandler_EditChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful edited channel", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockChannelService.On("UpdateChannel", mockChannel).Return(nil)

		mockSocketService := new(mocks.SocketService)
		response := mockChannel.SerializeChannel()
		mockSocketService.On("EmitEditChannel", mockGuild.ID, &response)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Guild not found", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockError := apperrors.NewNotFound("guild", mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(nil, mockError)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockChannelService.AssertNotCalled(t, "UpdateChannel")
		mockSocketService.AssertNotCalled(t, "EmitEditChannel")
	})

	t.Run("Channel not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("channel", id)

		mockGuildService := new(mocks.GuildService)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", id).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": fixture.RandStr(8),
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", id)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", id)
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockChannelService.AssertNotCalled(t, "UpdateChannel")
		mockSocketService.AssertNotCalled(t, "EmitEditChannel")
	})

	t.Run("Not the guild owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, _ := json.Marshal(gin.H{
			"error": e,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, e.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockChannelService.AssertNotCalled(t, "UpdateChannel")
		mockSocketService.AssertNotCalled(t, "EmitEditChannel")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockGuildService := new(mocks.GuildService)

		mockChannelService := new(mocks.ChannelService)

		mockSocketService := new(mocks.SocketService)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": fixture.RandStr(8),
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", id)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertNotCalled(t, "Get")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockChannelService.AssertNotCalled(t, "UpdateChannel")
		mockSocketService.AssertNotCalled(t, "EmitEditChannel")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockError := apperrors.NewInternal()
		mockChannelService.On("UpdateChannel", mockChannel).Return(mockError)

		mockSocketService := new(mocks.SocketService)
		response := mockChannel.SerializeChannel()
		mockSocketService.On("EmitEditChannel", mockGuild.ID, &response)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name": mockChannel.Name,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertNotCalled(t, "EmitEditChannel")
	})

	t.Run("Private channel made public", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockChannel.IsPublic = false

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		response := mockChannel.SerializeChannel()

		mockChannelService.On("CleanPCMembers", mockChannel.ID).Return(nil)
		mockChannelService.On("UpdateChannel", mockChannel).
			Run(func(args mock.Arguments) {
				mockChannel.IsPublic = true
				response = mockChannel.SerializeChannel()
			}).
			Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitEditChannel", mockGuild.ID, &response)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name":     mockChannel.Name,
			"isPublic": true,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Channel made private", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("AddPrivateChannelMembers", []string{authUser.ID}, mockChannel.ID).Return(nil)
		mockChannelService.On("RemovePrivateChannelMembers", []string(nil), mockChannel.ID).Return(nil)

		response := mockChannel.SerializeChannel()
		mockChannelService.On("UpdateChannel", mockChannel).
			Run(func(args mock.Arguments) {
				mockChannel.IsPublic = false
				response = mockChannel.SerializeChannel()
			}).
			Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitEditChannel", mockGuild.ID, &response)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		reqBody, err := json.Marshal(gin.H{
			"name":     mockChannel.Name,
			"isPublic": false,
		})
		assert.NoError(t, err)

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})
}

func TestHandler_EditChannel_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	router := getAuthenticatedTestRouter(mockUser.ID)

	mockGuildService := new(mocks.GuildService)
	mockChannelService := new(mocks.ChannelService)

	NewHandler(&Config{
		R:              router,
		GuildService:   mockGuildService,
		ChannelService: mockChannelService,
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

			url := fmt.Sprintf("/api/channels/%s", fixture.RandID())
			request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(reqBody))
			assert.NoError(t, err)

			request.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockGuildService.AssertNotCalled(t, "GetGuild")
			mockChannelService.AssertNotCalled(t, "CreateChannel")
		})
	}
}

func TestHandler_DeleteChannel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully deleted", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockGuild.Channels = append(mockGuild.Channels, *mockChannel)
		mockGuild.Channels = append(mockGuild.Channels, *fixture.GetMockChannel(mockGuild.ID))

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("DeleteChannel", mockChannel).Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitDeleteChannel", mockChannel)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Channel not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("channel", id)

		mockGuildService := new(mocks.GuildService)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", id).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", id)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", id)
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "DeleteChannel")
		mockSocketService.AssertNotCalled(t, "EmitDeleteChannel")
	})

	t.Run("Guild not found", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel(fixture.RandID())
		mockError := apperrors.NewNotFound("channel", mockChannel.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", *mockChannel.GuildID).Return(nil, mockError)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockGuildService.AssertCalled(t, "GetGuild", *mockChannel.GuildID)
		mockGuildService.AssertNotCalled(t, "DeleteChannel")
		mockSocketService.AssertNotCalled(t, "EmitDeleteChannel")
	})

	t.Run("Not the guild owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", *mockChannel.GuildID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		respBody, _ := json.Marshal(gin.H{
			"error": e,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, e.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockGuildService.AssertCalled(t, "GetGuild", *mockChannel.GuildID)
		mockGuildService.AssertNotCalled(t, "DeleteChannel")
		mockSocketService.AssertNotCalled(t, "EmitDeleteChannel")
	})

	t.Run("Channel is last channel of the guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockGuild.Channels = append(mockGuild.Channels, *mockChannel)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", *mockChannel.GuildID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewBadRequest(apperrors.OneChannelRequired)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockGuildService.AssertCalled(t, "GetGuild", *mockChannel.GuildID)
		mockGuildService.AssertNotCalled(t, "DeleteChannel")
		mockSocketService.AssertNotCalled(t, "EmitDeleteChannel")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockGuildService := new(mocks.GuildService)
		mockChannelService := new(mocks.ChannelService)
		mockSocketService := new(mocks.SocketService)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", id)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertNotCalled(t, "Get")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockGuildService.AssertNotCalled(t, "DeleteChannel")
		mockSocketService.AssertNotCalled(t, "EmitDeleteChannel")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockGuild.Channels = append(mockGuild.Channels, *mockChannel)
		mockGuild.Channels = append(mockGuild.Channels, *fixture.GetMockChannel(mockGuild.ID))

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockError := apperrors.NewInternal()
		mockChannelService.On("DeleteChannel", mockChannel).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s", mockChannel.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockGuildService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})
}

func TestHandler_CloseDM(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully close", func(t *testing.T) {
		channelId := fixture.RandID()
		dmId := fixture.RandID()

		mockChannelService := new(mocks.ChannelService)

		findArgs := mock.Arguments{
			authUser.ID,
			channelId,
		}
		mockChannelService.On("GetDMByUserAndChannel", findArgs...).Return(dmId, nil)

		setArgs := mock.Arguments{
			channelId,
			authUser.ID,
			false,
		}
		mockChannelService.On("SetDirectMessageStatus", setArgs...).Return(nil)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", channelId)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
	})

	t.Run("DM not found", func(t *testing.T) {
		channelId := fixture.RandID()
		mockChannelService := new(mocks.ChannelService)

		findArgs := mock.Arguments{
			authUser.ID,
			channelId,
		}
		mockError := apperrors.NewNotFound("dms", authUser.ID)
		mockChannelService.On("GetDMByUserAndChannel", findArgs...).Return("", mockError)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", channelId)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "GetDMByUserAndChannel", findArgs...)
		mockChannelService.AssertNotCalled(t, "SetDirectMessageStatus")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		channelId := fixture.RandID()
		mockChannelService := new(mocks.ChannelService)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", channelId)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertNotCalled(t, "GetDMByUserAndChannel")
		mockChannelService.AssertNotCalled(t, "SetDirectMessageStatus")
	})

	t.Run("Server Error", func(t *testing.T) {
		channelId := fixture.RandID()
		dmId := fixture.RandID()

		mockChannelService := new(mocks.ChannelService)

		findArgs := mock.Arguments{
			authUser.ID,
			channelId,
		}
		mockChannelService.On("GetDMByUserAndChannel", findArgs...).Return(dmId, nil)

		setArgs := mock.Arguments{
			channelId,
			authUser.ID,
			false,
		}
		mockError := apperrors.NewInternal()
		mockChannelService.On("SetDirectMessageStatus", setArgs...).Return(mockError)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
		})

		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		url := fmt.Sprintf("/api/channels/%s/dm", channelId)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
	})
}

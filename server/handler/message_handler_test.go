package handler

import (
	"bytes"
	"encoding/json"
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
	"strings"
	"testing"
	"time"
)

func TestHandler_GetMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successful fetch", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel(fixture.RandID())

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(nil)

		mockMessageService := new(mocks.MessageService)

		args := mock.Arguments{
			authUser.ID,
			mockChannel,
			"",
		}

		response := make([]model.MessageResponse, 0)

		for i := 0; i < 25; i++ {
			message := fixture.GetMockMessageResponse("", mockChannel.ID)
			response = append(response, *message)
		}

		mockMessageService.On("GetMessages", args...).Return(&response, nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/messages/"+mockChannel.ID, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockChannelService.AssertExpectations(t)
		mockMessageService.AssertExpectations(t)
	})

	t.Run("No channel found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("channel", id)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", id).Return(nil, mockError)

		mockMessageService := new(mocks.MessageService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/messages/"+id, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", id)
		mockChannelService.AssertNotCalled(t, "IsChannelMember")
		mockMessageService.AssertNotCalled(t, "GetMessages")
	})

	t.Run("Not a member of the channel", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel(fixture.RandID())

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockError := apperrors.NewAuthorization(apperrors.Unauthorized)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(mockError)

		mockMessageService := new(mocks.MessageService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/messages/"+mockChannel.ID, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
		mockMessageService.AssertNotCalled(t, "GetMessages")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockChannelService := new(mocks.ChannelService)
		mockMessageService := new(mocks.MessageService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/messages/"+id, nil)
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
		mockChannelService.AssertNotCalled(t, "IsChannelMember")
		mockMessageService.AssertNotCalled(t, "GetMessages")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel(fixture.RandID())

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(nil)

		mockMessageService := new(mocks.MessageService)

		args := mock.Arguments{
			authUser.ID,
			mockChannel,
			"",
		}
		mockError := apperrors.NewNotFound("messages", mockChannel.ID)
		mockMessageService.On("GetMessages", args...).Return(nil, mockError)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/messages/"+mockChannel.ID, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
		mockMessageService.AssertExpectations(t)
	})
}

func TestHandler_CreateMessage(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully created text message", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockMessage := fixture.GetMockMessage(authUser.ID, mockChannel.ID)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(nil)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", authUser.ID).Return(authUser, nil)

		params := model.Message{
			UserId:    mockMessage.UserId,
			ChannelId: mockMessage.ChannelId,
			Text:      mockMessage.Text,
		}
		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("CreateMessage", &params).Return(mockMessage, nil)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetMemberSettings", authUser.ID, mockGuild.ID).Return(&model.MemberSettings{}, nil)

		mockSocketService := new(mocks.SocketService)
		response := model.MessageResponse{
			Id:         mockMessage.ID,
			Text:       mockMessage.Text,
			CreatedAt:  mockMessage.CreatedAt,
			UpdatedAt:  mockMessage.UpdatedAt,
			Attachment: mockMessage.Attachment,
			User: model.MemberResponse{
				Id:        authUser.ID,
				Username:  authUser.Username,
				Image:     authUser.Image,
				IsOnline:  authUser.IsOnline,
				CreatedAt: authUser.CreatedAt,
				UpdatedAt: authUser.UpdatedAt,
				IsFriend:  false,
			},
		}

		mockSocketService.On("EmitNewMessage", mockChannel.ID, &response).Return()
		mockChannelService.On("UpdateChannel", mockChannel).Return(nil)
		mockSocketService.On("EmitNewNotification", mockGuild.ID, mockChannel.ID)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		form := url.Values{}
		form.Add("text", *mockMessage.Text)

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+mockChannel.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
		mockMessageService.AssertExpectations(t)
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})

	t.Run("Channel not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("channel", id)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", id).Return(nil, mockError)

		mockUserService := new(mocks.UserService)
		mockMessageService := new(mocks.MessageService)
		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		form := url.Values{}
		form.Add("text", fixture.RandStringRunes(8))

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+id, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", id)
		mockChannelService.AssertNotCalled(t, "IsChannelMember")
		mockUserService.AssertNotCalled(t, "Get")
		mockMessageService.AssertNotCalled(t, "CreateMessage")
		mockSocketService.AssertNotCalled(t, "EmitNewMessage")
	})

	t.Run("Not a member of the channel", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel("")
		mockError := apperrors.NewAuthorization(apperrors.Unauthorized)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(mockError)

		mockUserService := new(mocks.UserService)
		mockMessageService := new(mocks.MessageService)
		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		form := url.Values{}
		form.Add("text", fixture.RandStringRunes(8))

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+mockChannel.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockChannelService.AssertCalled(t, "IsChannelMember", mockChannel, authUser.ID)
		mockUserService.AssertNotCalled(t, "Get")
		mockMessageService.AssertNotCalled(t, "CreateMessage")
		mockSocketService.AssertNotCalled(t, "EmitNewMessage")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockChannelService := new(mocks.ChannelService)
		mockUserService := new(mocks.UserService)
		mockMessageService := new(mocks.MessageService)
		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		form := url.Values{}
		form.Add("text", fixture.RandStringRunes(8))

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+id, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewAuthorization(apperrors.InvalidSession)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertNotCalled(t, "Get")
		mockChannelService.AssertNotCalled(t, "IsChannelMember")
		mockUserService.AssertNotCalled(t, "Get")
		mockMessageService.AssertNotCalled(t, "CreateMessage")
		mockSocketService.AssertNotCalled(t, "EmitNewMessage")
	})

	t.Run("Text Message Creation failure", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockMessage := fixture.GetMockMessage(authUser.ID, mockChannel.ID)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(nil)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", authUser.ID).Return(authUser, nil)

		params := model.Message{
			UserId:    mockMessage.UserId,
			ChannelId: mockMessage.ChannelId,
			Text:      mockMessage.Text,
		}
		mockError := apperrors.NewInternal()
		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("CreateMessage", &params).Return(nil, mockError)

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		form := url.Values{}
		form.Add("text", *mockMessage.Text)

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+mockChannel.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockChannelService.AssertCalled(t, "IsChannelMember", mockChannel, authUser.ID)
		mockMessageService.AssertCalled(t, "CreateMessage", &params)
		mockUserService.AssertCalled(t, "Get", authUser.ID)
		mockChannelService.AssertNotCalled(t, "UpdateChannel")
		mockSocketService.AssertNotCalled(t, "EmitNewMessage")
	})

	t.Run("Disallowed mimetype", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(nil)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", authUser.ID).Return(authUser, nil)

		mockMessageService := new(mocks.MessageService)
		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		multipartImageFixture := fixture.NewMultipartImage("image.txt", "image/txt")
		defer multipartImageFixture.Close()

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+mockChannel.ID, multipartImageFixture.MultipartBody)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockChannelService.AssertCalled(t, "IsChannelMember", mockChannel, authUser.ID)
		mockUserService.AssertCalled(t, "Get", authUser.ID)
		mockMessageService.AssertNotCalled(t, "UploadFile")
		mockMessageService.AssertNotCalled(t, "CreateMessage")
		mockChannelService.AssertNotCalled(t, "UpdateChannel")
		mockSocketService.AssertNotCalled(t, "EmitNewMessage")
	})

	t.Run("Image Message Creation Success", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockMessage := fixture.GetMockMessage(authUser.ID, mockChannel.ID)
		mockMessage.Text = nil

		uploadImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer uploadImageFixture.Close()
		formFile := uploadImageFixture.GetFormFile()

		attachment := &model.Attachment{
			ID:        fixture.RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Url:       fixture.RandStringRunes(8),
			FileType:  "image/png",
			Filename:  fixture.RandStringRunes(8),
			MessageId: mockMessage.ID,
		}
		mockMessage.Attachment = attachment

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(nil)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", authUser.ID).Return(authUser, nil)

		params := model.Message{
			UserId:     mockMessage.UserId,
			ChannelId:  mockMessage.ChannelId,
			Attachment: attachment,
		}
		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("UploadFile", formFile, mockChannel.ID).Return(attachment, nil)
		mockMessageService.On("CreateMessage", &params).Return(mockMessage, nil)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetMemberSettings", authUser.ID, mockGuild.ID).Return(&model.MemberSettings{}, nil)

		mockSocketService := new(mocks.SocketService)
		response := model.MessageResponse{
			Id:         mockMessage.ID,
			Text:       nil,
			CreatedAt:  mockMessage.CreatedAt,
			UpdatedAt:  mockMessage.UpdatedAt,
			Attachment: attachment,
			User: model.MemberResponse{
				Id:        authUser.ID,
				Username:  authUser.Username,
				Image:     authUser.Image,
				IsOnline:  authUser.IsOnline,
				CreatedAt: authUser.CreatedAt,
				UpdatedAt: authUser.UpdatedAt,
				IsFriend:  false,
			},
		}

		mockSocketService.On("EmitNewMessage", mockChannel.ID, &response).Return()
		mockChannelService.On("UpdateChannel", mockChannel).Return(nil)
		mockSocketService.On("EmitNewNotification", mockGuild.ID, mockChannel.ID)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+mockChannel.ID, multipartImageFixture.MultipartBody)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", multipartImageFixture.ContentType)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
		mockMessageService.AssertExpectations(t)
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})

	t.Run("DM channel message success", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel("")
		mockChannel.IsDM = true
		mockMessage := fixture.GetMockMessage(authUser.ID, mockChannel.ID)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)
		mockChannelService.On("IsChannelMember", mockChannel, authUser.ID).Return(nil)

		mockUserService := new(mocks.UserService)
		mockUserService.On("Get", authUser.ID).Return(authUser, nil)

		params := model.Message{
			UserId:    mockMessage.UserId,
			ChannelId: mockMessage.ChannelId,
			Text:      mockMessage.Text,
		}
		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("CreateMessage", &params).Return(mockMessage, nil)

		mockSocketService := new(mocks.SocketService)
		response := model.MessageResponse{
			Id:         mockMessage.ID,
			Text:       mockMessage.Text,
			CreatedAt:  mockMessage.CreatedAt,
			UpdatedAt:  mockMessage.UpdatedAt,
			Attachment: mockMessage.Attachment,
			User: model.MemberResponse{
				Id:        authUser.ID,
				Username:  authUser.Username,
				Image:     authUser.Image,
				IsOnline:  authUser.IsOnline,
				CreatedAt: authUser.CreatedAt,
				UpdatedAt: authUser.UpdatedAt,
				IsFriend:  false,
			},
		}

		mockSocketService.On("EmitNewMessage", mockChannel.ID, &response).Return()
		mockSocketService.On("EmitNewDMNotification", mockChannel.ID, authUser).Return()
		mockChannelService.On("OpenDMForAll", mockChannel.ID).Return(nil)

		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			ChannelService: mockChannelService,
			MessageService: mockMessageService,
			SocketService:  mockSocketService,
			UserService:    mockUserService,
		})

		form := url.Values{}
		form.Add("text", *mockMessage.Text)

		request, err := http.NewRequest(http.MethodPost, "/api/messages/"+mockChannel.ID, strings.NewReader(form.Encode()))
		assert.NoError(t, err)
		request.Form = form

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockChannelService.AssertExpectations(t)
		mockMessageService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
		mockUserService.AssertExpectations(t)
	})
}

func TestHandler_CreateMessage_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	router := getAuthenticatedTestRouter(mockUser.ID)

	mockMessageService := new(mocks.MessageService)

	NewHandler(&Config{
		R:              router,
		MessageService: mockMessageService,
		MaxBodyBytes:   4 * 1024 * 1024,
	})

	testCases := []struct {
		name string
		body url.Values
	}{
		{
			name: "No file nor text",
			body: map[string][]string{},
		},
		{
			name: "Text empty and no file",
			body: map[string][]string{
				"text": {""},
			},
		},
		{
			name: "Text too long",
			body: map[string][]string{
				"text": {fixture.RandStringRunes(2001)},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			form := tc.body
			request, _ := http.NewRequest(http.MethodPost, "/api/messages/"+fixture.RandID(), strings.NewReader(form.Encode()))
			request.Form = form

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockMessageService.AssertNotCalled(t, "CreateMessage")
		})
	}
}

func TestHandler_UpdateMessage(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully updated", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage(authUser.ID, "")

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)
		mockMessageService.On("UpdateMessage", mockMessage).Return(nil)

		response := model.MessageResponse{
			Id:         mockMessage.ID,
			Text:       mockMessage.Text,
			CreatedAt:  mockMessage.CreatedAt,
			UpdatedAt:  mockMessage.UpdatedAt,
			Attachment: mockMessage.Attachment,
			User: model.MemberResponse{
				Id: authUser.ID,
			},
		}

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitEditMessage", mockMessage.ChannelId, &response).Return()

		reqBody, err := json.Marshal(gin.H{
			"text": *mockMessage.Text,
		})
		assert.NoError(t, err)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPut, "/api/messages/"+mockMessage.ID, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Message not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("message", id)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", id).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		reqBody, err := json.Marshal(gin.H{
			"text": fixture.RandStringRunes(12),
		})
		assert.NoError(t, err)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPut, "/api/messages/"+id, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertCalled(t, "Get", id)
		mockMessageService.AssertNotCalled(t, "UpdateMessage")
		mockSocketService.AssertNotCalled(t, "EmitEditMessage")
	})

	t.Run("Not the message author", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage("", "")

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)

		mockSocketService := new(mocks.SocketService)

		reqBody, err := json.Marshal(gin.H{
			"text": fixture.RandStringRunes(12),
		})
		assert.NoError(t, err)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPut, "/api/messages/"+mockMessage.ID, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.EditMessageError)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertCalled(t, "Get", mockMessage.ID)
		mockMessageService.AssertNotCalled(t, "UpdateMessage")
		mockSocketService.AssertNotCalled(t, "EmitEditMessage")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage("", "")

		mockMessageService := new(mocks.MessageService)
		mockSocketService := new(mocks.SocketService)

		reqBody, err := json.Marshal(gin.H{
			"text": fixture.RandStringRunes(12),
		})
		assert.NoError(t, err)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPut, "/api/messages/"+mockMessage.ID, bytes.NewBuffer(reqBody))
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

		mockMessageService.AssertNotCalled(t, "Get")
		mockMessageService.AssertNotCalled(t, "UpdateMessage")
		mockSocketService.AssertNotCalled(t, "EmitEditMessage")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage(authUser.ID, "")
		mockError := apperrors.NewInternal()

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)
		mockMessageService.On("UpdateMessage", mockMessage).Return(mockError)

		mockSocketService := new(mocks.SocketService)

		reqBody, err := json.Marshal(gin.H{
			"text": *mockMessage.Text,
		})
		assert.NoError(t, err)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodPut, "/api/messages/"+mockMessage.ID, bytes.NewBuffer(reqBody))
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})
}

func TestHandler_UpdateMessage_BadRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)

	mockUser := fixture.GetMockUser()
	router := getAuthenticatedTestRouter(mockUser.ID)

	mockMessageService := new(mocks.MessageService)

	NewHandler(&Config{
		R:              router,
		MessageService: mockMessageService,
		MaxBodyBytes:   4 * 1024 * 1024,
	})

	testCases := []struct {
		name string
		body url.Values
	}{
		{
			name: "Text is required",
			body: map[string][]string{},
		},
		{
			name: "Text empty",
			body: map[string][]string{
				"text": {""},
			},
		},
		{
			name: "Text too long",
			body: map[string][]string{
				"text": {fixture.RandStringRunes(2001)},
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {

			rr := httptest.NewRecorder()

			form := tc.body
			request, _ := http.NewRequest(http.MethodPut, "/api/messages/"+fixture.RandID(), strings.NewReader(form.Encode()))
			request.Form = form

			router.ServeHTTP(rr, request)

			assert.Equal(t, http.StatusBadRequest, rr.Code)
			mockMessageService.AssertNotCalled(t, "UpdateMessage")
		})
	}
}

func TestHandler_Delete_Message(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	t.Run("Successfully deleted", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockMessage := fixture.GetMockMessage(authUser.ID, mockChannel.ID)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)
		mockMessageService.On("DeleteMessage", mockMessage).Return(nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitDeleteMessage", mockChannel.ID, mockMessage.ID)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Message not found", func(t *testing.T) {
		id := fixture.RandID()
		mockError := apperrors.NewNotFound("message", id)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", id).Return(nil, mockError)

		mockChannelService := new(mocks.ChannelService)
		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+id, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertCalled(t, "Get", id)
		mockMessageService.AssertNotCalled(t, "DeleteMessage")
		mockChannelService.AssertNotCalled(t, "Get")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockSocketService.AssertNotCalled(t, "EmitDeleteMessage")
	})

	t.Run("Channel not found", func(t *testing.T) {
		id := fixture.RandID()
		mockMessage := fixture.GetMockMessage("", id)
		mockError := apperrors.NewNotFound("message", mockMessage.ID)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", id).Return(nil, mockError)

		mockGuildService := new(mocks.GuildService)
		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertCalled(t, "Get", mockMessage.ID)
		mockChannelService.AssertCalled(t, "Get", id)
		mockMessageService.AssertNotCalled(t, "DeleteMessage")
		mockGuildService.AssertNotCalled(t, "GetGuild")
		mockSocketService.AssertNotCalled(t, "EmitDeleteMessage")
	})

	t.Run("Delete in guild - guild owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild(authUser.ID)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockMessage := fixture.GetMockMessage("", mockChannel.ID)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)
		mockMessageService.On("DeleteMessage", mockMessage).Return(nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitDeleteMessage", mockChannel.ID, mockMessage.ID)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Delete in guild - not the guild owner", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockMessage := fixture.GetMockMessage("", mockChannel.ID)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.DeleteMessageError)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertCalled(t, "Get", mockMessage.ID)
		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockGuildService.AssertCalled(t, "GetGuild", mockGuild.ID)
		mockMessageService.AssertNotCalled(t, "DeleteMessage")
		mockSocketService.AssertNotCalled(t, "EmitDeleteMessage")
	})

	t.Run("Delete in DM - Author", func(t *testing.T) {
		mockChannel := fixture.GetMockDMChannel()
		mockMessage := fixture.GetMockMessage(authUser.ID, mockChannel.ID)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)
		mockMessageService.On("DeleteMessage", mockMessage).Return(nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitDeleteMessage", mockChannel.ID, mockMessage.ID)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(true)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockSocketService.AssertExpectations(t)
	})

	t.Run("Delete in DM - Not the author", func(t *testing.T) {
		mockChannel := fixture.GetMockDMChannel()
		mockMessage := fixture.GetMockMessage("", mockChannel.ID)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)
		mockMessageService.On("DeleteMessage", mockMessage).Return(nil)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		mockError := apperrors.NewAuthorization(apperrors.DeleteDMMessageError)
		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertCalled(t, "Get", mockMessage.ID)
		mockChannelService.AssertCalled(t, "Get", mockChannel.ID)
		mockMessageService.AssertNotCalled(t, "DeleteMessage")
		mockSocketService.AssertNotCalled(t, "EmitDeleteMessage")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockMessageService := new(mocks.MessageService)
		mockChannelService := new(mocks.ChannelService)
		mockSocketService := new(mocks.SocketService)

		router := getTestRouter()

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+fixture.RandID(), nil)
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

		mockMessageService.AssertNotCalled(t, "Get")
		mockChannelService.AssertNotCalled(t, "Get")
		mockMessageService.AssertNotCalled(t, "DeleteMessage")
		mockSocketService.AssertNotCalled(t, "EmitDeleteMessage")
	})

	t.Run("Server Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)
		mockMessage := fixture.GetMockMessage(authUser.ID, mockChannel.ID)

		mockMessageService := new(mocks.MessageService)
		mockMessageService.On("Get", mockMessage.ID).Return(mockMessage, nil)
		mockError := apperrors.NewInternal()
		mockMessageService.On("DeleteMessage", mockMessage).Return(mockError)

		mockChannelService := new(mocks.ChannelService)
		mockChannelService.On("Get", mockChannel.ID).Return(mockChannel, nil)

		mockGuildService := new(mocks.GuildService)
		mockGuildService.On("GetGuild", mockGuild.ID).Return(mockGuild, nil)

		mockSocketService := new(mocks.SocketService)

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:              router,
			MessageService: mockMessageService,
			GuildService:   mockGuildService,
			ChannelService: mockChannelService,
			SocketService:  mockSocketService,
		})

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		// use bytes.NewBuffer to create a reader
		request, err := http.NewRequest(http.MethodDelete, "/api/messages/"+mockMessage.ID, nil)
		assert.NoError(t, err)

		request.Header.Set("Content-Type", "application/json")

		respBody, _ := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockMessageService.AssertExpectations(t)
		mockChannelService.AssertExpectations(t)
		mockGuildService.AssertExpectations(t)
		mockSocketService.AssertNotCalled(t, "EmitDeleteMessage")
	})
}

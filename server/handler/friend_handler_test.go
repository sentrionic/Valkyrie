package handler

import (
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

func TestHandler_GetUserFriends(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	friends := make([]model.Friend, 0)

	for i := 0; i < 5; i++ {
		mockFriend := fixture.GetMockUser()
		friends = append(friends, model.Friend{
			Id:       mockFriend.ID,
			Username: mockFriend.Username,
			Image:    mockFriend.Image,
			IsOnline: false,
		})
	}

	t.Run("Successful Fetch", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetFriends", authUser.ID).Return(&friends, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account/me/friends", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(friends)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockFriendService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetFriends", authUser.ID).Return(&friends, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account/me/friends", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		mockFriendService.AssertNotCalled(t, "GetFriends", authUser.ID, "")
	})

	t.Run("Error", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetFriends", authUser.ID).Return(nil, fmt.Errorf("some error down call chain"))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account/me/friends", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewNotFound("user", authUser.ID)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockFriendService.AssertExpectations(t)
	})
}

func TestHandler_GetUserRequests(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	authUser := fixture.GetMockUser()

	requests := make([]model.FriendRequest, 0)

	for i := 0; i < 5; i++ {
		mockRequest := fixture.GetMockUser()
		requests = append(requests, model.FriendRequest{
			Id:       mockRequest.ID,
			Username: mockRequest.Username,
			Image:    mockRequest.Image,
			Type:     0,
		})
	}

	t.Run("Successful Fetch", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetRequests", authUser.ID).Return(&requests, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account/me/pending", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(requests)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockFriendService.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetRequests", authUser.ID).Return(&requests, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account/me/pending", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)

		mockFriendService.AssertNotCalled(t, "GetRequests", authUser.ID, "")
	})

	t.Run("Error", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetRequests", authUser.ID).Return(nil, fmt.Errorf("some error down call chain"))

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(authUser.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		request, err := http.NewRequest(http.MethodGet, "/api/account/me/pending", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewNotFound("user", authUser.ID)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())
		mockFriendService.AssertExpectations(t)
	})
}

func TestHandler_AcceptFriendRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	current := fixture.GetMockUser()

	t.Run("Successfully accepted request", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		mockUser.Requests = append(mockUser.Requests, *current)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockFriendService.On("SaveRequests", current).
			Run(func(args mock.Arguments) {
				current.Friends = append(current.Friends, *mockUser)
			}).
			Return(nil)

		mockFriendService.On("SaveRequests", mockUser).
			Run(func(args mock.Arguments) {
				mockUser.Friends = append(mockUser.Friends, *current)
			}).
			Return(nil)

		mockFriendService.On("DeleteRequest", current.ID, mockUser.ID).Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitAddFriend", current, mockUser).Return()

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/accept", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
	})

	t.Run("Member does not contain a request", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/accept", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "SaveRequests")
		mockSocketService.AssertNotCalled(t, "EmitAddFriendRequest")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockUserService := new(mocks.UserService)
		mockUserService.On("GetMemberById", id).Return(nil, nil)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/accept", id)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertNotCalled(t, "GetMemberById", id)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)

		mockError := apperrors.NewNotFound("user", mockUser.ID)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/accept", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "SaveRequests")
		mockSocketService.AssertNotCalled(t, "EmitAddFriendRequest")
	})

	t.Run("MemberId and UserId are the same", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/accept", current.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.AcceptYourselfError)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertNotCalled(t, "GetMemberById")
		mockFriendService.AssertNotCalled(t, "SaveRequests")
		mockFriendService.AssertNotCalled(t, "DeleteRequest")
		mockSocketService.AssertNotCalled(t, "EmitAddFriend")
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		mockUser.Requests = append(mockUser.Requests, *current)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockError := apperrors.NewBadRequest(apperrors.UnableAcceptError)
		mockFriendService.On("SaveRequests", mockUser).
			Return(mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/accept", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockSocketService.AssertNotCalled(t, "EmitAddFriend")
	})
}

func TestHandler_SendFriendRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	current := fixture.GetMockUser()

	t.Run("Successfully send request", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockFriendService.On("SaveRequests", current).
			Run(func(args mock.Arguments) {
				current.Requests = append(current.Requests, *mockUser)
			}).
			Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitAddFriendRequest", mockUser.ID, &model.FriendRequest{
			Id:       current.ID,
			Username: current.Username,
			Image:    current.Image,
			Type:     1,
		}).Return()

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
	})

	t.Run("Member already is a friend", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		current.Friends = append(current.Friends, *mockUser)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "SaveRequests")
		mockSocketService.AssertNotCalled(t, "EmitAddFriendRequest")
	})

	t.Run("Member already contains request", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		current.Requests = append(current.Requests, *mockUser)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "SaveRequests")
		mockSocketService.AssertNotCalled(t, "EmitAddFriendRequest")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockUserService := new(mocks.UserService)
		mockUserService.On("GetMemberById", id).Return(nil, nil)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", id)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertNotCalled(t, "GetMemberById", id)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)

		mockError := apperrors.NewNotFound("user", mockUser.ID)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "SaveRequests")
		mockSocketService.AssertNotCalled(t, "EmitAddFriendRequest")
	})

	t.Run("MemberId and UserId are the same", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", current.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.AddYourselfError)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertNotCalled(t, "GetMemberById")
		mockFriendService.AssertNotCalled(t, "SaveRequests")
		mockSocketService.AssertNotCalled(t, "EmitAddFriendRequest")
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockError := apperrors.NewBadRequest(apperrors.UnableAddError)
		mockFriendService.On("SaveRequests", current).
			Return(mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockSocketService.AssertNotCalled(t, "EmitAddFriendRequest")
	})
}

func TestHandler_RemoveFriend(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	current := fixture.GetMockUser()

	t.Run("Successfully removed friend", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		current.Friends = append(current.Friends, *mockUser)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockFriendService.On("RemoveFriend", mockUser.ID, current.ID).
			Run(func(args mock.Arguments) {
				current.Friends = make([]model.User, 0)
			}).
			Return(nil)

		mockSocketService := new(mocks.SocketService)
		mockSocketService.On("EmitRemoveFriend", current.ID, mockUser.ID).Return()

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
	})

	t.Run("Member was not a friend", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "RemoveFriend")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFriend")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockUserService := new(mocks.UserService)
		mockUserService.On("GetMemberById", id).Return(nil, nil)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", id)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertNotCalled(t, "GetMemberById", id)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)

		mockError := apperrors.NewNotFound("user", mockUser.ID)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(nil, mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "RemoveFriend")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFriend")
	})

	t.Run("MemberId and UserId are the same", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)
		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", current.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.RemoveYourselfError)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertNotCalled(t, "GetMemberById")
		mockFriendService.AssertNotCalled(t, "RemoveFriend")
		mockSocketService.AssertNotCalled(t, "EmitRemoveFriend")
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		current.Friends = append(current.Friends, *mockUser)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockError := apperrors.NewBadRequest(apperrors.UnableRemoveError)
		mockFriendService.On("RemoveFriend", mockUser.ID, current.ID).
			Return(mockError)

		mockSocketService := new(mocks.SocketService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
			SocketService: mockSocketService,
		})

		url := fmt.Sprintf("/api/account/%s/friend", mockUser.ID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockSocketService.AssertNotCalled(t, "EmitRemoveFriend")
	})
}

func TestHandler_CancelFriendRequest(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	current := fixture.GetMockUser()

	t.Run("Successfully canceled request", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		current.Requests = append(current.Requests, *mockUser)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockFriendService.On("DeleteRequest", mockUser.ID, current.ID).
			Run(func(args mock.Arguments) {
				current.Requests = make([]model.User, 0)
			}).
			Return(nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/cancel", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
	})

	t.Run("Current does not contain a request", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/cancel", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(true)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "DeleteRequest")
	})

	t.Run("Unauthorized", func(t *testing.T) {
		id := fixture.RandID()

		mockUserService := new(mocks.UserService)
		mockUserService.On("GetMemberById", id).Return(nil, nil)

		rr := httptest.NewRecorder()

		router := getTestRouter()

		NewHandler(&Config{
			R:           router,
			UserService: mockUserService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/cancel", id)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		assert.Equal(t, http.StatusUnauthorized, rr.Code)
		mockUserService.AssertNotCalled(t, "GetMemberById", id)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUser := fixture.GetMockUser()

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)

		mockError := apperrors.NewNotFound("user", mockUser.ID)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(nil, mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/cancel", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
		mockFriendService.AssertNotCalled(t, "DeleteRequest")
	})

	t.Run("MemberId and UserId are the same", func(t *testing.T) {
		mockFriendService := new(mocks.FriendService)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/cancel", current.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		mockError := apperrors.NewBadRequest(apperrors.CancelYourselfError)
		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertNotCalled(t, "GetMemberById")
		mockFriendService.AssertNotCalled(t, "DeleteRequest")
	})

	t.Run("Error", func(t *testing.T) {
		mockUser := fixture.GetMockUser()
		current.Requests = append(current.Requests, *mockUser)

		mockFriendService := new(mocks.FriendService)
		mockFriendService.On("GetMemberById", current.ID).Return(current, nil)
		mockFriendService.On("GetMemberById", mockUser.ID).Return(mockUser, nil)

		mockError := apperrors.NewBadRequest(apperrors.UnableRemoveError)
		mockFriendService.On("DeleteRequest", mockUser.ID, current.ID).
			Return(mockError)

		// a response recorder for getting written http response
		rr := httptest.NewRecorder()

		router := getAuthenticatedTestRouter(current.ID)

		NewHandler(&Config{
			R:             router,
			FriendService: mockFriendService,
		})

		url := fmt.Sprintf("/api/account/%s/friend/cancel", mockUser.ID)
		request, err := http.NewRequest(http.MethodPost, url, nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)

		respBody, err := json.Marshal(gin.H{
			"error": mockError,
		})
		assert.NoError(t, err)

		assert.Equal(t, mockError.Status(), rr.Code)
		assert.Equal(t, respBody, rr.Body.Bytes())

		mockFriendService.AssertExpectations(t)
	})
}

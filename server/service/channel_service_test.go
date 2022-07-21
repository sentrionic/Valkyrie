package service

import (
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestChannelService_CreateChannel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid := GenerateId()
		mockChannel := fixture.GetMockChannel("")

		params := &model.Channel{
			GuildID:  mockChannel.GuildID,
			Name:     mockChannel.Name,
			IsPublic: true,
		}

		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		mockChannelRepository.
			On("Create", params).
			Run(func(args mock.Arguments) {
				mockChannel.ID = uid
			}).Return(mockChannel, nil)

		channel, err := cs.CreateChannel(params)

		assert.NoError(t, err)

		assert.Equal(t, uid, mockChannel.ID)
		assert.Equal(t, channel, mockChannel)

		mockChannelRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel("")

		params := &model.Channel{
			GuildID:  mockChannel.GuildID,
			Name:     mockChannel.Name,
			IsPublic: true,
		}

		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		mockErr := apperrors.NewInternal()

		mockChannelRepository.
			On("Create", params).
			Return(nil, mockErr)

		channel, err := cs.CreateChannel(params)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())
		assert.Nil(t, channel)

		mockChannelRepository.AssertExpectations(t)
	})
}

func TestChannelService_AddDMChannelMembers(t *testing.T) {
	userId := fixture.RandID()
	ids := []string{userId, fixture.RandID()}
	channelId := fixture.RandID()

	t.Run("Success", func(t *testing.T) {
		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		mockChannelRepository.
			On("AddDMChannelMembers", mock.AnythingOfType("[]model.DMMember")).
			// Has to be added so the different IDs get accepted
			Return(nil)

		err := cs.AddDMChannelMembers(ids, channelId, userId)
		assert.NoError(t, err)

		mockChannelRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		mockError := apperrors.NewInternal()
		mockChannelRepository.
			On("AddDMChannelMembers", mock.AnythingOfType("[]model.DMMember")).
			// Has to be added so the different IDs get accepted
			Return(mockError)

		err := cs.AddDMChannelMembers(ids, channelId, userId)
		assert.Error(t, err)

		mockChannelRepository.AssertExpectations(t)
	})
}

func TestChannelService_IsChannelMember(t *testing.T) {
	mockUser := fixture.GetMockUser()

	t.Run("User is member of the DM", func(t *testing.T) {
		mockChannel := fixture.GetMockDMChannel()

		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		mockChannelRepository.On("FindDMByUserAndChannelId", mockChannel.ID, mockUser.ID).Return(fixture.RandID(), nil)

		err := cs.IsChannelMember(mockChannel, mockUser.ID)
		assert.NoError(t, err)
	})

	t.Run("User is not member of the DM", func(t *testing.T) {
		mockChannel := fixture.GetMockDMChannel()

		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		mockError := apperrors.NewAuthorization(apperrors.Unauthorized)
		mockChannelRepository.On("FindDMByUserAndChannelId", mockChannel.ID, mockUser.ID).Return("", mockError)

		err := cs.IsChannelMember(mockChannel, mockUser.ID)
		assert.Error(t, err)
		assert.Equal(t, err, mockError)
	})

	t.Run("User is member of the private channel", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel("")
		mockChannel.IsPublic = false
		mockChannel.PCMembers = append(mockChannel.PCMembers, *mockUser)

		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		err := cs.IsChannelMember(mockChannel, mockUser.ID)
		assert.NoError(t, err)
	})

	t.Run("User is not member of the private channel", func(t *testing.T) {
		mockChannel := fixture.GetMockChannel("")
		mockChannel.IsPublic = false

		mockChannelRepository := new(mocks.ChannelRepository)
		cs := NewChannelService(&CSConfig{
			ChannelRepository: mockChannelRepository,
		})

		err := cs.IsChannelMember(mockChannel, mockUser.ID)
		assert.Error(t, err)
		assert.Equal(t, err, apperrors.NewAuthorization(apperrors.Unauthorized))
	})

	t.Run("User is a guild member", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockGuild.Members = append(mockGuild.Members, *mockUser)
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildRepository := new(mocks.GuildRepository)
		cs := NewChannelService(&CSConfig{
			GuildRepository: mockGuildRepository,
		})

		mockGuildRepository.On("GetMember", mockUser.ID, *mockChannel.GuildID).Return(mockUser, nil)

		err := cs.IsChannelMember(mockChannel, mockUser.ID)
		assert.NoError(t, err)
	})

	t.Run("User is not a member of the guild", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")
		mockChannel := fixture.GetMockChannel(mockGuild.ID)

		mockGuildRepository := new(mocks.GuildRepository)
		cs := NewChannelService(&CSConfig{
			GuildRepository: mockGuildRepository,
		})

		mockError := apperrors.NewAuthorization(apperrors.Unauthorized)
		mockGuildRepository.On("GetMember", mockUser.ID, *mockChannel.GuildID).Return(nil, mockError)

		err := cs.IsChannelMember(mockChannel, mockUser.ID)
		assert.Error(t, err)
		assert.Equal(t, err, mockError)
	})
}

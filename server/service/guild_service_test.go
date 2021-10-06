package service

import (
	"context"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGuildService_CreateGuild(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()
		mockGuild := fixture.GetMockGuild("")

		params := &model.Guild{
			Name: mockGuild.Name,
		}

		mockGuildRepository := new(mocks.GuildRepository)
		gs := NewGuildService(&GSConfig{
			GuildRepository: mockGuildRepository,
		})

		mockGuildRepository.
			On("Create", params).
			Run(func(args mock.Arguments) {
				mockGuild.ID = uid
			}).Return(mockGuild, nil)

		guild, err := gs.CreateGuild(params)

		assert.NoError(t, err)

		assert.Equal(t, uid, mockGuild.ID)
		assert.Equal(t, guild, mockGuild)

		mockGuildRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockGuild := fixture.GetMockGuild("")

		params := &model.Guild{
			Name: mockGuild.Name,
		}

		mockGuildRepository := new(mocks.GuildRepository)
		gs := NewGuildService(&GSConfig{
			GuildRepository: mockGuildRepository,
		})

		mockErr := apperrors.NewInternal()

		mockGuildRepository.
			On("Create", params).
			Return(nil, mockErr)

		guild, err := gs.CreateGuild(params)

		// assert error is error we response with in mock
		assert.EqualError(t, err, mockErr.Error())
		assert.Nil(t, guild)

		mockGuildRepository.AssertExpectations(t)
	})
}

func TestGuildService_GenerateInviteLink(t *testing.T) {
	guildId := fixture.RandID()
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRedisRepository := new(mocks.RedisRepository)
		gs := NewGuildService(&GSConfig{
			RedisRepository: mockRedisRepository,
		})

		args := mock.Arguments{
			ctx,
			guildId,
			mock.AnythingOfType("string"),
			false,
		}

		mockRedisRepository.
			On("SaveInvite", args...).
			Run(func(args mock.Arguments) {}).
			Return(nil)

		link, err := gs.GenerateInviteLink(ctx, guildId, false)

		assert.NoError(t, err)

		assert.NotEqual(t, link, "")

		mockRedisRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRedisRepository := new(mocks.RedisRepository)
		gs := NewGuildService(&GSConfig{
			RedisRepository: mockRedisRepository,
		})

		args := mock.Arguments{
			ctx,
			guildId,
			mock.AnythingOfType("string"),
			false,
		}

		mockError := apperrors.NewInternal()
		mockRedisRepository.
			On("SaveInvite", args...).
			Run(func(args mock.Arguments) {}).
			Return(mockError)

		link, err := gs.GenerateInviteLink(ctx, guildId, false)

		assert.Error(t, err)

		assert.Equal(t, link, "")
		assert.Equal(t, err, mockError)

		mockRedisRepository.AssertExpectations(t)
	})
}

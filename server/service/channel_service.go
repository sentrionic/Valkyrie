package service

import (
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
)

// channelService acts as a struct for injecting an implementation of ChannelRepository
// for use in service methods
type channelService struct {
	ChannelRepository model.ChannelRepository
	GuildRepository   model.GuildRepository
}

// CSConfig will hold repositories that will eventually be injected into
// this service layer
type CSConfig struct {
	ChannelRepository model.ChannelRepository
	GuildRepository   model.GuildRepository
}

// NewChannelService is a factory function for
// initializing a ChannelService with its repository layer dependencies
func NewChannelService(c *CSConfig) model.ChannelService {
	return &channelService{
		ChannelRepository: c.ChannelRepository,
		GuildRepository:   c.GuildRepository,
	}
}

func (c *channelService) CreateChannel(channel *model.Channel) (*model.Channel, error) {
	id, err := GenerateId()

	if err != nil {
		return nil, err
	}

	channel.ID = id

	return c.ChannelRepository.Create(channel)
}

func (c *channelService) GetChannels(userId string, guildId string) (*[]model.ChannelResponse, error) {
	return c.ChannelRepository.Get(userId, guildId)
}

func (c *channelService) Get(channelId string) (*model.Channel, error) {
	return c.ChannelRepository.GetById(channelId)
}

func (c *channelService) GetPrivateChannelMembers(channelId string) (*[]string, error) {
	return c.ChannelRepository.GetPrivateChannelMembers(channelId)
}

func (c *channelService) GetDirectMessages(userId string) (*[]model.DirectMessage, error) {
	return c.ChannelRepository.GetDirectMessages(userId)
}

func (c *channelService) GetDirectMessageChannel(userId string, memberId string) (*string, error) {
	return c.ChannelRepository.GetDirectMessageChannel(userId, memberId)
}

func (c *channelService) AddDMChannelMembers(memberIds []string, channelId string, userId string) error {
	var members []model.DMMember
	for _, mId := range memberIds {
		id, err := GenerateId()

		if err != nil {
			return err
		}

		member := model.DMMember{
			ID:        id,
			UserID:    mId,
			ChannelId: channelId,
			IsOpen:    userId == mId,
		}
		members = append(members, member)
	}

	return c.ChannelRepository.AddDMChannelMembers(members)
}

func (c *channelService) SetDirectMessageStatus(dmId string, userId string, isOpen bool) error {
	return c.ChannelRepository.SetDirectMessageStatus(dmId, userId, isOpen)
}

func (c *channelService) DeleteChannel(channel *model.Channel) error {
	return c.ChannelRepository.DeleteChannel(channel)
}

func (c *channelService) UpdateChannel(channel *model.Channel) error {
	return c.ChannelRepository.UpdateChannel(channel)
}

func (c *channelService) CleanPCMembers(channelId string) error {
	return c.ChannelRepository.CleanPCMembers(channelId)
}

func (c *channelService) AddPrivateChannelMembers(memberIds []string, channelId string) error {
	return c.ChannelRepository.AddPrivateChannelMembers(memberIds, channelId)
}

func (c *channelService) RemovePrivateChannelMembers(memberIds []string, channelId string) error {
	return c.ChannelRepository.RemovePrivateChannelMembers(memberIds, channelId)
}

func (c *channelService) OpenDMForAll(dmId string) error {
	return c.ChannelRepository.OpenDMForAll(dmId)
}

func (c *channelService) GetDMByUserAndChannel(userId string, channelId string) (string, error) {
	return c.ChannelRepository.FindDMByUserAndChannelId(channelId, userId)
}

// IsChannelMember checks if the user has access to the given channel.
// Returns an error if they do not, otherwise nil
func (c *channelService) IsChannelMember(channel *model.Channel, userId string) error {
	// Check if user has access to the channel if it's private
	if !channel.IsPublic {
		// Channel is DM -> Check if one of the members
		if channel.IsDM {
			id, err := c.ChannelRepository.FindDMByUserAndChannelId(channel.ID, userId)

			if err != nil || id == "" {
				return apperrors.NewAuthorization(apperrors.Unauthorized)
			}
			return nil
		}
		// Channel is private
		for _, member := range channel.PCMembers {
			if member.ID == userId {
				return nil
			}
		}
		return apperrors.NewAuthorization(apperrors.Unauthorized)
	}

	// Check if user has access to the channel
	member, err := c.GuildRepository.GetMember(userId, *channel.GuildID)
	if err != nil || member.ID == "" {
		return apperrors.NewAuthorization(apperrors.Unauthorized)
	}
	return nil
}

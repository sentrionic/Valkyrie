package service

import (
	"context"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/sentrionic/valkyrie/model"
)

// GuildService acts as a struct for injecting an implementation of GuildRepository
// for use in service methods
type guildService struct {
	UserRepository    model.UserRepository
	FileRepository    model.FileRepository
	RedisRepository   model.RedisRepository
	GuildRepository   model.GuildRepository
	ChannelRepository model.ChannelRepository
}

// GSConfig will hold repositories that will eventually be injected into
// this service layer
type GSConfig struct {
	UserRepository    model.UserRepository
	FileRepository    model.FileRepository
	RedisRepository   model.RedisRepository
	GuildRepository   model.GuildRepository
	ChannelRepository model.ChannelRepository
}

// NewGuildService is a factory function for
// initializing a GuildService with its repository layer dependencies
func NewGuildService(c *GSConfig) model.GuildService {
	return &guildService{
		UserRepository:    c.UserRepository,
		FileRepository:    c.FileRepository,
		RedisRepository:   c.RedisRepository,
		GuildRepository:   c.GuildRepository,
		ChannelRepository: c.ChannelRepository,
	}
}

func (g *guildService) GetUserGuilds(uid string) (*[]model.GuildResponse, error) {
	return g.GuildRepository.List(uid)
}

func (g *guildService) GetGuildMembers(userId string, guildId string) (*[]model.MemberResponse, error) {
	return g.GuildRepository.GuildMembers(userId, guildId)
}

func (g *guildService) CreateGuild(guild *model.Guild) (*model.Guild, error) {
	id, err := GenerateId()

	if err != nil {
		return nil, err
	}

	guild.ID = id

	return g.GuildRepository.Create(guild)
}

func (g *guildService) GetUser(uid string) (*model.User, error) {
	return g.GuildRepository.FindUserByID(uid)
}

func (g *guildService) GetGuild(id string) (*model.Guild, error) {
	return g.GuildRepository.FindByID(id)
}

func (g *guildService) GenerateInviteLink(ctx context.Context, guildId string, isPermanent bool) (string, error) {
	id, err := gonanoid.Nanoid(8)

	if err != nil {
		return "", err
	}

	if err = g.RedisRepository.SaveInvite(ctx, guildId, id, isPermanent); err != nil {
		return "", err
	}

	return id, nil
}

func (g *guildService) UpdateGuild(guild *model.Guild) error {
	return g.GuildRepository.Save(guild)
}

func (g *guildService) GetGuildIdFromInvite(ctx context.Context, token string) (string, error) {
	return g.RedisRepository.GetInvite(ctx, token)
}

func (g *guildService) GetDefaultChannel(guildId string) (*model.Channel, error) {
	return g.ChannelRepository.GetGuildDefault(guildId)
}

func (g *guildService) InvalidateInvites(ctx context.Context, guild *model.Guild) {
	g.RedisRepository.InvalidateInvites(ctx, guild)
}

func (g *guildService) RemoveMember(userId string, guildId string) error {
	return g.GuildRepository.RemoveMember(userId, guildId)
}

func (g *guildService) DeleteGuild(guildId string) error {
	return g.GuildRepository.Delete(guildId)
}

func (g *guildService) UnbanMember(userId string, guildId string) error {
	return g.GuildRepository.UnbanMember(userId, guildId)
}

func (g *guildService) GetBanList(guildId string) (*[]model.BanResponse, error) {
	return g.GuildRepository.GetBanList(guildId)
}

func (g *guildService) GetMemberSettings(userId string, guildId string) (*model.MemberSettings, error) {
	return g.GuildRepository.GetMemberSettings(userId, guildId)
}

func (g *guildService) UpdateMemberSettings(settings *model.MemberSettings, userId string, guildId string) error {
	return g.GuildRepository.UpdateMemberSettings(settings, userId, guildId)
}

func (g *guildService) FindUsersByIds(ids []string, guildId string) (*[]model.User, error) {
	return g.GuildRepository.FindUsersByIds(ids, guildId)
}

func (g *guildService) UpdateMemberLastSeen(userId, guildId string) error {
	return g.GuildRepository.UpdateMemberLastSeen(userId, guildId)
}

func (g *guildService) RemoveVCMember(userId, guildId string) error {
	return g.GuildRepository.RemoveVCMember(userId, guildId)
}

func (g *guildService) GetVCMembers(guildId string) (*[]model.VCMemberResponse, error) {
	return g.GuildRepository.VCMembers(guildId)
}

func (g *guildService) UpdateVCMember(isMuted, isDeafened bool, userId, guildId string) error {
	return g.GuildRepository.UpdateVCMember(isMuted, isDeafened, userId, guildId)
}

func (g *guildService) GetVCMember(userId, guildId string) (*model.VCMember, error) {
	return g.GuildRepository.GetVCMember(userId, guildId)
}

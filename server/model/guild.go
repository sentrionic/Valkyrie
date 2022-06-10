package model

import (
	"context"
	"github.com/lib/pq"
	"time"
)

// Guild represents the server many users can chat in.
type Guild struct {
	BaseModel
	Name        string `gorm:"not null"`
	OwnerId     string `gorm:"not null"`
	Icon        *string
	InviteLinks pq.StringArray `gorm:"type:text[]"`
	Members     []User         `gorm:"many2many:members;constraint:OnDelete:CASCADE;"`
	Channels    []Channel      `gorm:"constraint:OnDelete:CASCADE;"`
	Bans        []User         `gorm:"many2many:bans;constraint:OnDelete:CASCADE;"`
	VCMembers   []User         `gorm:"many2many:vc_members;constraint:OnDelete:CASCADE;"`
}

// GuildResponse contains all info to display a guild.
// The DefaultChannelId is the channel the user first gets directed to
// and is the oldest channel of the guild.
type GuildResponse struct {
	Id               string    `json:"id"`
	Name             string    `json:"name"`
	OwnerId          string    `json:"ownerId"`
	Icon             *string   `json:"icon"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	HasNotification  bool      `json:"hasNotification"`
	DefaultChannelId string    `json:"default_channel_id"`
} //@name GuildResponse

// SerializeGuild returns the guild API response.
// The DefaultChannelId represents the default channel the user gets send to.
func (g Guild) SerializeGuild(channelId string) GuildResponse {
	return GuildResponse{
		Id:               g.ID,
		Name:             g.Name,
		OwnerId:          g.OwnerId,
		Icon:             g.Icon,
		CreatedAt:        g.CreatedAt,
		UpdatedAt:        g.UpdatedAt,
		HasNotification:  false,
		DefaultChannelId: channelId,
	}
}

// GuildService defines methods related to guild operations the handler layer expects
// any service it interacts with to implement
type GuildService interface {
	GetUser(uid string) (*User, error)
	GetGuild(id string) (*Guild, error)
	GetUserGuilds(uid string) (*[]GuildResponse, error)
	GetGuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	GetVCMembers(guildId string) (*[]VCMemberResponse, error)
	CreateGuild(guild *Guild) (*Guild, error)
	GenerateInviteLink(ctx context.Context, guildId string, isPermanent bool) (string, error)
	UpdateGuild(guild *Guild) error
	GetGuildIdFromInvite(ctx context.Context, token string) (string, error)
	GetDefaultChannel(guildId string) (*Channel, error)
	InvalidateInvites(ctx context.Context, guild *Guild)
	RemoveMember(userId string, guildId string) error
	UnbanMember(userId string, guildId string) error
	DeleteGuild(guildId string) error
	GetBanList(guildId string) (*[]BanResponse, error)
	GetMemberSettings(userId string, guildId string) (*MemberSettings, error)
	UpdateMemberSettings(settings *MemberSettings, userId string, guildId string) error
	FindUsersByIds(ids []string, guildId string) (*[]User, error)
	UpdateMemberLastSeen(userId, guildId string) error
	RemoveVCMember(userId, guildId string) error
	UpdateVCMember(isMuted, isDeafened bool, userId, guildId string) error
	GetVCMember(userId, guildId string) (*VCMember, error)
}

// GuildRepository defines methods related to guild db operations the service layer expects
// any repository it interacts with to implement
type GuildRepository interface {
	FindUserByID(uid string) (*User, error)
	FindByID(id string) (*Guild, error)
	List(uid string) (*[]GuildResponse, error)
	GuildMembers(userId string, guildId string) (*[]MemberResponse, error)
	VCMembers(guildId string) (*[]VCMemberResponse, error)
	Create(guild *Guild) (*Guild, error)
	Save(guild *Guild) error
	RemoveMember(userId string, guildId string) error
	Delete(guildId string) error
	UnbanMember(userId string, guildId string) error
	GetBanList(guildId string) (*[]BanResponse, error)
	GetMemberSettings(userId string, guildId string) (*MemberSettings, error)
	UpdateMemberSettings(settings *MemberSettings, userId string, guildId string) error
	FindUsersByIds(ids []string, guildId string) (*[]User, error)
	GetMember(userId, guildId string) (*User, error)
	UpdateMemberLastSeen(userId, guildId string) error
	RemoveVCMember(userId, guildId string) error
	GetMemberIds(guildId string) (*[]string, error)
	UpdateVCMember(isMuted, isDeafened bool, userId, guildId string) error
	GetVCMember(userId, guildId string) (*VCMember, error)
}

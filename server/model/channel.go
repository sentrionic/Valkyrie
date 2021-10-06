package model

import "time"

// Channel represents a text channel in a guild
// or a text channel for DMs between users.
// GuildID should only be nil if it is a DM channel
// PCMembers should only be used if the channel is private.
type Channel struct {
	BaseModel
	GuildID      *string   `gorm:"index"`
	Name         string    `gorm:"name"`
	IsPublic     bool      `gorm:"index"`
	IsDM         bool      `gorm:"is_dm"`
	LastActivity time.Time `gorm:"autoCreateTime"`
	PCMembers    []User    `gorm:"many2many:pcmembers;constraint:OnDelete:CASCADE;"`
	Messages     []Message `gorm:"constraint:OnDelete:CASCADE;"`
}

// ChannelResponse is the JSON response of the channel
type ChannelResponse struct {
	Id              string    `json:"id"`
	Name            string    `json:"name"`
	IsPublic        bool      `json:"isPublic"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	HasNotification bool      `json:"hasNotification"`
} //@name Channel

// SerializeChannel returns the channel API response.
func (c Channel) SerializeChannel() ChannelResponse {
	return ChannelResponse{
		Id:              c.ID,
		Name:            c.Name,
		IsPublic:        c.IsPublic,
		CreatedAt:       c.CreatedAt,
		UpdatedAt:       c.UpdatedAt,
		HasNotification: false,
	}
}

// ChannelService defines methods related to channel operations the handler layer expects
// any service it interacts with to implement
type ChannelService interface {
	CreateChannel(channel *Channel) (*Channel, error)
	GetChannels(userId string, guildId string) (*[]ChannelResponse, error)
	Get(channelId string) (*Channel, error)
	GetPrivateChannelMembers(channelId string) (*[]string, error)
	GetDirectMessages(userId string) (*[]DirectMessage, error)
	GetDirectMessageChannel(userId string, memberId string) (*string, error)
	GetDMByUserAndChannel(userId string, channelId string) (string, error)
	AddDMChannelMembers(memberIds []string, channelId string, userId string) error
	SetDirectMessageStatus(dmId string, userId string, isOpen bool) error
	DeleteChannel(channel *Channel) error
	UpdateChannel(channel *Channel) error
	CleanPCMembers(channelId string) error
	AddPrivateChannelMembers(memberIds []string, channelId string) error
	RemovePrivateChannelMembers(memberIds []string, channelId string) error
	IsChannelMember(channel *Channel, userId string) error
	OpenDMForAll(dmId string) error
}

// ChannelRepository defines methods related to channel db operations the service layer expects
// any repository it interacts with to implement
type ChannelRepository interface {
	Create(channel *Channel) (*Channel, error)
	GetGuildDefault(guildId string) (*Channel, error)
	Get(userId string, guildId string) (*[]ChannelResponse, error)
	GetDirectMessages(userId string) (*[]DirectMessage, error)
	GetDirectMessageChannel(userId string, memberId string) (*string, error)
	GetById(channelId string) (*Channel, error)
	GetPrivateChannelMembers(channelId string) (*[]string, error)
	AddDMChannelMembers(members []DMMember) error
	SetDirectMessageStatus(dmId string, userId string, isOpen bool) error
	DeleteChannel(channel *Channel) error
	UpdateChannel(channel *Channel) error
	CleanPCMembers(channelId string) error
	AddPrivateChannelMembers(memberIds []string, channelId string) error
	RemovePrivateChannelMembers(memberIds []string, channelId string) error
	FindDMByUserAndChannelId(channelId, userId string) (string, error)
	OpenDMForAll(dmId string) error
	GetDMMemberIds(channelId string) (*[]string, error)
}

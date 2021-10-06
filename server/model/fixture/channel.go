package fixture

import (
	"github.com/sentrionic/valkyrie/model"
	"time"
)

// GetMockChannel returns a mock channel. If guildId is not empty it will set GuildID to that id.
func GetMockChannel(guildId string) *model.Channel {

	var guild *string = nil
	if guildId != "" {
		guild = &guildId
	}

	return &model.Channel{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		GuildID:      guild,
		Name:         RandStr(8),
		IsPublic:     true,
		LastActivity: time.Now(),
	}
}

// GetMockDMChannel returns a mock channel that has IsDM set to true and does not belong to a guild.
func GetMockDMChannel() *model.Channel {
	return &model.Channel{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:         RandID(),
		IsDM:         true,
		LastActivity: time.Now(),
	}
}

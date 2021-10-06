package fixture

import (
	"github.com/sentrionic/valkyrie/model"
	"time"
)

// GetMockGuild returns a mock guild with the given uid as the owner.
func GetMockGuild(uid string) *model.Guild {
	ownerId := RandID()
	if uid != "" {
		ownerId = uid
	}

	return &model.Guild{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:    RandStr(8),
		OwnerId: ownerId,
	}
}

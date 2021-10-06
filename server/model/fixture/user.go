package fixture

import (
	"github.com/sentrionic/valkyrie/model"
	"time"
)

// GetMockUser returns a mock user
func GetMockUser() *model.User {
	email := Email()
	return &model.User{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Username: Username(),
		Email:    email,
		Password: RandStr(8),
		Image:    generateAvatar(email),
	}
}

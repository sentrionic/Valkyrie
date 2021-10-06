package fixture

import (
	"github.com/sentrionic/valkyrie/model"
	"time"
)

// GetMockMessage returns a mock message with the given uid as the owner and the cid as the ChannelId
func GetMockMessage(uid, cid string) *model.Message {
	text := RandStringRunes(100)

	ownerId := RandID()
	if uid != "" {
		ownerId = uid
	}

	return &model.Message{
		BaseModel: model.BaseModel{
			ID:        RandID(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Text:       &text,
		UserId:     ownerId,
		ChannelId:  cid,
		Attachment: nil,
	}
}

// GetMockMessageResponse returns a mock message response with the given uid as the owner and the cid as the ChannelId
func GetMockMessageResponse(uid, cid string) *model.MessageResponse {
	message := GetMockMessage(uid, cid)
	user := GetMockUser()

	return &model.MessageResponse{
		Id:         message.ID,
		Text:       message.Text,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
		Attachment: message.Attachment,
		User: model.MemberResponse{
			Id:        user.ID,
			Username:  user.Username,
			Image:     user.Image,
			IsOnline:  user.IsOnline,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Nickname:  nil,
			Color:     nil,
			IsFriend:  false,
		},
	}
}

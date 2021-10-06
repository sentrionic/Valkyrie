package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
	"log"
	"time"
)

// messageRepository is data/repository implementation
// of service layer MessageRepository
type messageRepository struct {
	DB *gorm.DB
}

// NewMessageRepository is a factory for initializing Message Repositories
func NewMessageRepository(db *gorm.DB) model.MessageRepository {
	return &messageRepository{
		DB: db,
	}
}

// messageQuery represents the fetched fields for GetMessages
type messageQuery struct {
	Id            string
	Text          *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FileType      *string
	Url           *string
	Filename      *string
	AttachmentId  *string
	UserId        string
	UserCreatedAt time.Time
	UserUpdatedAt time.Time
	Username      string
	Image         string
	IsOnline      bool
	Nickname      *string
	Color         *string
	IsFriend      bool
}

// GetMessages returns the 35 most recent messages for the given channel.
// If a cursor is specified it returns the 35 messages after the cursor.
func (r *messageRepository) GetMessages(userId string, channel *model.Channel, cursor string) (*[]model.MessageResponse, error) {
	var result []messageQuery

	memberSelect := ""
	memberJoin := ""
	memberWhere := ""

	// If the channel is not a DM channel, also fetch the message author's settings
	if !channel.IsDM {
		memberSelect = "member.nickname, member.color,"
		memberJoin = "LEFT JOIN members member on messages.user_id = member.user_id"
		memberWhere = fmt.Sprintf("AND member.guild_id = %s::text", *channel.GuildID)
	}

	crs := ""
	if cursor != "" {
		// Remove the timezone from the string since it's stored differently in the DB
		date := cursor[:len(cursor)-6]
		crs = fmt.Sprintf("AND messages.created_at < '%s'", date)
	}

	err := r.DB.
		Raw(fmt.Sprintf(`
		SELECT messages.id,
			messages.text,
			messages.created_at,
			messages.updated_at,
			a.file_type,
			a.url,
			a.filename,
			a.id                as "attachment_id",
			users.id         as "user_id",
			users.created_at as "user_created_at",
			users.updated_at as "user_updated_at",
			users.username,
			users.image,
			users.is_online,
			%s 
			EXISTS(
			  SELECT 1
			  FROM users
			   LEFT JOIN friends f ON users.id = f.user_id
			  WHERE f.friend_id = messages.user_id
				AND f.user_id = @userId) as is_friend
		FROM messages
		LEFT JOIN "users"
		ON users.id = messages.user_id
		LEFT JOIN attachments a
		ON a.message_id = messages.id
		%s
		WHERE messages.channel_id = @channelId
		%s 
		%s 
		ORDER BY messages.created_at DESC
		LIMIT 35
`, memberSelect, memberJoin, memberWhere, crs),
			sql.Named("userId", userId),
			sql.Named("channelId", channel.ID)).
		Scan(&result).Error

	var messages []model.MessageResponse

	// Turn messageQuery results into MessageResponse
	for _, m := range result {

		var attachment *model.Attachment = nil
		if m.AttachmentId != nil {
			attachment = &model.Attachment{
				Url:      *m.Url,
				FileType: *m.FileType,
				Filename: *m.Filename,
			}
		}

		message := model.MessageResponse{
			Id:         m.Id,
			Text:       m.Text,
			CreatedAt:  m.CreatedAt,
			UpdatedAt:  m.UpdatedAt,
			Attachment: attachment,
			User: model.MemberResponse{
				Id:        m.UserId,
				Username:  m.Username,
				Image:     m.Image,
				IsOnline:  m.IsOnline,
				CreatedAt: m.UserCreatedAt,
				UpdatedAt: m.UserUpdatedAt,
				Nickname:  m.Nickname,
				Color:     m.Color,
				IsFriend:  m.IsFriend,
			},
		}
		messages = append(messages, message)
	}

	return &messages, err
}

// CreateMessage inserts the message in the DB
func (r *messageRepository) CreateMessage(message *model.Message) (*model.Message, error) {
	if result := r.DB.Create(&message); result.Error != nil {
		log.Printf("Could not create a message for user: %v. Reason: %v\n", message.UserId, result.Error)
		return nil, apperrors.NewInternal()
	}

	return message, nil
}

// UpdateMessage updates the message in the DB
func (r *messageRepository) UpdateMessage(message *model.Message) error {
	if result := r.DB.Save(&message); result.Error != nil {
		log.Printf("Could not update message with id: %v. Reason: %v\n", message.ID, result.Error)
		return apperrors.NewInternal()
	}
	return nil
}

// DeleteMessage removes the message from the DB
func (r *messageRepository) DeleteMessage(message *model.Message) error {
	if result := r.DB.Delete(message); result.Error != nil {
		log.Printf("Could not delete message with id: %v. Reason: %v\n", message.ID, result.Error)
		return apperrors.NewInternal()
	}
	return nil
}

// GetById fetches the message for the given id
func (r *messageRepository) GetById(messageId string) (*model.Message, error) {
	message := &model.Message{}

	if result := r.DB.Where("id = ?", messageId).First(message); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return message, apperrors.NewNotFound("message", messageId)
		}
		return message, apperrors.NewInternal()
	}

	return message, nil
}

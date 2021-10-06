package model

import (
	"mime/multipart"
	"time"
)

// Message represents a text message in a channel.
// It may contain an Attachment that is displayed instead of text.
type Message struct {
	BaseModel
	Text       *string
	UserId     string      `gorm:"index;constraint:OnDelete:CASCADE;"`
	ChannelId  string      `gorm:"index;constraint:OnDelete:CASCADE;"`
	Attachment *Attachment `gorm:"constraint:OnDelete:CASCADE;"`
}

// MessageResponse is the API response of a Message
type MessageResponse struct {
	Id         string         `json:"id"`
	Text       *string        `json:"text"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Attachment *Attachment    `json:"attachment"`
	User       MemberResponse `json:"user"`
} //@name Message

// Attachment represents a message attachment that displays
// a file instead of text.
type Attachment struct {
	ID        string    `gorm:"primaryKey" json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Url       string    `json:"url"`
	FileType  string    `json:"filetype"`
	Filename  string    `json:"filename"`
	MessageId string    `gorm:"index;constraint:OnDelete:CASCADE;" json:"-"`
} //@name Attachment

// MessageService defines methods related to message operations the handler layer expects
// any service it interacts with to implement
type MessageService interface {
	GetMessages(userId string, channel *Channel, cursor string) (*[]MessageResponse, error)
	CreateMessage(params *Message) (*Message, error)
	UpdateMessage(message *Message) error
	DeleteMessage(message *Message) error
	UploadFile(header *multipart.FileHeader, channelId string) (*Attachment, error)
	Get(messageId string) (*Message, error)
}

// MessageRepository defines methods related message db operations the service layer expects
// any repository it interacts with to implement
type MessageRepository interface {
	GetMessages(userId string, channel *Channel, cursor string) (*[]MessageResponse, error)
	CreateMessage(params *Message) (*Message, error)
	UpdateMessage(message *Message) error
	DeleteMessage(message *Message) error
	GetById(messageId string) (*Message, error)
}

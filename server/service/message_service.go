package service

import (
	"fmt"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/sentrionic/valkyrie/model"
	"log"
	"mime/multipart"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// messageService acts as a struct for injecting an implementation of MessageRepository
// for use in service methods
type messageService struct {
	MessageRepository model.MessageRepository
	FileRepository    model.FileRepository
}

// MSConfig will hold repositories that will eventually be injected into
// this service layer
type MSConfig struct {
	MessageRepository model.MessageRepository
	FileRepository    model.FileRepository
}

// NewMessageService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewMessageService(c *MSConfig) model.MessageService {
	return &messageService{
		MessageRepository: c.MessageRepository,
		FileRepository:    c.FileRepository,
	}
}

func (m *messageService) GetMessages(userId string, channel *model.Channel, cursor string) (*[]model.MessageResponse, error) {
	return m.MessageRepository.GetMessages(userId, channel, cursor)
}

func (m *messageService) CreateMessage(params *model.Message) (*model.Message, error) {
	params.ID = GenerateId()

	return m.MessageRepository.CreateMessage(params)
}

func (m *messageService) UpdateMessage(message *model.Message) error {
	return m.MessageRepository.UpdateMessage(message)
}

func (m *messageService) DeleteMessage(message *model.Message) error {
	if message.Attachment != nil {
		if err := m.FileRepository.DeleteImage(message.Attachment.Filename); err != nil {
			log.Printf("Error deleting file from S3: %s", err)
		}
	}

	return m.MessageRepository.DeleteMessage(message)
}

func (m *messageService) UploadFile(header *multipart.FileHeader, channelId string) (*model.Attachment, error) {

	filename := formatName(header.Filename)
	mimetype := header.Header.Get("Content-Type")

	attachment := model.Attachment{
		FileType: mimetype,
		Filename: filename,
	}

	attachment.ID = GenerateId()

	directory := fmt.Sprintf("channels/%s", channelId)
	url, err := m.FileRepository.UploadFile(header, directory, filename, mimetype)

	if err != nil {
		return nil, err
	}

	attachment.Url = url

	return &attachment, nil
}

func (m *messageService) Get(messageId string) (*model.Message, error) {
	return m.MessageRepository.GetById(messageId)
}

var re = regexp.MustCompile(`/[^a-z0-9]/g`)

func formatName(filename string) string {
	ext := path.Ext(filename)
	id, _ := gonanoid.Nanoid(5)
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	filename = strings.ToLower(filename)
	filename = re.ReplaceAllString(filename, "-")
	return fmt.Sprintf("%s-%s%s", id, filename, ext)
}

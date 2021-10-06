package service

import (
	"fmt"
	"github.com/sentrionic/valkyrie/mocks"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"github.com/sentrionic/valkyrie/model/fixture"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestGuildService_CreateMessage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		uid, _ := GenerateId()
		mockMessage := fixture.GetMockMessage("", "")

		params := &model.Message{
			UserId:    mockMessage.UserId,
			ChannelId: mockMessage.ChannelId,
			Text:      mockMessage.Text,
		}

		mockMessageRepository := new(mocks.MessageRepository)
		ms := NewMessageService(&MSConfig{
			MessageRepository: mockMessageRepository,
		})

		mockMessageRepository.
			On("CreateMessage", params).
			Run(func(args mock.Arguments) {
				mockMessage.ID = uid
			}).Return(mockMessage, nil)

		message, err := ms.CreateMessage(params)

		assert.NoError(t, err)

		assert.Equal(t, uid, mockMessage.ID)
		assert.Equal(t, message, mockMessage)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage("", "")

		params := &model.Message{
			UserId:    mockMessage.UserId,
			ChannelId: mockMessage.ChannelId,
			Text:      mockMessage.Text,
		}

		mockMessageRepository := new(mocks.MessageRepository)
		ms := NewMessageService(&MSConfig{
			MessageRepository: mockMessageRepository,
		})

		mockErr := apperrors.NewInternal()
		mockMessageRepository.
			On("CreateMessage", params).
			Return(nil, mockErr)

		message, err := ms.CreateMessage(params)

		assert.EqualError(t, err, mockErr.Error())
		assert.Nil(t, message)
	})
}

func TestGuildService_DeleteMessage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage("", "")

		mockMessageRepository := new(mocks.MessageRepository)
		ms := NewMessageService(&MSConfig{
			MessageRepository: mockMessageRepository,
		})

		mockMessageRepository.
			On("DeleteMessage", mockMessage).
			Return(nil)

		err := ms.DeleteMessage(mockMessage)

		assert.NoError(t, err)

		mockMessageRepository.AssertExpectations(t)
	})

	t.Run("Success with attachment", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage("", "")
		mockMessage.Attachment = &model.Attachment{
			Filename: fixture.RandStr(12),
		}

		mockMessageRepository := new(mocks.MessageRepository)
		mockFileRepository := new(mocks.FileRepository)
		ms := NewMessageService(&MSConfig{
			MessageRepository: mockMessageRepository,
			FileRepository:    mockFileRepository,
		})

		mockFileRepository.On("DeleteImage", mockMessage.Attachment.Filename).Return(nil)

		mockMessageRepository.
			On("DeleteMessage", mockMessage).
			Return(nil)

		err := ms.DeleteMessage(mockMessage)

		assert.NoError(t, err)

		mockMessageRepository.AssertExpectations(t)
		mockFileRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockMessage := fixture.GetMockMessage("", "")

		mockMessageRepository := new(mocks.MessageRepository)
		ms := NewMessageService(&MSConfig{
			MessageRepository: mockMessageRepository,
		})

		mockError := apperrors.NewInternal()

		mockMessageRepository.
			On("DeleteMessage", mockMessage).
			Return(mockError)

		err := ms.DeleteMessage(mockMessage)

		assert.EqualError(t, err, mockError.Error())

		mockMessageRepository.AssertExpectations(t)
	})
}

func TestMessageService_UploadFile(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		imageURL := "https://imageurl.com/jdfkj34kljl"
		channelId := fixture.RandID()
		id, err := GenerateId()
		assert.NoError(t, err)

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := fmt.Sprintf("channels/%s", channelId)

		attachment := model.Attachment{
			FileType: imageFileHeader.Header.Get("Content-Type"),
			Filename: formatName("image.png"),
		}

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
			mock.AnythingOfType("string"),
			attachment.FileType,
		}

		mockFileRepository := new(mocks.FileRepository)

		mockFileRepository.
			On("UploadFile", uploadFileArgs...).
			Run(func(args mock.Arguments) {
				attachment.ID = id
			}).
			Return(imageURL, nil)

		ms := NewMessageService(&MSConfig{
			FileRepository: mockFileRepository,
		})

		_, err = ms.UploadFile(imageFileHeader, channelId)
		assert.NoError(t, err)

		mockFileRepository.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		channelId := fixture.RandID()

		multipartImageFixture := fixture.NewMultipartImage("image.png", "image/png")
		defer multipartImageFixture.Close()
		imageFileHeader := multipartImageFixture.GetFormFile()
		directory := fmt.Sprintf("channels/%s", channelId)

		attachment := model.Attachment{
			FileType: imageFileHeader.Header.Get("Content-Type"),
			Filename: formatName("image.png"),
		}

		uploadFileArgs := mock.Arguments{
			imageFileHeader,
			directory,
			mock.AnythingOfType("string"),
			attachment.FileType,
		}

		mockFileRepository := new(mocks.FileRepository)

		mockError := apperrors.NewInternal()
		mockFileRepository.
			On("UploadFile", uploadFileArgs...).
			Return("", mockError)

		ms := NewMessageService(&MSConfig{
			FileRepository: mockFileRepository,
		})

		att, err := ms.UploadFile(imageFileHeader, channelId)
		assert.Error(t, err)
		assert.Equal(t, err, apperrors.NewInternal())
		assert.Nil(t, att)

		mockFileRepository.AssertExpectations(t)
	})
}

package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

/*
 * MessageHandler contains all routes related to message actions (/api/messages)
 */

// GetMessages returns messages for the given channel
// It returns the most recent 35 or the ones after the given cursor
// GetMessages godoc
// @Tags Messages
// @Summary Get Channel Messages
// @Produce  json
// @Param channelId path string true "Channel ID"
// @Param cursor query string false "Cursor Pagination using the createdAt field"
// @Success 200 {array} model.MessageResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /messages/{channelId} [get]
func (h *Handler) GetMessages(c *gin.Context) {
	channelId := c.Param("channelId")
	userId := c.MustGet("userId").(string)

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user has access to said channel
	err = h.channelService.IsChannelMember(channel, userId)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Cursor is based on the created_at field of the message
	cursor := c.Query("cursor")

	messages, err := h.messageService.GetMessages(userId, channel, cursor)

	if err != nil {
		e := apperrors.NewNotFound("messages", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// If the channel does not have any messages, return an empty array
	if len(*messages) == 0 {
		var empty = make([]model.MessageResponse, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	c.JSON(http.StatusOK, messages)
}

// messageRequest contains all field required to create a message.
// Either text or file must be provided
type messageRequest struct {
	// Maximum 2000 characters
	Text *string `form:"text"`
	// image/* or audio/*
	File *multipart.FileHeader `form:"file" swaggertype:"string" format:"binary"`
} //@name MessageRequest

func (r messageRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Text,
			validation.NilOrNotEmpty,
			validation.Required.When(r.File == nil).
				Error(apperrors.MessageOrFileRequired),
			validation.Length(1, 2000),
		),
	)
}

func (r *messageRequest) sanitize() {
	if r.Text != nil {
		text := strings.TrimSpace(*r.Text)
		r.Text = &text
	}
}

// CreateMessage creates a message in the given channel
// CreateMessage godoc
// @Tags Messages
// @Summary Create Messages
// @Accepts  mpfd
// @Produce  json
// @Param channelId path string true "Channel ID"
// @Param request body messageRequest true "Create Message"
// @Success 201 {object} model.Success
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /messages/{channelId} [post]
func (h *Handler) CreateMessage(c *gin.Context) {
	channelId := c.Param("channelId")
	userId := c.MustGet("userId").(string)

	var req messageRequest
	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user has access to said channel
	err = h.channelService.IsChannelMember(channel, userId)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	author, err := h.userService.Get(userId)

	if err != nil {
		e := apperrors.NewNotFound("user", userId)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	params := model.Message{
		UserId:    userId,
		ChannelId: channel.ID,
	}

	params.Text = req.Text

	if req.File != nil {
		mimeType := req.File.Header.Get("Content-Type")

		if valid := isAllowedFileType(mimeType); !valid {
			toFieldErrorResponse(c, "File", apperrors.InvalidImageType)
			return
		}

		// Prevent file upload on the live server.
		// Remove the if part if you do want upload
		var attachment *model.Attachment
		if gin.Mode() == gin.ReleaseMode {
			id, _ := gonanoid.Nanoid(20)

			// Random image to test files in the app
			attachment = &model.Attachment{
				ID:       id,
				Url:      fmt.Sprintf("https://picsum.photos/seed/%s/600", id),
				FileType: "image/jpeg",
				Filename: id,
			}
		} else {
			attachment, err = h.messageService.UploadFile(req.File, channel.ID)

			if err != nil {
				c.JSON(apperrors.Status(err), gin.H{
					"error": err,
				})
				return
			}
		}

		params.Attachment = attachment
	}

	message, err := h.messageService.CreateMessage(&params)

	if err != nil {
		log.Printf("Failed to create message: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	response := model.MessageResponse{
		Id:         message.ID,
		Text:       message.Text,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
		Attachment: message.Attachment,
		User: model.MemberResponse{
			Id:        author.ID,
			Username:  author.Username,
			Image:     author.Image,
			IsOnline:  author.IsOnline,
			CreatedAt: author.CreatedAt,
			UpdatedAt: author.UpdatedAt,
			IsFriend:  false,
		},
	}

	// Get member settings if it is not a DM
	if !channel.IsDM {
		settings, _ := h.guildService.GetMemberSettings(userId, *channel.GuildID)
		response.User.Nickname = settings.Nickname
		response.User.Color = settings.Color
	}

	// Emit new message to the channel
	h.socketService.EmitNewMessage(channelId, &response)

	if channel.IsDM {
		// Open the DM and push it to the top
		_ = h.channelService.OpenDMForAll(channelId)
		// Post a notification
		h.socketService.EmitNewDMNotification(channelId, author)
	} else {
		// Update last activity in channel
		channel.LastActivity = time.Now()
		_ = h.channelService.UpdateChannel(channel)
		// Post a notification
		h.socketService.EmitNewNotification(*channel.GuildID, channelId)
	}

	c.JSON(http.StatusCreated, true)
}

// EditMessage edits the given message with the given text
// EditMessage godoc
// @Tags Messages
// @Summary Edit Messages
// @Accepts  json
// @Produce  json
// @Param messageId path string true "Message ID"
// @Param request body messageRequest true "Edit Message"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /messages/{messageId} [put]
func (h *Handler) EditMessage(c *gin.Context) {
	messageId := c.Param("messageId")
	userId := c.MustGet("userId").(string)

	var req messageRequest
	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	message, err := h.messageService.Get(messageId)

	if err != nil {
		e := apperrors.NewNotFound("message", messageId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if message.UserId != userId {
		e := apperrors.NewAuthorization(apperrors.EditMessageError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	message.Text = req.Text

	if err = h.messageService.UpdateMessage(message); err != nil {
		log.Printf("Failed to edit message: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	response := model.MessageResponse{
		Id:         message.ID,
		Text:       message.Text,
		CreatedAt:  message.CreatedAt,
		UpdatedAt:  message.UpdatedAt,
		Attachment: message.Attachment,
		User: model.MemberResponse{
			Id: userId,
		},
	}

	// Emit edited message to the channel
	h.socketService.EmitEditMessage(message.ChannelId, &response)

	c.JSON(http.StatusOK, true)
}

// DeleteMessage deletes the given message
// DeleteMessage godoc
// @Tags Messages
// @Summary Delete Messages
// @Produce  json
// @Param messageId path string true "Message ID"
// @Success 200 {object} model.Success
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /messages/{messageId} [delete]
func (h *Handler) DeleteMessage(c *gin.Context) {
	messageId := c.Param("messageId")
	userId := c.MustGet("userId").(string)
	message, err := h.messageService.Get(messageId)

	if err != nil {
		e := apperrors.NewNotFound("message", messageId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	channel, err := h.channelService.Get(message.ChannelId)

	if err != nil {
		e := apperrors.NewNotFound("message", messageId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if message author or guild owner
	if !channel.IsDM {
		guild, err := h.guildService.GetGuild(*channel.GuildID)

		if err != nil {
			e := apperrors.NewNotFound("message", messageId)

			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		if message.UserId != userId && guild.OwnerId != userId {
			e := apperrors.NewAuthorization(apperrors.DeleteMessageError)
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}
		// Only message author check required
	} else {
		if message.UserId != userId {
			e := apperrors.NewAuthorization(apperrors.DeleteDMMessageError)
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}
	}

	if err = h.messageService.DeleteMessage(message); err != nil {
		log.Printf("Failed to delete message: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit delete message to the channel
	h.socketService.EmitDeleteMessage(message.ChannelId, message.ID)

	c.JSON(http.StatusOK, true)
}

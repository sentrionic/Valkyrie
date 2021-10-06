package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
	"strings"
)

/*
 * ChannelHandler contains all routes related to channel actions (/api/channels)
 */

// GuildChannels returns the given guild's channels
// GuildChannels godoc
// @Tags Channels
// @Summary Get Guild Channels
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {array} model.ChannelResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /channels/{guildId} [get]
func (h *Handler) GuildChannels(c *gin.Context) {
	guildId := c.Param("id")
	userId := c.MustGet("userId").(string)

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Only get the channels if the user is a member
	if !isMember(guild, userId) {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	channels, err := h.channelService.GetChannels(userId, guildId)

	if err != nil {
		log.Printf("Unable to find channels for guild id: %v\n%v", guildId, err)
		e := apperrors.NewNotFound("channels", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, channels)
}

// channelReq specifies the input form for creating a channel
// IsPublic and Members do not need to be specified if you want
// to create a public channel
type channelReq struct {
	// Channel Name. 3 to 30 character
	Name string `json:"name"`
	// Default is true
	IsPublic *bool `json:"isPublic"`
	// Array of memberIds
	Members []string `json:"members"`
} //@name ChannelRequest

func (r channelReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 30)),
	)
}

func (r *channelReq) sanitize() {
	r.Name = strings.TrimSpace(r.Name)
}

// CreateChannel creates a channel for the given guild param
// CreateChannel godoc
// @Tags Channels
// @Summary Create Channel
// @Accepts json
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {array} model.ChannelResponse
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /channels/{guildId} [post]
func (h *Handler) CreateChannel(c *gin.Context) {
	var req channelReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	userId := c.MustGet("userId").(string)
	guildId := c.Param("id")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the server already has 50 channels
	if len(guild.Channels) >= model.MaximumChannels {
		e := apperrors.NewBadRequest(apperrors.ChannelLimitError)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	channelParams := model.Channel{
		Name:     req.Name,
		IsPublic: true,
		GuildID:  &guildId,
	}

	// Channel is private
	if req.IsPublic != nil && !*req.IsPublic {
		channelParams.IsPublic = false

		// Add the current user to the members if they are not in there
		if !containsUser(req.Members, userId) {
			req.Members = append(req.Members, userId)
		}
		members, err := h.guildService.FindUsersByIds(req.Members, guildId)

		if err != nil {
			e := apperrors.NewInternal()
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		// Create private channel members
		channelParams.PCMembers = append(channelParams.PCMembers, *members...)
	}

	channel, err := h.channelService.CreateChannel(&channelParams)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	guild.Channels = append(guild.Channels, *channel)

	if err = h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to update guild: %v\n", err.Error())
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	response := channel.SerializeChannel()

	// Emit the new channel to the guild members
	if channel.IsPublic {
		h.socketService.EmitNewChannel(guildId, &response)
		// Emit to private channel members
	} else {
		h.socketService.EmitNewPrivateChannel(req.Members, &response)
	}

	c.JSON(http.StatusCreated, response)
}

// PrivateChannelMembers returns the ids of all members
// that are part of the channel
// PrivateChannelMembers godoc
// @Tags Channels
// @Summary Get Members of the given Channel
// @Produce  json
// @Param channelId path string true "Channel ID"
// @Success 200 {array} string
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /channels/{channelId}/members [get]
func (h *Handler) PrivateChannelMembers(c *gin.Context) {
	channelId := c.Param("id")
	userId := c.MustGet("userId").(string)

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if channel.GuildID == nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(*channel.GuildID)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Public channels do not have any private members
	if channel.IsPublic {
		var empty = make([]string, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	members, err := h.channelService.GetPrivateChannelMembers(channelId)

	if err != nil {
		log.Printf("Unable to find members for channel: %v\n%v", channelId, err)
		e := apperrors.NewNotFound("members", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, members)
}

// DirectMessages returns a list of the current users DMs
// DirectMessages godoc
// @Tags Channels
// @Summary Get User's DMs
// @Produce  json
// @Success 200 {array} model.DirectMessage
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /channels/me/dm [get]
func (h *Handler) DirectMessages(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	channels, err := h.channelService.GetDirectMessages(userId)

	if err != nil {
		log.Printf("Unable to find dms for user id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("dms", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// If the user does not have any dms, return an empty array
	if len(*channels) == 0 {
		var empty = make([]model.DirectMessage, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	c.JSON(http.StatusOK, channels)
}

// GetOrCreateDM gets the DM with the given member and creates it
// if it does not already exist
// DirectMessages godoc
// @Tags Channels
// @Summary Get or Create DM
// @Produce  json
// @Param channelId path string true "Member ID"
// @Success 200 {object} model.DirectMessage
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /channels/{channelId}/dm [post]
func (h *Handler) GetOrCreateDM(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	memberId := c.Param("id")

	if userId == memberId {
		e := apperrors.NewBadRequest(apperrors.DMYourselfError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	member, err := h.friendService.GetMemberById(memberId)

	if err != nil {
		log.Printf("Unable to find member for id: %v\n%v", memberId, err)
		e := apperrors.NewNotFound("member", memberId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// check if dm channel already exists with these members
	dmId, err := h.channelService.GetDirectMessageChannel(userId, memberId)

	if err != nil {
		log.Printf("Unable to find or create dms for user id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("dms", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// dm already exists
	if dmId != nil && *dmId != "" {
		_ = h.channelService.SetDirectMessageStatus(*dmId, userId, true)
		c.JSON(http.StatusOK, toDMChannel(member, *dmId, userId))
		return
	}

	// Create the dm channel between the current user and the member
	id := fmt.Sprintf("%s-%s", userId, memberId)
	channelParams := model.Channel{
		Name:     id,
		IsPublic: false,
		IsDM:     true,
	}

	channel, err := h.channelService.CreateChannel(&channelParams)

	// Create the DM channel
	if err != nil {
		log.Printf("Failed to create channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Add the users to it
	ids := []string{userId, memberId}
	err = h.channelService.AddDMChannelMembers(ids, channel.ID, userId)

	if err != nil {
		log.Printf("Failed to create channel: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, toDMChannel(member, channel.ID, userId))
}

// toDMChannel returns the DM response for the given channel and member
func toDMChannel(member *model.User, channelId string, userId string) model.DirectMessage {
	return model.DirectMessage{
		Id: channelId,
		User: model.DMUser{
			Id:       member.ID,
			Username: member.Username,
			Image:    member.Image,
			IsOnline: member.IsOnline,
			IsFriend: isFriend(member, userId),
		},
	}
}

// EditChannel edits the specified channel
// EditChannel godoc
// @Tags Channels
// @Summary Edit Channel
// @Accepts json
// @Produce  json
// @Param channelId path string true "Channel ID"
// @Param request body channelReq true "Edit Channel"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /channels/{channelId} [put]
func (h *Handler) EditChannel(c *gin.Context) {
	var req channelReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	userId := c.MustGet("userId").(string)
	channelId := c.Param("id")

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(*channel.GuildID)

	if err != nil {
		e := apperrors.NewNotFound("guild", *channel.GuildID)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	isPublic := true
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}

	// Used to be private and now is public
	if isPublic && !channel.IsPublic {
		err = h.channelService.CleanPCMembers(channelId)
		if err != nil {
			log.Printf("error removing pc members: %v", err)
		}
	}

	channel.IsPublic = isPublic
	channel.Name = req.Name

	// Member Changes
	if !isPublic {
		// Check if the array contains the current member
		if !containsUser(req.Members, userId) {
			req.Members = append(req.Members, userId)
		}

		// Current members of the channel
		current := make([]string, 0)
		for _, member := range channel.PCMembers {
			current = append(current, member.ID)
		}

		// Newly added members
		newMembers := difference(req.Members, current)
		// Members that got removed
		toRemove := difference(current, req.Members)

		err = h.channelService.AddPrivateChannelMembers(newMembers, channelId)
		if err != nil {
			c.JSON(apperrors.Status(err), gin.H{
				"error": err,
			})
			return
		}

		err = h.channelService.RemovePrivateChannelMembers(toRemove, channelId)
		if err != nil {
			c.JSON(apperrors.Status(err), gin.H{
				"error": err,
			})
			return
		}
	}

	if err = h.channelService.UpdateChannel(channel); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit the channel changes to the guild members
	response := channel.SerializeChannel()
	h.socketService.EmitEditChannel(*channel.GuildID, &response)

	c.JSON(http.StatusOK, true)
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// DeleteChannel removes the given channel from the guild
// DeleteChannel godoc
// @Tags Channels
// @Summary Delete Channel
// @Produce  json
// @Param id path string true "Channel ID"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /channels/{id} [delete]
func (h *Handler) DeleteChannel(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	channelId := c.Param("id")

	channel, err := h.channelService.Get(channelId)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(*channel.GuildID)

	if err != nil {
		e := apperrors.NewNotFound("channel", channelId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the guild has the minimum amount of channels
	if len(guild.Channels) <= model.MinimumChannels {
		e := apperrors.NewBadRequest(apperrors.OneChannelRequired)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if err = h.channelService.DeleteChannel(channel); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signal to remove the channel from the guild
	h.socketService.EmitDeleteChannel(channel)

	c.JSON(http.StatusOK, true)
}

// CloseDM closes the DM on the current users side
// CloseDM godoc
// @Tags Channels
// @Summary Close DM
// @Produce  json
// @Param id path string true "DM Channel ID"
// @Success 200 {object} model.Success
// @Failure 404 {object} model.ErrorResponse
// @Router /channels/{id}/dm [delete]
func (h *Handler) CloseDM(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	channelId := c.Param("id")

	dmId, err := h.channelService.GetDMByUserAndChannel(userId, channelId)

	if err != nil || dmId == "" {
		log.Printf("Unable to find or create dms for user id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("dms", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	_ = h.channelService.SetDirectMessageStatus(channelId, userId, false)

	c.JSON(http.StatusOK, true)
}

// containsUser checks if the array contains the user
func containsUser(members []string, userId string) bool {
	for _, m := range members {
		if m == userId {
			return true
		}
	}
	return false
}

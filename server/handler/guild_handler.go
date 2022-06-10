package handler

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/lib/pq"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
)

/*
 * GuildHandler contains all routes related to guild actions (/api/guilds)
 */

// GetUserGuilds returns the current users guilds
// GetUserGuilds godoc
// @Tags Guilds
// @Summary Get Current User's Guilds
// @Produce  json
// @Success 200 {array} model.GuildResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /guilds [get]
func (h *Handler) GetUserGuilds(c *gin.Context) {
	userId := c.MustGet("userId").(string)

	guilds, err := h.guildService.GetUserGuilds(userId)

	if err != nil {
		log.Printf("Unable to find guilds for id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, guilds)
}

// GetGuildMembers returns the given guild's members
// GetGuildMembers godoc
// @Tags Guilds
// @Summary Get Guild Members
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {array} model.MemberResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /guilds/{guildId}/members [get]
func (h *Handler) GetGuildMembers(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		log.Printf("Unable to find guilds for id: %v\n%v", guildId, err)
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if a member
	if !isMember(guild, userId) {
		e := apperrors.NewAuthorization(apperrors.NotAMember)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	members, err := h.guildService.GetGuildMembers(userId, guildId)

	if err != nil {
		log.Printf("Unable to find guilds for id: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, members)
}

// GetVCMembers returns the given guild's members that are currently in the VC
// GetVCMembers godoc
// @Tags Guilds
// @Summary Get Guild VC Members
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {array} model.VCMemberResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /guilds/{guildId}/vcmembers [get]
func (h *Handler) GetVCMembers(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		log.Printf("Unable to find guilds for id: %v\n%v", guildId, err)
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if a member
	if !isMember(guild, userId) {
		e := apperrors.NewAuthorization(apperrors.NotAMember)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	members, err := h.guildService.GetVCMembers(guild.ID)

	if err != nil {
		log.Printf("Unable to find vc members for id: %v\n%v", guildId, err)
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, members)
}

type createGuildRequest struct {
	// Guild Name. 3 to 30 characters
	Name string `json:"name"`
} //@name CreateGuildRequest

func (r createGuildRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 30)),
	)
}

func (r *createGuildRequest) sanitize() {
	r.Name = strings.TrimSpace(r.Name)
}

// CreateGuild creates a guild
// CreateGuild godoc
// @Tags Guilds
// @Summary Create Guild
// @Accepts  json
// @Produce  json
// @Param request body createGuildRequest true "Create Guild"
// @Success 201 {array} model.GuildResponse
// @Failure 400 {object} model.ErrorsResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/create [post]
func (h *Handler) CreateGuild(c *gin.Context) {
	var req createGuildRequest

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	userId := c.MustGet("userId").(string)

	authUser, err := h.guildService.GetUser(userId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user is already in 100 guilds
	if len(authUser.Guilds) >= model.MaximumGuilds {
		e := apperrors.NewBadRequest(apperrors.GuildLimitReached)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guildParams := model.Guild{
		Name:    req.Name,
		OwnerId: userId,
	}

	// Add the current user as a member
	guildParams.Members = append(guildParams.Members, *authUser)

	guild, err := h.guildService.CreateGuild(&guildParams)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Create the default 'general' channel for the guild
	channelParams := model.Channel{
		GuildID:  &guild.ID,
		Name:     "general",
		IsPublic: true,
	}

	channel, err := h.channelService.CreateChannel(&channelParams)

	if err != nil {
		log.Printf("Failed to create channel for guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusCreated, guild.SerializeGuild(channel.ID))
}

// editGuildRequest specifies the form to edit the guild.
// If Image is not nil then the guild's icon got changed.
// If Icon is not nil then the guild kept its old one.
// If both are nil then the icon got reset.
type editGuildRequest struct {
	// Guild Name. 3 to 30 characters
	Name string `form:"name"`
	// image/png or image/jpeg
	Image *multipart.FileHeader `form:"image" swaggertype:"string" format:"binary"`
	// The old guild icon url if no new image is selected. Set to null to reset the guild icon
	Icon *string `form:"icon"`
} //@name EditGuildRequest

func (r editGuildRequest) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required, validation.Length(3, 30)),
	)
}

func (r *editGuildRequest) sanitize() {
	r.Name = strings.TrimSpace(r.Name)
}

// EditGuild edits the given guild
// EditGuild godoc
// @Tags Guilds
// @Summary Edit Guild
// @Accepts  mpfd
// @Produce  json
// @Param request body editGuildRequest true "Edit Guild"
// @Param guildId path string true "Guild ID"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.ErrorsResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId} [put]
func (h *Handler) EditGuild(c *gin.Context) {
	var req editGuildRequest

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

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

	guild.Name = req.Name

	// Guild icon got changed
	if req.Image != nil {
		// Validate image mime-type is allowable
		mimeType := req.Image.Header.Get("Content-Type")

		if valid := isAllowedImageType(mimeType); !valid {
			toFieldErrorResponse(c, "Image", apperrors.InvalidImageType)
			return
		}

		directory := fmt.Sprintf("valkyrie/guilds/%s", guild.ID)
		url, err := h.userService.ChangeAvatar(req.Image, directory)

		if err != nil {
			e := apperrors.NewInternal()
			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}

		if guild.Icon != nil {
			_ = h.userService.DeleteImage(*guild.Icon)
		}
		guild.Icon = &url
		// Guild kept its old icon
	} else if req.Icon != nil {
		guild.Icon = req.Icon
		// Guild reset its icon
	} else {
		guild.Icon = nil
	}

	if err = h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to update guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit guild changes to guild members
	h.socketService.EmitEditGuild(guild)

	c.JSON(http.StatusOK, true)
}

// GetInvite creates an invitation for the given guild
// The isPermanent query parameter specifies if the invite
// should not be deleted after it got used
// GetInvite godoc
// @Tags Guilds
// @Summary Get Guild Invite
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Param isPermanent query boolean false "Is Permanent"
// @Success 200 string link
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/invite [get]
func (h *Handler) GetInvite(c *gin.Context) {
	guildId := c.Param("guildId")
	permanent := c.Query("isPermanent")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)
	// Must be a member to create an invitation
	if !isMember(guild, userId) {
		e := apperrors.NewAuthorization(apperrors.MustBeMemberInvite)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	isPermanent := false
	if permanent != "" {
		isPermanent, err = strconv.ParseBool(permanent)

		if err != nil {
			e := apperrors.NewBadRequest(apperrors.IsPermanentError)

			c.JSON(e.Status(), gin.H{
				"error": e,
			})
			return
		}
	}

	ctx := context.Background()
	link, err := h.guildService.GenerateInviteLink(ctx, guild.ID, isPermanent)

	if err != nil {
		e := apperrors.NewInternal()
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if isPermanent {
		guild.InviteLinks = append(guild.InviteLinks, link)
		_ = h.guildService.UpdateGuild(guild)
	}

	origin := os.Getenv("CORS_ORIGIN")
	c.JSON(http.StatusOK, fmt.Sprintf("%s/%s", origin, link))
}

// DeleteGuildInvites removes all permanent invites from the given guild
// DeleteGuildInvites godoc
// @Tags Guilds
// @Summary Delete all permanent invite links
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {object} model.Success
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/invite [delete]
func (h *Handler) DeleteGuildInvites(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.InvalidateInvitesError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	ctx := context.Background()
	h.guildService.InvalidateInvites(ctx, guild)
	guild.InviteLinks = make(pq.StringArray, 0)

	if err = h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to delete guild invites: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

type joinReq struct {
	Link string `json:"link"`
} //@name JoinRequest

func (r joinReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Link, validation.Required),
	)
}

func (r *joinReq) sanitize() {
	r.Link = strings.TrimSpace(r.Link)
}

// JoinGuild adds the current user to invited guild
// JoinGuild godoc
// @Tags Guilds
// @Summary Join Guild
// @Produce  json
// @Param request body joinReq true "Join Guild"
// @Success 200 {object} model.GuildResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/join [post]
func (h *Handler) JoinGuild(c *gin.Context) {
	var req joinReq

	// Bind incoming json to struct and check for validation errors
	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	userId := c.MustGet("userId").(string)

	authUser, err := h.guildService.GetUser(userId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user has reached the guild limit
	if len(authUser.Guilds) >= model.MaximumGuilds {
		e := apperrors.NewBadRequest(apperrors.GuildLimitReached)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// If the link contains the domain, remove it
	if strings.Contains(req.Link, "/") {
		req.Link = req.Link[strings.LastIndex(req.Link, "/")+1:]
	}

	ctx := context.Background()
	guildId, err := h.guildService.GetGuildIdFromInvite(ctx, req.Link)

	if err != nil {
		e := apperrors.NewBadRequest(apperrors.InvalidInviteError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewBadRequest(apperrors.InvalidInviteError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user is banned from the guild
	if isBanned(guild, authUser.ID) {
		e := apperrors.NewBadRequest(apperrors.BannedFromServer)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Check if the user is already a member
	if isMember(guild, authUser.ID) {
		e := apperrors.NewBadRequest(apperrors.AlreadyMember)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild.Members = append(guild.Members, *authUser)

	if err = h.guildService.UpdateGuild(guild); err != nil {
		log.Printf("Failed to join guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit new member to the guild
	h.socketService.EmitAddMember(guild.ID, authUser)

	channel, _ := h.guildService.GetDefaultChannel(guildId)

	c.JSON(http.StatusCreated, guild.SerializeGuild(channel.ID))
}

// LeaveGuild leaves the given guild
// LeaveGuild godoc
// @Tags Guilds
// @Summary Leave Guild
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {object} model.Success
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId} [delete]
func (h *Handler) LeaveGuild(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId == userId {
		e := apperrors.NewAuthorization(apperrors.OwnerCantLeave)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if err := h.guildService.RemoveMember(userId, guildId); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signal to remove the member from the guild
	h.socketService.EmitRemoveMember(guild.ID, userId)

	c.JSON(http.StatusOK, true)
}

// DeleteGuild deletes the given guild
// DeleteGuild godoc
// @Tags Guilds
// @Summary Delete Guild
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {object} model.Success
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/delete [delete]
func (h *Handler) DeleteGuild(c *gin.Context) {
	userId := c.MustGet("userId").(string)
	guildId := c.Param("guildId")

	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.DeleteGuildError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	// Get the ID of all members to emit the deletion to
	members := make([]string, 0)
	for _, member := range guild.Members {
		members = append(members, member.ID)
	}

	if err := h.guildService.DeleteGuild(guildId); err != nil {
		log.Printf("Failed to leave guild: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signal to remove the guild to its members
	h.socketService.EmitDeleteGuild(guildId, members)

	c.JSON(http.StatusOK, true)
}

// isMember checks if the given user is a member of the guild
func isMember(guild *model.Guild, userId string) bool {
	for _, v := range guild.Members {
		if v.ID == userId {
			return true
		}
	}
	return false
}

// isBanned checks if the given user is banned from the guild
func isBanned(guild *model.Guild, userId string) bool {
	for _, v := range guild.Bans {
		if v.ID == userId {
			return true
		}
	}
	return false
}

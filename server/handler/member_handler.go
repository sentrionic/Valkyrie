package handler

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"net/http"
	"strings"
)

/*
 * MemberHandler contains all routes related to mod actions (/api/guilds)
 */

// memberReq contains the MemberId of the user
// that needs to be moderated
type memberReq struct {
	MemberId string `json:"memberId"`
} //@name MemberRequest

func (r memberReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.MemberId, validation.Required, is.UTFDigit),
	)
}

// GetMemberSettings gets the current user's role color and nickname
// for the given guild
// GetMemberSettings godoc
// @Tags Members
// @Summary Get Member Settings
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {object} model.MemberSettings
// @Failure 404 {object} model.ErrorResponse
// @Router /guilds/{guildId}/member [get]
func (h *Handler) GetMemberSettings(c *gin.Context) {
	guildId := c.Param("guildId")
	userId := c.MustGet("userId").(string)
	_, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	settings, err := h.guildService.GetMemberSettings(userId, guildId)

	if err != nil {
		log.Printf("Unable to find settings: %v\n%v", userId, err)
		e := apperrors.NewNotFound("user", userId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, settings)
}

type memberSettingsReq struct {
	Nickname *string `json:"nickname"`
	Color    *string `json:"color"`
} //@name MemberSettingsRequest

func (r memberSettingsReq) validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Nickname, validation.NilOrNotEmpty, validation.Length(3, 30)),
		validation.Field(&r.Color, validation.NilOrNotEmpty, is.HexColor),
	)
}

func (r *memberSettingsReq) sanitize() {
	if r.Nickname != nil {
		nickname := strings.TrimSpace(*r.Nickname)
		r.Nickname = &nickname
	}
}

// EditMemberSettings changes the current user's role color and nickname
// for the given guild
// EditMemberSettings godoc
// @Tags Members
// @Summary Edit Member Settings
// @Accepts json
// @Produce  json
// @Param request body memberSettingsReq true "Edit Member"
// @Param guildId path string true "Guild ID"
// @Success 200 {object} model.Success
// @Failure 400 {object} model.ErrorsResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/member [put]
func (h *Handler) EditMemberSettings(c *gin.Context) {
	var req memberSettingsReq

	if ok := bindData(c, &req); !ok {
		return
	}

	req.sanitize()

	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	// Check if the user is a member of the guild
	if !isMember(guild, userId) {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	settings := &model.MemberSettings{
		Nickname: req.Nickname,
		Color:    req.Color,
	}

	err = h.guildService.UpdateMemberSettings(settings, userId, guildId)

	if err != nil {
		log.Printf("Unable to update settings for user: %v\n%v", userId, err)
		e := apperrors.NewInternal()

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

// GetBanList returns a list of all banned users for the given guild
// GetBanList godoc
// @Tags Members
// @Summary Get Guild Ban list
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Success 200 {array} model.BanResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/bans [get]
func (h *Handler) GetBanList(c *gin.Context) {
	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	bans, err := h.guildService.GetBanList(guildId)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// If the guild does not have any bans, return an empty array
	if len(*bans) == 0 {
		empty := make([]model.BanResponse, 0)
		c.JSON(http.StatusOK, empty)
		return
	}

	c.JSON(http.StatusOK, bans)
}

// BanMember bans the provided member from the given guild
// BanMember godoc
// @Tags Members
// @Summary Ban Member
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Param request body memberReq true "Member ID"
// @Success 200 {array} model.Success
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/bans [post]
func (h *Handler) BanMember(c *gin.Context) {
	var req memberReq

	if ok := bindData(c, &req); !ok {
		return
	}

	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	member, err := h.guildService.GetUser(req.MemberId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", req.MemberId, err)
		e := apperrors.NewNotFound("user", req.MemberId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if member.ID == userId {
		e := apperrors.NewBadRequest(apperrors.BanYourselfError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	guild.Bans = append(guild.Bans, *member)

	if err = h.guildService.UpdateGuild(guild); err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	err = h.guildService.RemoveMember(req.MemberId, guildId)

	if err != nil {
		log.Printf("Failed to ban member: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signals to remove the member from the guild
	h.socketService.EmitRemoveMember(guild.ID, member.ID)
	h.socketService.EmitRemoveFromGuild(member.ID, guildId)

	c.JSON(http.StatusOK, true)
}

// UnbanMember unbans the specified user from the given guild
// BanMember godoc
// @Tags Members
// @Summary Unban Member
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Param request body memberReq true "Member ID"
// @Success 200 {array} model.Success
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/bans [delete]
func (h *Handler) UnbanMember(c *gin.Context) {
	var req memberReq

	if ok := bindData(c, &req); !ok {
		return
	}

	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if req.MemberId == userId {
		e := apperrors.NewBadRequest(apperrors.UnbanYourselfError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if err := h.guildService.UnbanMember(req.MemberId, guildId); err != nil {
		log.Printf("Failed to unban member: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, true)
}

// KickMember kicks the provided member from the given guild
// KickMember godoc
// @Tags Members
// @Summary Kick Member
// @Produce  json
// @Param guildId path string true "Guild ID"
// @Param request body memberReq true "Member ID"
// @Success 200 {array} model.Success
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /guilds/{guildId}/kick [post]
func (h *Handler) KickMember(c *gin.Context) {
	var req memberReq

	if ok := bindData(c, &req); !ok {
		return
	}

	guildId := c.Param("guildId")
	guild, err := h.guildService.GetGuild(guildId)

	if err != nil {
		e := apperrors.NewNotFound("guild", guildId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	userId := c.MustGet("userId").(string)

	if guild.OwnerId != userId {
		e := apperrors.NewAuthorization(apperrors.MustBeOwner)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	member, err := h.guildService.GetUser(req.MemberId)

	if err != nil {
		log.Printf("Unable to find user: %v\n%v", req.MemberId, err)
		e := apperrors.NewNotFound("user", req.MemberId)

		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	if member.ID == userId {
		e := apperrors.NewBadRequest(apperrors.KickYourselfError)
		c.JSON(e.Status(), gin.H{
			"error": e,
		})
		return
	}

	err = h.guildService.RemoveMember(req.MemberId, guildId)

	if err != nil {
		log.Printf("Failed to kick member: %v\n", err.Error())
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
		return
	}

	// Emit signals to remove the member from the guild
	h.socketService.EmitRemoveMember(guild.ID, member.ID)
	h.socketService.EmitRemoveFromGuild(member.ID, guildId)

	c.JSON(http.StatusOK, true)
}

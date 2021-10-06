package model

// Invite represents an invite link for a guild
// IsPermanent indicates if the invite should not expire
type Invite struct {
	GuildId     string `json:"guild_id"`
	IsPermanent bool   `json:"is_permanent"`
}

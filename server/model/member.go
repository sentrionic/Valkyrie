package model

import "time"

// Member represents a user in a guild and is the join table between
// User and Guild.
type Member struct {
	UserID    string    `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	GuildID   string    `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	Nickname  *string   `gorm:"nickname"`
	Color     *string   `gorm:"color"`
	LastSeen  time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

type VCMember struct {
	UserID     string `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	GuildID    string `gorm:"primaryKey;constraint:OnDelete:CASCADE;"`
	IsMuted    bool
	IsDeafened bool
}

// MemberResponse is the API response of a member.
type MemberResponse struct {
	Id        string    `json:"id"`
	Username  string    `json:"username"`
	Image     string    `json:"image"`
	IsOnline  bool      `json:"isOnline"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Nickname  *string   `json:"nickname"`
	Color     *string   `json:"color"`
	IsFriend  bool      `json:"isFriend"`
} //@name Member

// BanResponse is the API response of a banned member.
type BanResponse struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
} //@name BanResponse

// MemberSettings is the API response of a member's guild settings.
type MemberSettings struct {
	Nickname *string `json:"nickname"`
	Color    *string `json:"color"`
} //@name MemberSettings

// VCMemberResponse is the API response of a member that is currently in a VC.
type VCMemberResponse struct {
	Id         string  `json:"id"`
	Username   string  `json:"username"`
	Image      string  `json:"image"`
	IsMuted    bool    `json:"isMuted"`
	IsDeafened bool    `json:"IsDeafened"`
	Nickname   *string `json:"nickname"`
} //@name VCMember

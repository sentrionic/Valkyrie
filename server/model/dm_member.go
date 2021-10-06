package model

import "time"

// DMMember represents a member of a DM channel.
// IsOpen indicates if the DM is open on the client.
type DMMember struct {
	ID        string `gorm:"primaryKey"`
	UserID    string `gorm:"primaryKey"`
	ChannelId string `gorm:"primaryKey;"`
	IsOpen    bool
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
}

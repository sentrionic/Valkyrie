package model

import (
	"context"
	"mime/multipart"
)

// User represents the user of the website.
type User struct {
	BaseModel
	Username string    `gorm:"not null" json:"username"`
	Email    string    `gorm:"not null;uniqueIndex" json:"email"`
	Password string    `gorm:"not null" json:"-"`
	Image    string    `json:"image"`
	IsOnline bool      `gorm:"index;default:true" json:"isOnline"`
	Friends  []User    `gorm:"many2many:friends;" json:"-"`
	Requests []User    `gorm:"many2many:friend_requests;joinForeignKey:sender_id;joinReferences:receiver_id" json:"-"`
	Guilds   []Guild   `gorm:"many2many:members;" json:"-"`
	Message  []Message `json:"-"`
} //@name User

// UserService defines methods related to account operations the handler layer expects
// any service it interacts with to implement
type UserService interface {
	Get(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	Register(user *User) (*User, error)
	Login(email, password string) (*User, error)
	UpdateAccount(user *User) error
	IsEmailAlreadyInUse(email string) bool
	ChangeAvatar(header *multipart.FileHeader, directory string) (string, error)
	DeleteImage(key string) error
	ChangePassword(currentPassword, newPassword string, user *User) error
	ForgotPassword(ctx context.Context, user *User) error
	ResetPassword(ctx context.Context, password string, token string) (*User, error)
	GetFriendAndGuildIds(userId string) (*[]string, error)
	GetRequestCount(userId string) (*int64, error)
}

// UserRepository defines methods related to account db operations the service layer expects
// any repository it interacts with to implement
type UserRepository interface {
	FindByID(id string) (*User, error)
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) error
	GetFriendAndGuildIds(userId string) (*[]string, error)
	GetRequestCount(userId string) (*int64, error)
}

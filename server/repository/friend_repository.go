package repository

import (
	"database/sql"
	"errors"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
)

// friendRepository is data/repository implementation
// of service layer FriendRepository
type friendRepository struct {
	DB *gorm.DB
}

// NewFriendRepository is a factory for initializing Friend Repositories
func NewFriendRepository(db *gorm.DB) model.FriendRepository {
	return &friendRepository{
		DB: db,
	}
}

// FriendsList returns the friends for the given user ID from the DB
func (r *friendRepository) FriendsList(id string) (*[]model.Friend, error) {
	var friends []model.Friend

	result := r.DB.
		Table("users").
		Joins(`JOIN friends ON friends.user_id = "users".id`).
		Where("friends.friend_id = ?", id).
		Find(&friends)

	return &friends, result.Error
}

// RequestList returns the friend requests for the given user ID from the DB
func (r *friendRepository) RequestList(id string) (*[]model.FriendRequest, error) {
	var requests []model.FriendRequest

	result := r.DB.
		Raw(`
		  select u.id, u.username, u.image, 1 as "type" from users u
		  join friend_requests fr on u.id = fr."sender_id"
		  where fr."receiver_id" = @id
		  UNION
		  select u.id, u.username, u.image, 0 as "type" from users u
		  join friend_requests fr on u.id = fr."receiver_id"
		  where fr."sender_id" = @id
		  order by username;
		`, sql.Named("id", id)).
		Find(&requests)

	return &requests, result.Error
}

// FindByID returns a User containing their friends and requests from the DB
func (r *friendRepository) FindByID(id string) (*model.User, error) {
	user := &model.User{}

	if err := r.DB.
		Preload("Friends").
		Preload("Requests").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", id)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

// DeleteRequest removes the given member and user from the friend requests
func (r *friendRepository) DeleteRequest(memberId string, userId string) error {
	return r.DB.Exec(`
		DELETE
		FROM friend_requests
		WHERE receiver_id = @memberId AND sender_id = @userId
		   OR receiver_id = @userId AND sender_id = @memberId
`, sql.Named("memberId", memberId), sql.Named("userId", userId)).Error
}

// RemoveFriend removes members from the friends table.
func (r *friendRepository) RemoveFriend(memberId string, userId string) error {
	return r.DB.
		Exec("DELETE FROM friends WHERE user_id = ? AND friend_id = ?", memberId, userId).
		Exec("DELETE FROM friends WHERE user_id = ? AND friend_id = ?", userId, memberId).
		Error
}

// Save inserts the given user in the DB
func (r *friendRepository) Save(user *model.User) error {
	return r.DB.Save(&user).Error
}

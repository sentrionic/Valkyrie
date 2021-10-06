package repository

import (
	"database/sql"
	"errors"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
	"log"
	"regexp"
)

// userRepository is data/repository implementation
// of service layer UserRepository
type userRepository struct {
	DB *gorm.DB
}

// NewUserRepository is a factory for initializing User Repositories
func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		DB: db,
	}
}

// FindByID returns a user for the given ID
func (r *userRepository) FindByID(id string) (*model.User, error) {
	user := &model.User{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", id)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

// Create inserts the user in the DB
func (r *userRepository) Create(user *model.User) (*model.User, error) {
	if result := r.DB.Create(&user); result.Error != nil {
		// check unique constraint
		if isDuplicateKeyError(result.Error) {
			return nil, apperrors.NewBadRequest(apperrors.DuplicateEmail)
		}

		log.Printf("Could not create a user with email: %v. Reason: %v\n", user.Email, result.Error)
		return nil, apperrors.NewInternal()
	}

	return user, nil
}

// FindByEmail retrieves user row by email address
func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	user := &model.User{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("email", email)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

// Update updates the user in the DB
func (r *userRepository) Update(user *model.User) error {
	return r.DB.Save(&user).Error
}

// GetFriendAndGuildIds returns the id of the users friends and the guilds they are part of
func (r *userRepository) GetFriendAndGuildIds(userId string) (*[]string, error) {
	var ids []string
	result := r.DB.Raw(`
          SELECT g.id
          FROM guilds g
          JOIN members m on m.guild_id = g."id"
          where m.user_id = @userId
          UNION
          SELECT "User__friends"."id"
          FROM "users" "User" LEFT JOIN "friends" "User_User__friends" ON "User_User__friends"."user_id"="User"."id" LEFT
              JOIN "users" "User__friends" ON "User__friends"."id"="User_User__friends"."friend_id"
          WHERE ( "User"."id" = @userId )
	`, sql.Named("userId", userId)).Find(&ids)

	return &ids, result.Error
}

// GetRequestCount returns the amount of incoming friend requests the current user currently has
func (r *userRepository) GetRequestCount(userId string) (*int64, error) {
	var count int64
	err := r.DB.
		Table("users").
		Joins("JOIN friend_requests fr ON users.id = fr.sender_id").
		Where("fr.receiver_id = ?", userId).
		Count(&count).
		Error

	return &count, err
}

// isDuplicateKeyError checks if the provided error is a PostgreSQL duplicate key error
func isDuplicateKeyError(err error) bool {
	duplicate := regexp.MustCompile(`\(SQLSTATE 23505\)$`)
	return duplicate.MatchString(err.Error())
}

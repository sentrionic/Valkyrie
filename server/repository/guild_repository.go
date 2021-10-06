package repository

import (
	"errors"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log"
	"time"
)

// guildRepository is data/repository implementation
// of service layer GuildRepository
type guildRepository struct {
	DB *gorm.DB
}

// NewGuildRepository is a factory for initializing Guild Repositories
func NewGuildRepository(db *gorm.DB) model.GuildRepository {
	return &guildRepository{
		DB: db,
	}
}

// List returns all of the given users guilds
func (r *guildRepository) List(uid string) (*[]model.GuildResponse, error) {
	var guilds []model.GuildResponse
	result := r.DB.Raw(`
		SELECT distinct g."id",
		g."name",
		g."owner_id",
		g."icon",
		g."created_at",
		g."updated_at",
		((SELECT c."last_activity"
		 FROM channels c
		 JOIN guilds g ON g.id = c."guild_id"
		 WHERE g.id = member."guild_id"
		 order by c."last_activity" DESC
		 limit 1) > member."last_seen") AS "hasNotification",
		(SELECT c.id AS "default_channel_id"
		FROM channels c
	    JOIN guilds g ON g.id = c."guild_id"
		WHERE g.id = member."guild_id"
		ORDER BY c."created_at"
		LIMIT 1)
		FROM guilds g
		JOIN members as member
		on g."id"::text = member."guild_id"
		WHERE member."user_id" = ?
		ORDER BY g."created_at";
	`, uid).Find(&guilds)

	return &guilds, result.Error
}

// GuildMembers returns all members of the given guild and
// whether they are the given user IDs friend
func (r *guildRepository) GuildMembers(userId string, guildId string) (*[]model.MemberResponse, error) {
	var members []model.MemberResponse
	result := r.DB.Raw(`
		SELECT u.id,
		u.username,
		u.image,
		u."is_online",
		u."created_at",
		u."updated_at",
		m.nickname,
		m.color,
		EXISTS(
			SELECT 1
			FROM users
			LEFT JOIN friends f ON users.id = f."user_id"
			WHERE f."friend_id" = u.id
			AND f."user_id" = ?
		) AS is_friend
		FROM users AS u
		JOIN members m ON u."id"::text = m."user_id"
		WHERE m."guild_id" = ?
		ORDER BY (CASE WHEN m.nickname notnull THEN m.nickname ELSE u.username END)
	`, userId, guildId).Find(&members)

	return &members, result.Error
}

// Create inserts the given guild in the DB
func (r *guildRepository) Create(guild *model.Guild) (*model.Guild, error) {
	if result := r.DB.Create(&guild); result.Error != nil {
		log.Printf("Could not create a guild for user: %v. Reason: %v\n", guild.OwnerId, result.Error)
		return nil, apperrors.NewInternal()
	}

	return guild, nil
}

// FindUserByID returns a user containing all of their guilds
func (r *guildRepository) FindUserByID(uid string) (*model.User, error) {
	user := &model.User{}

	// we need to actually check errors as it could be something other than not found
	if err := r.DB.
		Preload("Guilds").
		Where("id = ?", uid).
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, apperrors.NewNotFound("uid", uid)
		}
		return user, apperrors.NewInternal()
	}

	return user, nil
}

// FindByID returns the guild for the given id containing all of their fields
func (r *guildRepository) FindByID(id string) (*model.Guild, error) {
	guild := &model.Guild{}

	if err := r.DB.
		Preload(clause.Associations).
		Where("id = ?", id).
		First(&guild).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return guild, apperrors.NewNotFound("id", id)
		}
		return guild, apperrors.NewInternal()
	}

	return guild, nil
}

// Save updates the given guild
func (r *guildRepository) Save(guild *model.Guild) error {
	if result := r.DB.Save(&guild); result.Error != nil {
		log.Printf("Could not update the guild with id: %v. Reason: %v\n", guild.ID, result.Error)
		return apperrors.NewInternal()
	}

	return nil
}

// RemoveMember removes the given user from the given guild
func (r *guildRepository) RemoveMember(userId string, guildId string) error {
	if result := r.DB.
		Exec("DELETE FROM members WHERE user_id = ? AND guild_id = ?", userId, guildId); result.Error != nil {
		log.Printf("Could not remove member with id: %s from the guild with id: %v. Reason: %v\n", userId, guildId, result.Error)
		return apperrors.NewInternal()
	}

	return nil
}

// Delete removes the given guild and all its associations
func (r *guildRepository) Delete(guildId string) error {
	if result := r.DB.
		Exec("DELETE FROM members WHERE guild_id = ?", guildId).
		Exec("DELETE FROM bans WHERE guild_id = ?", guildId).
		Exec("DELETE FROM guilds WHERE id = ?", guildId); result.Error != nil {
		log.Printf("Could not delete the guild with id: %v. Reason: %v\n", guildId, result.Error)
		return apperrors.NewInternal()
	}

	return nil
}

// UnbanMember removes the given user from the bans of the given guild
func (r *guildRepository) UnbanMember(userId string, guildId string) error {
	if result := r.DB.Exec("DELETE FROM bans WHERE guild_id = ? AND user_id = ?", guildId, userId); result.Error != nil {
		log.Printf("Could not unban the user with id: %v from the guild with id: %v. Reason: %v\n", userId, guildId, result.Error)
		return apperrors.NewInternal()
	}
	return nil
}

// GetBanList returns a list of all banned users from the given guild
func (r *guildRepository) GetBanList(guildId string) (*[]model.BanResponse, error) {
	var bans []model.BanResponse
	if result := r.DB.Raw(`
			select u.id, u.username, u.image
			from bans b
			join users u on b."user_id" = u.id
			where b."guild_id" = ?
		`, guildId).Scan(&bans); result.Error != nil {
		log.Printf("Could not get the ban list for the guild with id: %v. Reason: %v\n", guildId, result.Error)
		return &bans, apperrors.NewInternal()
	}

	return &bans, nil
}

// GetMemberSettings returns the given members settings in the given guild
func (r *guildRepository) GetMemberSettings(userId string, guildId string) (*model.MemberSettings, error) {
	settings := model.MemberSettings{}
	err := r.DB.
		Table("members").
		Where("user_id = ? AND guild_id = ?", userId, guildId).
		First(&settings)
	return &settings, err.Error
}

// UpdateMemberSettings updates the settings of the given member in the given guild
func (r *guildRepository) UpdateMemberSettings(settings *model.MemberSettings, userId string, guildId string) error {
	err := r.DB.
		Table("members").
		Where("user_id = ? AND guild_id = ?", userId, guildId).
		Updates(map[string]interface{}{
			"color":      settings.Color,
			"nickname":   settings.Nickname,
			"updated_at": time.Now(),
		}).
		Error
	return err
}

// FindUsersByIds returns the found users for the given user IDs and guild ID
func (r *guildRepository) FindUsersByIds(ids []string, guildId string) (*[]model.User, error) {
	var users []model.User
	result := r.DB.Raw(`
		SELECT u.*
		FROM users AS u
		JOIN members m ON u."id"::text = m."user_id"
		WHERE m."guild_id" = ?
		AND m."user_id" IN ?
	`, guildId, ids).Find(&users)

	return &users, result.Error
}

// GetMember returns user for the given userId and guildId
func (r *guildRepository) GetMember(userId, guildId string) (*model.User, error) {
	var user model.User
	result := r.DB.Raw(`
		SELECT u.*
		FROM users AS u
		JOIN members m ON u."id"::text = m."user_id"
		WHERE m."guild_id" = ?
		AND m."user_id" = ?
	`, guildId, userId).Find(&user)

	return &user, result.Error
}

// UpdateMemberLastSeen sets the LastSeen field of the given user to the current date
func (r *guildRepository) UpdateMemberLastSeen(userId, guildId string) error {
	err := r.DB.
		Table("members").
		Where("user_id = ? AND guild_id = ?", userId, guildId).
		Updates(map[string]interface{}{
			"last_seen": time.Now(),
		}).
		Error
	return err
}

// GetMemberIds returns the ids of all members of the given guild
func (r *guildRepository) GetMemberIds(guildId string) (*[]string, error) {
	var users []string
	result := r.DB.Raw(`
		SELECT u.id
		FROM users AS u
		JOIN members m ON u."id"::text = m."user_id"
		WHERE m."guild_id" = ?
	`, guildId).Find(&users)

	return &users, result.Error
}

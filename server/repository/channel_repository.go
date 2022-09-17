package repository

import (
	"database/sql"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"gorm.io/gorm"
	"log"
	"time"
)

// channelRepository is data/repository implementation
// of service layer ChannelRepository
type channelRepository struct {
	DB *gorm.DB
}

// NewChannelRepository is a factory for initializing Channel Repositories
func NewChannelRepository(db *gorm.DB) model.ChannelRepository {
	return &channelRepository{
		DB: db,
	}
}

// Create inserts a channel in the DB
func (r *channelRepository) Create(channel *model.Channel) (*model.Channel, error) {
	if result := r.DB.Create(&channel); result.Error != nil {
		log.Printf("Could not create a channel for guild: %v. Reason: %v\n", channel.GuildID, result.Error)
		return nil, apperrors.NewInternal()
	}

	return channel, nil
}

// GetGuildDefault fetches the oldest channel for the given guildId from the DB
func (r *channelRepository) GetGuildDefault(guildId string) (*model.Channel, error) {
	channel := model.Channel{}
	result := r.DB.
		Where("guild_id = ?", guildId).
		Order("created_at ASC").
		First(&channel)

	return &channel, result.Error
}

// Get fetches all public channels for the given guildId
// and the private channels the given user is part in
func (r *channelRepository) Get(userId string, guildId string) (*[]model.ChannelResponse, error) {
	var channels []model.ChannelResponse

	result := r.DB.
		Raw(`
			SELECT DISTINCT ON (c.id, c."created_at") c.id, c.name, 
			c."is_public", c."created_at", c."updated_at",
			(c."last_activity" > m."last_seen") AS "hasNotification"
			FROM channels AS c
			LEFT OUTER JOIN pcmembers as pc
			ON c."id"::text = pc."channel_id"::text
			LEFT OUTER JOIN members m on c."guild_id" = m."guild_id"
			WHERE c."guild_id"::text = ?
			AND (c."is_public" = true or pc."user_id"::text = ?)
			ORDER BY c."created_at"
		`, guildId, userId).
		Scan(&channels)

	return &channels, result.Error
}

// dmQuery represents the fetched fields for GetDirectMessages
type dmQuery struct {
	ChannelId string
	Id        string
	Username  string
	Image     string
	IsOnline  bool
	IsFriend  bool
}

// GetDirectMessages returns all DMs for the given user
func (r *channelRepository) GetDirectMessages(userId string) (*[]model.DirectMessage, error) {
	var results []dmQuery

	err := r.DB.
		Raw(`
			SELECT dm."channel_id", u.username, u.image, u.id, u."is_online", u."created_at", u."updated_at"
			FROM users u
			JOIN dm_members dm ON dm."user_id" = u.id
			WHERE u.id != @id
			AND dm."channel_id" IN (
				SELECT DISTINCT c.id
				FROM channels as c
				LEFT OUTER JOIN dm_members as dm
				ON c."id" = dm."channel_id"
				JOIN users u on dm."user_id" = u.id
				WHERE c."is_public" = false
				AND c.is_dm = true
				AND dm."is_open" = true
				AND dm."user_id" = @id
			)
			order by dm."updated_at" DESC 
		`, sql.Named("id", userId)).
		Scan(&results)

	var channels []model.DirectMessage

	// Turn into DirectMessage response
	for _, dm := range results {
		channel := model.DirectMessage{
			Id: dm.ChannelId,
			User: model.DMUser{
				Id:       dm.Id,
				Username: dm.Username,
				Image:    dm.Image,
				IsOnline: dm.IsOnline,
				IsFriend: dm.IsFriend,
			},
		}
		channels = append(channels, channel)
	}

	return &channels, err.Error
}

// GetDirectMessageChannel returns the dm channel ID of the given members
// if it exists.
func (r *channelRepository) GetDirectMessageChannel(userId string, memberId string) (*string, error) {
	var id string

	result := r.DB.
		Raw(`
			SELECT c.id
			FROM channels as c, dm_members dm 
			WHERE dm."channel_id" = c."id" AND c.is_dm = true AND c."is_public" = false
			GROUP BY c."id"
			HAVING array_agg(dm."user_id"::text) @> Array[?,?]
			AND count(dm."user_id") = 2;
		`, userId, memberId).
		Scan(&id)

	return &id, result.Error
}

// GetById returns the channel with its PCMembers for the given channel id
func (r *channelRepository) GetById(channelId string) (*model.Channel, error) {
	var channel model.Channel
	err := r.DB.Preload("PCMembers").Where("id = ?", channelId).First(&channel).Error
	return &channel, err
}

// GetPrivateChannelMembers returns the ids of all users
// that are members of the given channel
func (r *channelRepository) GetPrivateChannelMembers(channelId string) (*[]string, error) {
	var members []string
	err := r.DB.
		Raw("SELECT pc.user_id FROM pcmembers pc JOIN channels c ON pc.channel_id = c.id WHERE c.id = ?", channelId).
		Scan(&members).Error
	return &members, err
}

// AddDMChannelMembers inserts the given DM members in the DB
func (r *channelRepository) AddDMChannelMembers(members []model.DMMember) error {
	if err := r.DB.CreateInBatches(&members, len(members)).Error; err != nil {
		log.Printf("Could not add members to DM. Reason: %v\n", err)
		return apperrors.NewInternal()
	}
	return nil
}

// SetDirectMessageStatus opens or closes the dm channel for the given
// userId
func (r *channelRepository) SetDirectMessageStatus(dmId string, userId string, isOpen bool) error {
	err := r.DB.
		Table("dm_members").
		Where("channel_id = ? AND user_id = ?", dmId, userId).
		Updates(map[string]any{
			"is_open":    isOpen,
			"updated_at": time.Now(),
		}).
		Error
	return err
}

// OpenDMForAll opens the given dm channel for all users in the channel
func (r *channelRepository) OpenDMForAll(dmId string) error {
	err := r.DB.
		Table("dm_members").
		Where("channel_id = ? ", dmId).
		Updates(map[string]any{
			"is_open":    true,
			"updated_at": time.Now(),
		}).
		Error
	return err
}

// DeleteChannel deletes the given channel from the DB
func (r *channelRepository) DeleteChannel(channel *model.Channel) error {
	if result := r.DB.Delete(&channel); result.Error != nil {
		log.Printf("Could not delete the channel with id: %v. Reason: %v\n", channel, result.Error)
		return apperrors.NewInternal()
	}
	return nil
}

// UpdateChannel updates the given channel in the DB
func (r *channelRepository) UpdateChannel(channel *model.Channel) error {
	if result := r.DB.Save(&channel); result.Error != nil {
		log.Printf("Could not update the given channel: %v. Reason: %v\n", channel.ID, result.Error)
		return apperrors.NewInternal()
	}
	return nil
}

// CleanPCMembers removes all private channel members from the given channel
func (r *channelRepository) CleanPCMembers(channelId string) error {
	if result := r.DB.Exec("DELETE FROM pcmembers WHERE channel_id = ?", channelId); result.Error != nil {
		log.Printf("Could not clean members from the channel with id: %v. Reason: %v\n", channelId, result.Error)
		return apperrors.NewInternal()
	}
	return nil
}

// AddPrivateChannelMembers inserts the given member as PCMembers in the given channel
func (r *channelRepository) AddPrivateChannelMembers(memberIds []string, channelId string) error {
	var err error = nil
	for _, id := range memberIds {
		err = r.DB.Exec("INSERT INTO pcmembers VALUES (?, ?)", channelId, id).Error
	}

	if err != nil {
		log.Printf("Could not add members to private channel %s. Reason: %v\n", channelId, err)
		return apperrors.NewInternal()
	}
	return nil
}

// RemovePrivateChannelMembers removes the given ids from the PCMembers of the given channel
func (r *channelRepository) RemovePrivateChannelMembers(memberIds []string, channelId string) error {
	if err := r.DB.Exec("DELETE FROM pcmembers WHERE channel_id = ? AND user_id IN ?", channelId, memberIds).
		Error; err != nil {
		log.Printf("Could not remove members from private channel %s. Reason: %v\n", channelId, err)
		return apperrors.NewInternal()
	}

	return nil
}

// FindDMByUserAndChannelId returns the id of dm channel for the given channelId and userId
func (r *channelRepository) FindDMByUserAndChannelId(channelId, userId string) (string, error) {
	var id string
	err := r.DB.
		Raw("SELECT id FROM dm_members WHERE user_id = ? AND channel_id = ?", userId, channelId).
		Scan(&id).Error
	return id, err
}

// GetDMMemberIds returns the ids of all dm members for the given channel
func (r *channelRepository) GetDMMemberIds(channelId string) (*[]string, error) {
	var members []string
	err := r.DB.
		Raw("SELECT u.id FROM users u JOIN dm_members dm on u.id = dm.user_id WHERE dm.channel_id = ?", channelId).
		Scan(&members).Error
	return &members, err
}

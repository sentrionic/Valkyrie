package repository

import (
	"context"
	"encoding/json"
	"fmt"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"time"
)

type redisRepository struct {
	rds *redis.Client
}

// NewRedisRepository is a factory for initializing Redis Repositories
func NewRedisRepository(rds *redis.Client) model.RedisRepository {
	return &redisRepository{
		rds: rds,
	}
}

// Redis Prefixes
const (
	InviteLinkPrefix     = "inviteLink"
	ForgotPasswordPrefix = "forgot-password"
)

// SetResetToken inserts a password reset token in the DB and returns the generated token
func (r *redisRepository) SetResetToken(ctx context.Context, id string) (string, error) {
	uid, err := gonanoid.New()

	if err != nil {
		log.Printf("Failed to generate id: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if err = r.rds.Set(ctx, fmt.Sprintf("%s:%s", ForgotPasswordPrefix, uid), id, 24*time.Hour).Err(); err != nil {
		log.Printf("Failed to set link in redis: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	return uid, nil
}

// GetIdFromToken returns the user ID from the DB for the given token
func (r *redisRepository) GetIdFromToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", ForgotPasswordPrefix, token)
	val, err := r.rds.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", apperrors.NewBadRequest(apperrors.InvalidResetToken)
	}
	if err != nil {
		log.Printf("Failed to get value from redis: %v\n", err)
		return "", apperrors.NewInternal()
	}

	r.rds.Del(ctx, key)

	return val, nil
}

// SaveInvite inserts an invite for the given guild in the DB.
// If isPermanent is true, the invite won't expire
func (r *redisRepository) SaveInvite(ctx context.Context, guildId string, id string, isPermanent bool) error {

	invite := model.Invite{GuildId: guildId, IsPermanent: isPermanent}

	value, err := json.Marshal(invite)

	if err != nil {
		log.Printf("Error marshalling: %v\n", err.Error())
		return apperrors.NewInternal()
	}

	expiration := 24 * time.Hour
	if isPermanent {
		expiration = 0
	}

	if result := r.rds.Set(ctx, fmt.Sprintf("%s:%s", InviteLinkPrefix, id), value, expiration); result.Err() != nil {
		log.Printf("Failed to set invite link in redis: %v\n", err.Error())
		return apperrors.NewInternal()
	}

	return nil
}

// GetInvite returns the stored guild Id for the given token.
func (r *redisRepository) GetInvite(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", InviteLinkPrefix, token)
	val, err := r.rds.Get(ctx, key).Result()

	if err != nil {
		log.Printf("Failed to get invite link from redis: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	var invite model.Invite
	err = json.Unmarshal([]byte(val), &invite)

	if err != nil {
		log.Printf("Error unmarshalling: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if !invite.IsPermanent {
		r.rds.Del(ctx, key)
	}

	return invite.GuildId, nil
}

// InvalidateInvites deletes all permanent invites in the DB for the given guild
func (r *redisRepository) InvalidateInvites(ctx context.Context, guild *model.Guild) {
	for _, v := range guild.InviteLinks {
		key := fmt.Sprintf("%s:%s", InviteLinkPrefix, v)
		r.rds.Del(ctx, key)
	}
}

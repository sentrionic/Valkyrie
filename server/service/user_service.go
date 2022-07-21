package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
	"mime/multipart"
	"strings"
)

// UserService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository  model.UserRepository
	FileRepository  model.FileRepository
	RedisRepository model.RedisRepository
	MailRepository  model.MailRepository
}

// USConfig will hold repositories that will eventually be injected into
// this service layer
type USConfig struct {
	UserRepository  model.UserRepository
	FileRepository  model.FileRepository
	RedisRepository model.RedisRepository
	MailRepository  model.MailRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &userService{
		UserRepository:  c.UserRepository,
		FileRepository:  c.FileRepository,
		RedisRepository: c.RedisRepository,
		MailRepository:  c.MailRepository,
	}
}

// Get retrieves a user based on their uid
func (s *userService) Get(uid string) (*model.User, error) {
	return s.UserRepository.FindByID(uid)
}

// GetByEmail retrieves a user based on their email
func (s *userService) GetByEmail(email string) (*model.User, error) {

	// Sanitize email
	email = strings.ToLower(email)
	email = strings.TrimSpace(email)

	return s.UserRepository.FindByEmail(email)
}

// Register creates a user
func (s *userService) Register(user *model.User) (*model.User, error) {
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", user.Email)
		return nil, apperrors.NewInternal()
	}

	user.ID = GenerateId()
	user.Image = generateAvatar(user.Email)
	user.Password = hashedPassword

	return s.UserRepository.Create(user)
}

// Login reaches out to the UserRepository check if the user exists
// and then compares the supplied password with the provided password
// if a valid email/password combo is provided, u will hold all
// available user fields
func (s *userService) Login(email, password string) (*model.User, error) {
	user, err := s.UserRepository.FindByEmail(email)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	// verify
	match, err := comparePasswords(user.Password, password)

	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !match {
		return nil, apperrors.NewAuthorization(apperrors.InvalidCredentials)
	}

	return user, nil
}

func (s *userService) UpdateAccount(u *model.User) error {
	return s.UserRepository.Update(u)
}

func (s *userService) IsEmailAlreadyInUse(email string) bool {
	user, err := s.UserRepository.FindByEmail(email)

	if err != nil {
		return true
	}

	return user.ID != ""
}

func (s *userService) ChangeAvatar(header *multipart.FileHeader, directory string) (string, error) {
	return s.FileRepository.UploadAvatar(header, directory)
}

func (s *userService) DeleteImage(key string) error {
	return s.FileRepository.DeleteImage(key)
}

func (s *userService) ChangePassword(currentPassword, newPassword string, user *model.User) error {
	// verify
	match, err := comparePasswords(user.Password, currentPassword)

	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization(apperrors.InvalidOldPassword)
	}

	hashedPassword, err := hashPassword(newPassword)

	if err != nil {
		log.Printf("Unable to change password for email: %v\n", user.Email)
		return apperrors.NewInternal()
	}

	user.Password = hashedPassword

	return s.UserRepository.Update(user)
}

func (s *userService) ForgotPassword(ctx context.Context, user *model.User) error {
	token, err := s.RedisRepository.SetResetToken(ctx, user.ID)

	if err != nil {
		return err
	}

	return s.MailRepository.SendResetMail(user.Email, token)
}

func (s *userService) ResetPassword(ctx context.Context, password string, token string) (*model.User, error) {
	id, err := s.RedisRepository.GetIdFromToken(ctx, token)

	if err != nil {
		return nil, err
	}

	user, err := s.UserRepository.FindByID(id)

	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(password)

	if err != nil {
		log.Printf("Unable to reset password")
		return nil, apperrors.NewInternal()
	}

	user.Password = hashedPassword

	if err = s.UserRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetFriendAndGuildIds(userId string) (*[]string, error) {
	return s.UserRepository.GetFriendAndGuildIds(userId)
}

func (s *userService) GetRequestCount(userId string) (*int64, error) {
	return s.UserRepository.GetRequestCount(userId)
}

// generateAvatar returns a gravatar using the md5 hash of the email
func generateAvatar(email string) string {
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://gravatar.com/avatar/%s?d=identicon", hex.EncodeToString(hash[:]))
}

package service

import (
	"github.com/sentrionic/valkyrie/model"
)

// friendService acts as a struct for injecting an implementation of UserRepository
// and FriendRepository for use in service methods
type friendService struct {
	UserRepository   model.UserRepository
	FriendRepository model.FriendRepository
}

// FSConfig will hold repositories that will eventually be injected into
// this service layer
type FSConfig struct {
	UserRepository   model.UserRepository
	FriendRepository model.FriendRepository
}

// NewFriendService is a factory function for
// initializing a FriendService with its repository layer dependencies
func NewFriendService(c *FSConfig) model.FriendService {
	return &friendService{
		UserRepository:   c.UserRepository,
		FriendRepository: c.FriendRepository,
	}
}

func (f *friendService) GetFriends(id string) (*[]model.Friend, error) {
	return f.FriendRepository.FriendsList(id)
}

func (f *friendService) GetRequests(id string) (*[]model.FriendRequest, error) {
	return f.FriendRepository.RequestList(id)
}

func (f *friendService) GetMemberById(id string) (*model.User, error) {
	return f.FriendRepository.FindByID(id)
}

func (f *friendService) DeleteRequest(memberId string, userId string) error {
	return f.FriendRepository.DeleteRequest(memberId, userId)
}

func (f *friendService) RemoveFriend(memberId string, userId string) error {
	return f.FriendRepository.RemoveFriend(memberId, userId)
}

func (f *friendService) SaveRequests(user *model.User) error {
	return f.FriendRepository.Save(user)
}

package model

// Friend represents the api response of a user's friend.
type Friend struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	IsOnline bool   `json:"isOnline"`
} //@name Friend

// FriendService defines methods related to friend operations the handler layer expects
// any service it interacts with to implement
type FriendService interface {
	GetFriends(id string) (*[]Friend, error)
	GetRequests(id string) (*[]FriendRequest, error)
	GetMemberById(id string) (*User, error)
	DeleteRequest(memberId string, userId string) error
	RemoveFriend(memberId string, userId string) error
	SaveRequests(user *User) error
}

// FriendRepository defines methods related to friend db operations the service layer expects
// any repository it interacts with to implement
type FriendRepository interface {
	FindByID(id string) (*User, error)
	FriendsList(id string) (*[]Friend, error)
	RequestList(id string) (*[]FriendRequest, error)
	DeleteRequest(memberId string, userId string) error
	RemoveFriend(memberId string, userId string) error
	Save(user *User) error
}

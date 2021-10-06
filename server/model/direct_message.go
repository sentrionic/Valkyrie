package model

// DirectMessage is the json response of the channel ID
// and the other user of the DM.
type DirectMessage struct {
	Id   string `json:"id"`
	User DMUser `json:"user"`
} //@name DirectMessage

// DMUser is the other member of the DM.
type DMUser struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	Image    string `json:"image"`
	IsOnline bool   `json:"isOnline"`
	IsFriend bool   `json:"isFriend"`
} //@name DMUser

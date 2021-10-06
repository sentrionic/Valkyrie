package service

import (
	"encoding/json"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/ws"
	"log"
)

type socketService struct {
	Hub               ws.Hub
	GuildRepository   model.GuildRepository
	ChannelRepository model.ChannelRepository
}

// SSConfig will hold repositories that will eventually be injected into
// this service layer
type SSConfig struct {
	Hub               ws.Hub
	GuildRepository   model.GuildRepository
	ChannelRepository model.ChannelRepository
}

// NewSocketService is a factory function for
// initializing a SocketService with its repository layer dependencies
func NewSocketService(c *SSConfig) model.SocketService {
	return &socketService{
		Hub:               c.Hub,
		GuildRepository:   c.GuildRepository,
		ChannelRepository: c.ChannelRepository,
	}
}

func (s *socketService) EmitNewMessage(room string, message *model.MessageResponse) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.NewMessageAction,
		Data:   message,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitEditMessage(room string, message *model.MessageResponse) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.EditMessageAction,
		Data:   message,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitDeleteMessage(room, messageId string) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.DeleteMessageAction,
		Data:   messageId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitNewChannel(room string, channel *model.ChannelResponse) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.AddChannelAction,
		Data:   channel,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitNewPrivateChannel(members []string, channel *model.ChannelResponse) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.AddPrivateChannelAction,
		Data:   channel,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	for _, id := range members {
		s.Hub.BroadcastToRoom(data, id)
	}
}

func (s *socketService) EmitEditChannel(room string, channel *model.ChannelResponse) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.EditChannelAction,
		Data:   channel,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitDeleteChannel(channel *model.Channel) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.DeleteChannelAction,
		Data:   channel.ID,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, *channel.GuildID)
}

func (s *socketService) EmitEditGuild(guild *model.Guild) {

	response := guild.SerializeGuild("")

	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.EditGuildAction,
		Data:   response,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	members, err := s.GuildRepository.GetMemberIds(guild.ID)

	if err != nil {
		log.Printf("error getting member ids: %v\n", err)
	}

	for _, id := range *members {
		s.Hub.BroadcastToRoom(data, id)
	}
}

func (s *socketService) EmitDeleteGuild(guildId string, members []string) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.DeleteGuildAction,
		Data:   guildId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	for _, id := range members {
		s.Hub.BroadcastToRoom(data, id)
	}
}

func (s *socketService) EmitRemoveFromGuild(memberId, guildId string) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.RemoveFromGuildAction,
		Data:   guildId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, memberId)
}

func (s *socketService) EmitAddMember(room string, member *model.User) {

	response := model.MemberResponse{
		Id:        member.ID,
		Username:  member.Username,
		Image:     member.Image,
		IsOnline:  member.IsOnline,
		CreatedAt: member.CreatedAt,
		UpdatedAt: member.UpdatedAt,
		IsFriend:  false,
	}

	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.AddMemberAction,
		Data:   response,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitRemoveMember(room, memberId string) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.RemoveMemberAction,
		Data:   memberId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitNewDMNotification(channelId string, user *model.User) {

	response := model.DirectMessage{
		Id: channelId,
		User: model.DMUser{
			Id:       user.ID,
			Username: user.Username,
			Image:    user.Image,
			IsOnline: user.IsOnline,
			IsFriend: false,
		},
	}

	notification, err := json.Marshal(model.WebsocketMessage{
		Action: ws.NewDMNotificationAction,
		Data:   response,
	})

	if err != nil {
		log.Printf("error marshalling notification: %v\n", err)
	}

	pushToTop, err := json.Marshal(model.WebsocketMessage{
		Action: ws.PushToTopAction,
		Data:   channelId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	members, err := s.ChannelRepository.GetDMMemberIds(channelId)

	if err != nil {
		log.Printf("error getting member ids: %v\n", err)
	}

	for _, id := range *members {
		if id != user.ID {
			s.Hub.BroadcastToRoom(notification, id)
		}
		s.Hub.BroadcastToRoom(pushToTop, id)
	}
}

func (s *socketService) EmitNewNotification(guildId, channelId string) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.NewNotificationAction,
		Data:   guildId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	members, err := s.GuildRepository.GetMemberIds(guildId)

	if err != nil {
		log.Printf("error getting member ids: %v\n", err)
	}

	for _, id := range *members {
		s.Hub.BroadcastToRoom(data, id)
	}

	notification, err := json.Marshal(model.WebsocketMessage{
		Action: ws.NewNotificationAction,
		Data:   channelId,
	})

	if err != nil {
		log.Printf("error marshalling notification: %v\n", err)
	}

	s.Hub.BroadcastToRoom(notification, guildId)
}

func (s *socketService) EmitSendRequest(room string) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.SendRequestAction,
		Data:   "",
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
}

func (s *socketService) EmitAddFriendRequest(room string, request *model.FriendRequest) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.AddRequestAction,
		Data:   request,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, room)
	s.EmitSendRequest(room)
}

func (s *socketService) EmitAddFriend(user, member *model.User) {

	userResponse := model.Friend{
		Id:       user.ID,
		Username: user.Username,
		Image:    user.Image,
		IsOnline: user.IsOnline,
	}

	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.AddFriendAction,
		Data:   userResponse,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, member.ID)

	memberResponse := model.Friend{
		Id:       member.ID,
		Username: member.Username,
		Image:    member.Image,
		IsOnline: member.IsOnline,
	}

	data, err = json.Marshal(model.WebsocketMessage{
		Action: ws.AddFriendAction,
		Data:   memberResponse,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, user.ID)
}

func (s *socketService) EmitRemoveFriend(userId, memberId string) {
	data, err := json.Marshal(model.WebsocketMessage{
		Action: ws.RemoveFriendAction,
		Data:   memberId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, userId)

	data, err = json.Marshal(model.WebsocketMessage{
		Action: ws.RemoveFriendAction,
		Data:   userId,
	})

	if err != nil {
		log.Printf("error marshalling response: %v\n", err)
	}

	s.Hub.BroadcastToRoom(data, memberId)
}

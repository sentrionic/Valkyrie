package ws

import (
	"github.com/redis/go-redis/v9"
	"github.com/sentrionic/valkyrie/model"
)

// Hub contains all rooms and clients
type Hub struct {
	clients        map[*Client]bool
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	rooms          map[*Room]bool
	channelService model.ChannelService
	guildService   model.GuildService
	userService    model.UserService
	redisClient    *redis.Client
}

// Config will hold services that will eventually be injected into this
// service layer
type Config struct {
	UserService    model.UserService
	GuildService   model.GuildService
	ChannelService model.ChannelService
	Redis          *redis.Client
}

// NewWebsocketHub creates a new Hub
func NewWebsocketHub(c *Config) *Hub {
	return &Hub{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		rooms:          make(map[*Room]bool),
		channelService: c.ChannelService,
		guildService:   c.GuildService,
		userService:    c.UserService,
		redisClient:    c.Redis,
	}
}

// Run our websocket server, accepting various requests
func (hub *Hub) Run() {
	for {
		select {

		case client := <-hub.register:
			hub.registerClient(client)

		case client := <-hub.unregister:
			hub.unregisterClient(client)

		case message := <-hub.broadcast:
			hub.broadcastToClients(message)
		}
	}
}

func (hub *Hub) registerClient(client *Client) {
	hub.clients[client] = true
}

func (hub *Hub) unregisterClient(client *Client) {
	delete(hub.clients, client)
}

func (hub *Hub) broadcastToClients(message []byte) {
	for client := range hub.clients {
		client.send <- message
	}
}

// BroadcastToRoom sends the given message to all clients connected to the given room
func (hub *Hub) BroadcastToRoom(message []byte, roomId string) {
	if room := hub.findRoomById(roomId); room != nil {
		room.publishRoomMessage(message)
	}
}

func (hub *Hub) findRoomById(id string) *Room {
	var foundRoom *Room
	for room := range hub.rooms {
		if room.GetId() == id {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (hub *Hub) createRoom(id string) *Room {
	room := NewRoom(id, hub.redisClient)
	go room.RunRoom()
	hub.rooms[room] = true

	return room
}

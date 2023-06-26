package ws

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/sentrionic/valkyrie/model"
	"log"
)

// Room represents a websocket room
type Room struct {
	id         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *model.WebsocketMessage
	redis      *redis.Client
}

var ctx = context.Background()

// NewRoom creates a new Room
func NewRoom(id string, rds *redis.Client) *Room {
	return &Room{
		id:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *model.WebsocketMessage),
		redis:      rds,
	}
}

// RunRoom runs our room, accepting various requests
func (room *Room) RunRoom() {

	go room.subscribeToRoomMessages()

	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			room.publishRoomMessage(message.Encode())
		}
	}
}

// registerClientInRoom adds the client to the room
func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

// unregisterClientInRoom removes the client from the room
func (room *Room) unregisterClientInRoom(client *Client) {
	delete(room.clients, client)
}

// broadcastToClientsInRoom sends the given message to all members in the room
func (room *Room) broadcastToClientsInRoom(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}

// GetId returns the ID of the room
func (room *Room) GetId() string {
	return room.id
}

// publishRoomMessage publishes the message to all clients subscribing to the room
func (room *Room) publishRoomMessage(message []byte) {
	err := room.redis.Publish(ctx, room.GetId(), message).Err()

	if err != nil {
		log.Println(err)
	}
}

// subscribeToRoomMessages subscribes to messages in this room
func (room *Room) subscribeToRoomMessages() {
	pubsub := room.redis.Subscribe(ctx, room.GetId())

	ch := pubsub.Channel()

	for msg := range ch {
		room.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}

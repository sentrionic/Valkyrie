package ws

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sentrionic/valkyrie/model"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var newline = []byte{'\n'}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client represents the websockets client at the server
type Client struct {
	// The actual websockets connection.
	ID    string
	conn  *websocket.Conn
	hub   *Hub
	send  chan []byte
	rooms map[*Room]bool
}

func newClient(conn *websocket.Conn, hub *Hub, id string) *Client {
	return &Client{
		ID:    id,
		conn:  conn,
		hub:   hub,
		send:  make(chan []byte, 256),
		rooms: make(map[*Room]bool),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)

	_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))

	client.conn.SetPongHandler(func(string) error {
		_ = client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				_ = client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			_, _ = w.Write(message)

			// Attach queued chat messages to the current websockets message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				_, _ = w.Write(newline)
				_, _ = w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			_ = client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) disconnect() {
	client.hub.unregister <- client
	for room := range client.rooms {
		room.unregister <- client
	}
	close(client.send)
	_ = client.conn.Close()
}

// ServeWs handles websockets requests from clients requests.
func ServeWs(hub *Hub, ctx *gin.Context) {

	userId := ctx.MustGet("userId").(string)
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, hub, userId)

	go client.writePump()
	go client.readPump()

	hub.register <- client
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message model.ReceivedMessage
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	switch message.Action {
	// Join Room Actions
	case JoinChannelAction:
		client.handleJoinChannelMessage(message)
	case JoinGuildAction:
		client.handleJoinGuildMessage(message)
	case JoinUserAction:
		client.handleJoinRoomMessage(message)

	// Leave Room Actions
	case LeaveRoomAction:
		client.handleLeaveRoomMessage(message)
	case LeaveGuildAction:
		client.handleLeaveGuildMessage(message)

	// Chat Typing Actions
	case StartTypingAction:
		client.handleTypingEvent(message, AddToTypingAction)
	case StopTypingAction:
		client.handleTypingEvent(message, RemoveFromTypingAction)

	// Online Status Actions
	case ToggleOnlineAction:
		client.toggleOnlineStatus(true)
	case ToggleOfflineAction:
		client.toggleOnlineStatus(false)

	// Other
	case GetRequestCountAction:
		client.handleGetRequestCount()
	}
}

// handleJoinChannelMessage joins the given room if the user is a member in it
func (client *Client) handleJoinChannelMessage(message model.ReceivedMessage) {
	roomName := message.Room

	cs := client.hub.channelService
	channel, err := cs.Get(roomName)

	if err != nil {
		return
	}

	// Check if the user has access to the given channel
	if err = cs.IsChannelMember(channel, client.ID); err != nil {
		return
	}

	client.handleJoinRoomMessage(message)
}

// handleJoinGuildMessage joins the given guild if the user is member in it
func (client *Client) handleJoinGuildMessage(message model.ReceivedMessage) {
	roomName := message.Room

	gs := client.hub.guildService
	guild, err := gs.GetGuild(roomName)

	if err != nil {
		return
	}

	// Check if the user is member of the given guild
	if !isMember(guild, client.ID) {
		return
	}

	client.handleJoinRoomMessage(message)
}

// handleJoinRoomMessage joins the given room
func (client *Client) handleJoinRoomMessage(message model.ReceivedMessage) {
	roomName := message.Room

	room := client.hub.findRoomById(roomName)
	if room == nil {
		room = client.hub.createRoom(roomName)
	}

	client.rooms[room] = true

	room.register <- client
}

// handleLeaveGuildMessage leaves the room and updates the members last seen date
func (client *Client) handleLeaveGuildMessage(message model.ReceivedMessage) {
	_ = client.hub.guildService.UpdateMemberLastSeen(client.ID, message.Room)
	client.handleLeaveRoomMessage(message)
}

// handleLeaveRoomMessage leaves the room
func (client *Client) handleLeaveRoomMessage(message model.ReceivedMessage) {
	room := client.hub.findRoomById(message.Room)
	delete(client.rooms, room)

	if room != nil {
		room.unregister <- client
	}
}

// handleGetRequestCount returns the users incoming friend request count
func (client *Client) handleGetRequestCount() {
	if room := client.hub.findRoomById(client.ID); room != nil {
		count, err := client.hub.userService.GetRequestCount(client.ID)

		if err != nil {
			return
		}

		msg := model.WebsocketMessage{
			Action: RequestCountEmission,
			Data:   count,
		}
		room.broadcast <- &msg
	}
}

// handleTypingEvent emits the username of the currently typing user to the room
func (client *Client) handleTypingEvent(message model.ReceivedMessage, action string) {
	roomID := message.Room
	if room := client.hub.findRoomById(roomID); room != nil {
		msg := model.WebsocketMessage{
			Action: action,
			Data:   message.Message,
		}
		room.broadcast <- &msg
	}
}

// toggleOnlineStatus updates the users online status and emits it to all
// guilds the user is a member of and all of their friends
func (client *Client) toggleOnlineStatus(isOnline bool) {
	uid := client.ID
	us := client.hub.userService

	user, err := us.Get(uid)

	if err != nil {
		log.Printf("could not find user: %v", err)
		return
	}

	user.IsOnline = isOnline

	if err := us.UpdateAccount(user); err != nil {
		log.Printf("could not update user: %v", err)
		return
	}

	ids, err := us.GetFriendAndGuildIds(uid)

	if err != nil {
		log.Printf("could not find ids: %v", err)
		return
	}

	action := ToggleOfflineEmission
	if isOnline {
		action = ToggleOnlineEmission
	}

	for _, id := range *ids {
		if room := client.hub.findRoomById(id); room != nil {
			msg := model.WebsocketMessage{
				Action: action,
				Data:   uid,
			}
			room.broadcast <- &msg
		}
	}
}

// isMember checks if the user is member of the given guild
func isMember(guild *model.Guild, userId string) bool {
	for _, v := range guild.Members {
		if v.ID == userId {
			return true
		}
	}
	return false
}

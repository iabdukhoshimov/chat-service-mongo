package wsclient

import (
	"context"
	"encoding/json"
	"fmt"
	pb "golang-websocker/genproto/chat_service"
	"golang-websocker/pkg/helper"
	"golang-websocker/pkg/models"
	"log"
	"net/http"
	"reflect"
	"time"

	"golang-websocker/storage"

	"github.com/gorilla/websocket"
)

const (

	// time allowed to write message to peer
	writeWait = 10 * time.Second

	// time allowed to read message from peer
	pongWait = 60 * time.Second

	// ping period must be less from pongPeriod
	pingPeriod = (pongWait * 9) / 10

	// size allowed for message size
	maxMessageSize = 512 // kb
)

var upgrader = websocket.Upgrader{
	WriteBufferSize: 1024,
	ReadBufferSize:  1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Conn struct {

	// connection
	ws *websocket.Conn

	// send msg (in byte type)
	send chan []byte
}

// func read messages from the ws connects to the hub
func (s subscriptions) readPump(h *Hub) {
	c := s.conn
	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(appData string) error {
		err := c.ws.SetReadDeadline(time.Now().Add(pongWait))
		if err != nil {
			return err
		}

		return nil
	})

	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Println("---Error-> IsUnexpectedCloseError--->", err)
			}
			break
		}

		log.Printf("\n\nMessage received: from %s, in room %s\n", s.User.SenderName, s.room)
		dbm, err := sendMessage(context.TODO(), h.strg, s, msg)
		if err != nil {
			log.Println("Error: sendMessage --->", err.Error())
			break
		}

		mstrct := struct {
			Message      string `json:"message"`
			Message_Type string `json:"type"`
			Sender_Name  string `json:"sender_name"`
			Created_At   string `json:"created_at"`
			UserID       string `json:"user_id"`
		}{
			dbm.Message,
			dbm.Type,
			dbm.SenderName,
			dbm.CreatedAt,
			dbm.UserID,
		}

		respJSON, err := json.Marshal(&mstrct)
		if err != nil {
			log.Println("Error: json.Marshal --->", err.Error())
			break
		}

		m := message{respJSON, s.room}

		h.broadcast <- m
	}
}

func (c *Conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (s *subscriptions) writePump() {
	c := s.conn

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {

		select {
		case msg, ok := <-c.send:
			if !ok {
				if err := c.write(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Println("Error: ", err.Error())
					return
				}
			}

			log.Println("msg send ----> ")

			if err := c.write(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func ServeWS(ctx context.Context, w http.ResponseWriter, r *http.Request, roomId string, h *Hub) {

	var (
		user = models.CtxData{}
	)

	data := ctx.Value("data")
	if data == nil || reflect.TypeOf(data).String() != "models.CtxData" {
		log.Println("invalid ctx-data")
		fmt.Fprint(w, "invalid ctx-data")
		return
	}
	fmt.Println("ctx-data: ", data)

	user = data.(models.CtxData)

	if user.SenderName == "" || len(user.SenderName) > 60 {
		log.Println("invalid Sender-Name")
		fmt.Fprint(w, "invalid Sender-Name")
		return
	}

	var (
		userID     = user.UserID
		senderName = user.SenderName
		remoteAddr = r.RemoteAddr
	)

	if !helper.IsValidUUID(userID) {
		log.Println("invalid User-Id")
		fmt.Fprint(w, "invalid User-Id")
		return
	}

	if !helper.ValidName(senderName) {
		log.Println("invalid Sender-Name")
		fmt.Fprint(w, "invalid Sender-Name")
		return
	}

	u := User{
		UserID:     userID,
		SenderName: senderName,
		RemoteAddr: remoteAddr,
		EnterAt:    time.Now().Format(time.RFC3339),
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}

	c := &Conn{send: make(chan []byte, 256), ws: ws}

	s := subscriptions{c, roomId, u}

	h.register <- s
	go s.writePump()
	go s.readPump(h)
}

func sendMessage(ctx context.Context, strg storage.StorageI, s subscriptions, data []byte) (message *models.ChatMessage, err error) {

	um := pb.UserMessage{}
	err = json.Unmarshal(data, &um)
	if err != nil {
		return nil, err
	}
	um.UserId = s.User.UserID
	um.SenderName = s.User.SenderName
	// um.UserId =

	body := &pb.CreateMessageRequest{
		Message: &um,
		ChatId:  s.room,
	}
	message, err = strg.Chat().CreateMessage(
		ctx,
		body,
	)
	return message, err
}

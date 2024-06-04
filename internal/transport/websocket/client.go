package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"github.com/migmatore/study-platform-api/internal/core"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type ClientArgs struct {
	hub      *Hub
	conn     *websocket.Conn
	userId   int
	userRole core.RoleType
}

type ClientDeps struct {
	classroomUseCase ClassroomUseCase
}

type Client struct {
	userId   int
	userRole core.RoleType

	hub  *Hub
	conn *websocket.Conn

	send chan []byte

	classroomUseCase ClassroomUseCase
}

func NewClient(args ClientArgs, deps ClientDeps) *Client {
	return &Client{hub: args.hub, conn: args.conn, userId: args.userId, userRole: args.userRole, send: make(chan []byte, 256), classroomUseCase: deps.classroomUseCase}
}

func (c *Client) readPump() {
	defer func() {
		fmt.Println("close connection from readPump")
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		fmt.Println("read")

		if c.userRole != core.TeacherRole {
			continue
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		req := struct {
			Type        MessageType `json:"type"`
			ClassroomId int         `json:"classroom_id"`
			ElementId   *string     `json:"element_id,omitempty"`
		}{}

		if err := json.Unmarshal(message, &req); err != nil {
			log.Println("error while json unmarshalling")
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}

		students, err := c.classroomUseCase.Students(
			context.Background(),
			core.TokenMetadata{
				UserId: c.userId,
				Role:   string(c.userRole),
			},
			req.ClassroomId,
		)
		if err != nil {
			log.Println("error while getting classroom's students")
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			break
		}

		to := make([]Receiver, 0)

		for _, student := range students {
			to = append(to, Receiver{
				Id:   student.Id,
				role: core.StudentRole,
			})
		}

		//if req.Type == NewRoom {
		//	at := auth.NewAccessToken(os.GetEnv("LIVEKIT_API_KEY"), os.GetEnv("LIVEKIT_API_SECRET"))
		//	grant := &auth.VideoGrant{
		//		RoomJoin: true,
		//		Room:     "",
		//	}
		//	at.AddGrant(grant).
		//		SetIdentity("").
		//		SetValidFor(time.Hour)
		//
		//	c.send <- at.
		//}

		msg := NewMessage(message, to)

		c.hub.broadcast <- msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

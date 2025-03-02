package server

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/0x6flab/namegenerator"
	"golang.org/x/net/websocket"
)

type User struct {
	Conn     *websocket.Conn
	Username string
}

// Channel struct to store users and manage messages
type Channel struct {
	Name  string
	Users map[*User]bool
	mu    sync.Mutex
}

// Server struct to manage channels and connections
type Server struct {
	Channels map[string]*Channel
	Users    map[*websocket.Conn]*User
	mu       sync.Mutex
}

func NewServer() *Server {
	server := &Server{
		Channels: make(map[string]*Channel),
		Users:    make(map[*websocket.Conn]*User),
	}
	// creates a default channel
	server.CreateChannel("#general")
	return server
}

// handle incoming requests
func (s *Server) HandleConnection(ws *websocket.Conn) {
	defer ws.Close()

	ng := namegenerator.NewGenerator()
	username := ng.WithGender(namegenerator.Male).Generate()

	s.mu.Lock()
	user := &User{Conn: ws, Username: username}
	s.Users[ws] = user
	s.mu.Unlock()

	fmt.Println("User connected")

	s.JoinChannel(ws, user.Username, "#general")

	var msg = make([]byte, 1024)

	for {
		n, err := ws.Read(msg)
		if err != nil {
			if err == io.EOF {
				fmt.Println("BREAKS")
				break
			}
			fmt.Println("Error in your connection", err)
			continue
		}

		if n == 0 {
			continue
		}

		s.ProcessCommand(ws, msg[:n])
	}
}

// process input commands from client
func (s *Server) ProcessCommand(ws *websocket.Conn, msg []byte) {
	message := string(msg)

	parts := strings.SplitN(message, " ", 2)
	if len(parts) == 0 {
		return
	}

	command := strings.ToUpper(parts[0])
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	fmt.Println("command :: ", command)

	switch command {
	case "SHOW_CHANNELS":
		s.ShowChannels(ws)
	case "CREATE":
		s.CreateChannel(args)
	case "JOIN":
		s.JoinChannel(ws, s.Users[ws].Username, args)
	case "LEAVE":
		s.LeaveChannel(ws, args)
	case "MESSAGE":
		argsCommand := strings.SplitN(args, " ", 2)
		if len(argsCommand) < 2 {
			ws.Write([]byte("Usage is MESSAGE <channel_name> <message>"))
		}
		cn := argsCommand[0]
		data := argsCommand[1]
		from := s.Users[ws].Username
		s.Broadcast(from, cn, data)
	default:
		ws.Write([]byte("Unknown command\n"))
	}
}

func (s *Server) ShowChannels(ws *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	message := "Available Channels:\n"
	for name := range s.Channels {

		message += "- " + name + "\n"
	}
	fmt.Println("message :: ", message)
	ws.Write([]byte(message))

}

func (s *Server) Disconnect(ws *websocket.Conn) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Users[ws]
	if !exists {
		return
	}

	// Remove user from all channels
	for _, channel := range s.Channels {
		channel.mu.Lock()
		delete(channel.Users, user)
		channel.mu.Unlock()
	}

	delete(s.Users, ws)
	ws.Close()
	fmt.Println("Disconnected user:", user.Username)
}

func (s *Server) CreateChannel(name string) *Channel {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.Channels[name]; !exists {
		s.Channels[name] = &Channel{Name: name, Users: make(map[*User]bool)}
	}

	return s.Channels[name]
}

func (s *Server) GetChannel(name string) *Channel {
	s.mu.Lock()
	defer s.mu.Unlock()

	if channel, exists := s.Channels[name]; !exists {
		return channel
	}
	return nil
}

func (s *Server) JoinChannel(ws *websocket.Conn, username, channelName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var user *User
	var ok bool

	// get or create a user
	if user, ok = s.Users[ws]; !ok {
		user = &User{Conn: ws, Username: username}
	}

	// find/create channel to join
	channel, exists := s.Channels[channelName]

	if !exists {
		fmt.Println("channel does not exist, creating: ", channelName)
		channel = s.CreateChannel(channelName)
	}

	channel.mu.Lock()
	defer channel.mu.Unlock()

	// add user to channel
	channel.Users[user] = true
}

func (s *Server) LeaveChannel(ws *websocket.Conn, channelName string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user, exists := s.Users[ws]
	if !exists {
		fmt.Println("User is not connected")
		return
	}

	channel, exists := s.Channels[channelName]
	if !exists {
		fmt.Println("Channel does not exists")
		return
	}

	channel.mu.Lock()
	delete(channel.Users, user)
	channel.mu.Unlock()
	fmt.Println("Existed from channel")
}

func (s *Server) Broadcast(from string, channelName, msg string) {
	s.mu.Lock()
	channel, exists := s.Channels[channelName]
	s.mu.Unlock()

	if !exists {
		fmt.Println("Channel not found : ", channelName)
		return
	}

	channel.mu.Lock()
	msg = fmt.Sprintf("%s : %s", from, msg)
	fmt.Println("channel.Users :: ", channel.Users, channelName)
	channel.mu.Unlock()

	for user := range channel.Users {
		user.Conn.Write([]byte(msg))
	}

}

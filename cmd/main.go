package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gliderlabs/ssh"
	"github.com/teris-io/shortid"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/exp/rand"
	"golang.org/x/term"
)

type SSHHandler struct {
}
type UserMessage struct {
	message string
	time    string
	userTag string
}
type Room struct {
	users          []*User
	roomId         string
	roomName       string
	messageHistory []UserMessage
}
type User struct {
	session     ssh.Session
	nameTag     string
	currentRoom string
}

func roomExists(rooms []Room, selectedRoom Room) bool {
	for _, room := range rooms {
		if room.roomId == selectedRoom.roomId {
			return true
		}
	}
	return false
}

var clients sync.Map

// funtion to help populate a room for test purposes
func populateRoom(session ssh.Session, room *Room) {
	// Create random users
	for i := 0; i < rand.Intn(10)+1; i++ { // Generate 1 to 10 users
		nameTag := fmt.Sprintf("#%04d", rand.Intn(10000)) // Random 4-digit tag
		user := &User{session: session, nameTag: nameTag, currentRoom: room.roomId}
		room.users = append(room.users, user)
	}

	// Create random messages
	for i := 0; i < rand.Intn(20)+5; i++ { // Generate 5 to 20 messages
		user := room.users[rand.Intn(len(room.users))] // Pick a random user
		message := fmt.Sprintf("Message %d content", i+1)
		timestamp := time.Now().Add(time.Duration(-rand.Intn(3600)) * time.Second).Format("15:04:05") // Random timestamp in last hour
		room.messageHistory = append(room.messageHistory, UserMessage{
			message: message,
			time:    timestamp,
			userTag: user.nameTag,
		})
	}
}

func NewUser(session ssh.Session, userName string) *User {
	/*randomID, err := shortid.Generate()
	if err != nil {
		log.Fatal(err)
	}*/
	return &User{
		session:     session,
		nameTag:     userName,
		currentRoom: "000000",
	}
}
func (u *User) AddUserToMap() {
	clients.Store(u.currentRoom, u)
}
func (u *User) RemoveUserFromMap() {
	clients.Delete(u.currentRoom)
}

func (r *Room) AddUserToGroupChat(user *User, term *term.Terminal) {
	term.Write([]byte("\033[H\033[2J"))
	term.Write([]byte("\n"))
	for _, message := range r.messageHistory {
		term.Write([]byte(r.roomId))
		displayedMessage := fmt.Sprintf("%s at %s: %s\n ", message.userTag, message.time, message.message)
		term.Write([]byte(displayedMessage))
	}
}
func (r *Room) CreateRoom(user *User, rooms []Room, term *term.Terminal) { // --> CHANGE rooms []Room   --< pass as pointer cause its needed to add something to the slice

	for {
		term.Write([]byte("What is the name of the room you want to create?"))
		nameOfRoom, err := term.ReadLine()
		if err != nil {
			term.Write([]byte("Wrong command, try gain:\n"))
		}

		if roomExists(rooms, *r) {
			term.Write([]byte("Fail creating room, a room that name already exists\n"))
		} else {

			randomID, err := shortid.Generate()
			if err != nil {
				log.Fatal(err)
			}
			r.roomId = randomID
			r.roomName = nameOfRoom
			r.users = append(r.users, user)

			fmt.Printf("Id of the rrom just created  %v", r.roomId)

			return
		}

	}
}
func (u *User) JoinRoom(rooms []Room, term *term.Terminal) Room {
	for {
		term.Write([]byte("Whats the ID of the room you want to join?"))
		idOfRoom, err := term.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		for _, room := range rooms {
			if room.roomId == idOfRoom {
				u.currentRoom = room.roomId
				msg := fmt.Sprintf("Room found! Wecome to %s room (Room ID: %s)\n", room.roomName, room.roomId)
				term.Write([]byte(msg))
				return room
			}
		}
	}
}
func (h *SSHHandler) handleSSHSession(session ssh.Session) {
	rooms := []Room{}

	generalRoom := Room{
		users:    []*User{},
		roomId:   "000000",
		roomName: "General Room",
	}
	//populate a room to test
	populateRoom(session, &generalRoom)
	rooms = append(rooms, generalRoom)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	term := term.NewTerminal(session, "> ")
	term.Write([]byte("Welcome to secure chat!!!\n What's your Name Tag?\n"))
	nameTag, err := term.ReadLine()
	term.Write([]byte("What do you want to join?\n- Chat Room (cmd: CR)\n- Create a One On One Room (cmd: CCOOO)\n- Join a One On One Room (cmd: JCOOO)\n"))
	userChoice, err := term.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	for {
		switch userChoice {
		case "CR":
			term.Write([]byte("Joined a chat room"))
			return
		case "CCOOO": // creat one on one room
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			user := NewUser(session, nameTag) //create new user
			newRoom := Room{}
			newRoom.CreateRoom(user, rooms, term) // ------->> CHANGE --> rooms must be pass as a pointer
			fmt.Printf(newRoom.roomId)
			newRoom.AddUserToGroupChat(user, term)

			userChoice, err = term.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			user.AddUserToMap() // add to a connection
			fmt.Println("SSH connection established successfully!")
			fmt.Println("Print user info: ", user)
		case "JCOOO": // join a room
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			user := NewUser(session, nameTag) //create new user
			user.AddUserToMap()               // add to a connection

			room := user.JoinRoom(rooms, term)
			room.AddUserToGroupChat(user, term)
			fmt.Println("SSH connection established successfully!")
			fmt.Println("Print user info: ", user)
		default:
			term.Write([]byte("try again"))

		}
	}
}
func NewSSHHandler() *SSHHandler {
	return &SSHHandler{}
}
func main() {
	sshPort := ":2222"
	handler := NewSSHHandler()
	server := &ssh.Server{ //defining ssh server
		Addr:    sshPort,
		Handler: handler.handleSSHSession,
		PublicKeyHandler: (func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		ServerConfigCallback: func(ctx ssh.Context) *gossh.ServerConfig {
			cfg := &gossh.ServerConfig{
				ServerVersion: "SSH-2.0-OpenSSH_8.9p1",
			}
			cfg.Ciphers = []string{"chacha20-poly1305@openssh.com"}
			return cfg
		},
	}
	b, err := os.ReadFile("../keys/private_key")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("here : ", b)
	privateKey, err := gossh.ParsePrivateKey(b)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}
	//fmt.Println("Key in the main : ", privateKey)
	server.AddHostKey(privateKey)
	log.Printf("Starting SSH server on port %s", sshPort)
	log.Fatal(server.ListenAndServe())
}

package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/teris-io/shortid"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHHandler struct {
}
type chatMessages struct {
	message string
	time    string
	userTag string
}
type User struct {
	session     ssh.Session
	nameTag     string
	currentRoom string
}

var clients sync.Map

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

type Room struct {
	users          []User
	roomId         string
	roomName       string
	messageHistory []string
}

func (r *Room) CreateRoom(term *term.Terminal) {
	randomID, err := shortid.Generate()

	if err != nil {
		log.Fatal(err)
	}
	r.roomId = randomID
	term.Write([]byte("Whats the ID of the room you want to join?"))

}
func (u *User) JoinRoom(rooms []Room, idOfRoom string, term *term.Terminal) {
	for _, room := range rooms {
		if room.roomId == idOfRoom {
			u.currentRoom = room.roomId
			msg := fmt.Sprintf("Room found! Wecome to %s room (Room ID: %s)\n", room.roomName, room.roomId)
			term.Write([]byte(msg))
		}
	}
}
func (h *SSHHandler) handleSSHSession(session ssh.Session) {
	rooms := []Room{}

	newRoom := Room{
		users:    []User{},
		roomId:   "000000",
		roomName: "General Room",
	}
	rooms = append(rooms, newRoom)
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
			newRoom.CreateRoom(term)
			user.AddUserToMap() // add to a connection
			fmt.Println("SSH connection established successfully!")
			fmt.Println("Print user info: ", user)
		case "JCOOO": // join a room
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			user := NewUser(session, nameTag) //create new user
			user.AddUserToMap()               // add to a connection
			term.Write([]byte("Whats the ID of the room you want to join?"))
			idOfRoom, err := term.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			user.JoinRoom(rooms, idOfRoom, term)
			fmt.Println("SSH connection established successfully!")
			fmt.Println("Print user info: ", user)

		default:
			term.Write([]byte("Wrong command, try gain:\n"))
			userChoice, err = term.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
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
	b, err := os.ReadFile("./keys/private_key")
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

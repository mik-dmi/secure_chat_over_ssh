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
type User struct {
	session      ssh.Session
	nameTag      string
	connectionID string
}

var clients sync.Map

func NewUser(session ssh.Session, userName string) *User {
	randomID, err := shortid.Generate()
	if err != nil {
		log.Fatal(err)
	}
	return &User{
		session:      session,
		nameTag:      userName,
		connectionID: randomID,
	}
}
func (u *User) AddUserToMap() {
	clients.Store(u.connectionID, u)
}
func (u *User) RemoveUserFromMap() {
	clients.Delete(u.connectionID)
}

func (h *SSHHandler) handleSSHSession(session ssh.Session) {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	term := term.NewTerminal(session, "> ")
	term.Write([]byte("Welcome to secure chat!!!\n What's your Name Tag?\n"))
	nameTag, err := term.ReadLine()
	term.Write([]byte("What do you want to join?\n- Chat Room (cmd: CR)\n- Chat One On One (cmd: COOO)\n"))
	userChoice, err := term.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	for {
		switch userChoice {
		case "CR":
			term.Write([]byte("Joined a chat room"))

			return
		case "COOO":
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			user := NewUser(session, nameTag)

			fmt.Println("SSH connection established successfully!")
			fmt.Println("Print user info: ", user)
			return
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
	server := &ssh.Server{
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

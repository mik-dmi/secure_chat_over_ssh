package main

import (
	"log"
	"os"
	"sync"

	"secure_chat_over_ssh/chat"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

func main() {
	sshPort := ":2222"
	room := &chat.Room{
		RoomName:       "General Room", // is waiting room
		Users:          sync.Map{},
		MessageHistory: []*chat.UserMessage{},
	}
	waitingRoom := &chat.Room{
		RoomName:       "Waiting room", // is waiting room
		Users:          sync.Map{},
		MessageHistory: []*chat.UserMessage{},
	}

	handler := chat.NewSSHHandler(room)           // maybe delete later
	handler.RoomManager.Rooms.Store("0001", room) // rm id ---> 0001

	handler.RoomManager.Rooms.Store("0000", waitingRoom) // rm id ---> 0001

	server := &ssh.Server{
		Addr:    sshPort,
		Handler: handler.HandleSSHSession,

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
	b, err := os.ReadFile("../keys/dockerkey")
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
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start SSH server: %v\n", err)
	}
}

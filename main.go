package main

import (
	"io"
	"log"
	"os"

	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
	"golang.org/x/term"
)

type SSHHandler struct {
}

func (h *SSHHandler) handleSSHSession(session ssh.Session) {
	io.WriteString(session, "Remote forwarding available...\n")

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	term := term.NewTerminal(session, "> ")
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
			term.Write([]byte("Individual chat room created"))
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
	b, err := os.ReadFile("./keys/private_key")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("here : ", b)
	privateKey, err := gossh.ParsePrivateKey(b)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	handler := NewSSHHandler()
	server := &ssh.Server{
		Addr:        sshPort,
		Handler:     handler.handleSSHSession,
		HostSigners: []ssh.Signer{privateKey},
		PublicKeyHandler: (func(ctx ssh.Context, key ssh.PublicKey) bool {
			return true
		}),
		Version: "SSH-2.0-OpenSSH_8.9p1",
	}

	server.AddHostKey(privateKey)
	log.Printf("Starting SSH server on port %s", sshPort)
	log.Fatal(server.ListenAndServe())
}

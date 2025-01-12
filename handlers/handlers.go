package handlers

import (
	"fmt"

	"secure_chat_over_ssh/chat"
	"secure_chat_over_ssh/utils"
	"sync"

	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SSHHandler struct {
	RoomManager  *chat.RoomManager
	UsersManager *chat.UsersManager
}

func NewSSHHandler() *SSHHandler {
	return &SSHHandler{
		RoomManager:  chat.NewRoomManager(),
		UsersManager: chat.NewUsersManager(),
	}
}

var AllUsersMap sync.Map

func (h *SSHHandler) HandleSSHSession(session ssh.Session) {

	room := &chat.Room{
		RoomName: "General Room",
	}
	h.RoomManager.Rooms.Store("00000", room)

	//populate a room to test
	//utils.PopulateRoom(session, h.RoomManager)

	/*oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)
	term := term.NewTerminal(session, "> ")
	*/
	term := terminal.NewTerminal(session, "> ")

	term.Write([]byte("Welcome to secure chat!!!\n What's your User Tag?\n"))

	user, err := utils.NewUser(session, h.UsersManager, term)

	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Print in  main : ", user.UserTag)

	for {
		term.Write([]byte("What do you want to join?\n- Chat Room (cmd: CR)\n- Create a One On One Room (cmd: CCOOO)\n- Join a One On One Room (cmd: JCOOO)\n"))
		userChoice, err := term.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch userChoice {
		case "CR":
			term.Write([]byte("Joined a chat room"))
			return
		case "CCOOO": // creat one on one room
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			//create new user
			newRoom := h.RoomManager.CreateRoom(user, term) // ------->> CHANGE --> rooms must be pass as a pointer
			user.CurrentRoomName = newRoom.RoomName
			h.RoomManager.GetIntoAGroupChat(term, newRoom)
			newRoom.WriteMessageToChat(term, user)

		case "JCOOO": // join a room
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			room = h.RoomManager.JoinRoom(user, term)
			h.RoomManager.GetIntoAGroupChat(term, room)
			room.WriteMessageToChat(term, user)

		default:
			term.Write([]byte("try again"))
		}
	}
}

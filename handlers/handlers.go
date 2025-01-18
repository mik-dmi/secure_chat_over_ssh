package handlers

import (
	"fmt"

	"secure_chat_over_ssh/chat"
	"secure_chat_over_ssh/utils"
	"sync"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
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
	term := term.NewTerminal(session, "> ")

	term.Write([]byte("Welcome to secure chat!!!\n What's your User Tag?\n"))

	user, err := chat.NewUser(session, h.UsersManager, term)

	if err != nil {
		fmt.Println(err)
		return
	}
	//fmt.Println("Print in  main : ", user.UserTag)

	for {
		term.Write([]byte("Choose an option by typing the corresponding command (cmd):\n- Join General Chat Room (cmd: JGR)\n- Create A Chat Room (cmd: CR)\n- Join a Chat Room (cmd: JR)\n"))
		userChoice, err := term.ReadLine()
		if err != nil {
			fmt.Println(err)
			return
		}
		switch userChoice {
		case "JGR":
			term.Write([]byte("Joined a chat room"))
			return
		case "CR": // creat one on one room
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			//create new user
			newRoom := h.RoomManager.CreateRoom(user, term) // ------->> CHANGE --> rooms must be pass as a pointer
			user.CurrentRoomName = newRoom.RoomName
			h.RoomManager.GetIntoAGroupChat(term, newRoom)
			newRoom.WriteMessageToChat(term, user)
			continue
		case "JR": // join a room
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			room = h.RoomManager.JoinRoom(user, term)
			h.RoomManager.GetIntoAGroupChat(term, room)
			room.WriteMessageToChat(term, user)
			continue
		default:
			utils.ClearUserTerminal(term)
			term.Write([]byte("-> Command was invalid, try again!\n\n"))
		}
	}
}

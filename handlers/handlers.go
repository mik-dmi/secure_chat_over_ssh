package handlers

import (
	"fmt"

	"secure_chat_over_ssh/chat"
	"secure_chat_over_ssh/utils"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

type SSHHandler struct {
	RoomManager  *chat.RoomManager
	UsersManager *chat.UsersManager
	Room         *chat.Room
}

func NewSSHHandler(room *chat.Room) *SSHHandler {
	return &SSHHandler{
		RoomManager:  chat.NewRoomManager(),
		UsersManager: chat.NewUsersManager(),
		Room:         room,
	}
}

func (h *SSHHandler) HandleSSHSession(session ssh.Session) {
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

	user, err := h.UsersManager.NewUser(session, term)

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

			h.Room.Users.Store(user.UserTag, user)
			h.RoomManager.GetIntoAGroupChat(term, h.Room)
			h.Room.Users.Range(func(key, value any) bool {
				fmt.Printf("key: %v, value: %v\n", key, value)
				return true
			})

			h.Room.WriteMessageToChat(term, user)
			continue
		case "CR":
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			newRoom := h.RoomManager.CreateRoom(user, term) // ------->> CHANGE --> rooms must be pass as a pointer
			if newRoom == nil {                             // the user what to exit creating room option
				continue
			}
			user.CurrentRoomName = newRoom.RoomName
			h.RoomManager.GetIntoAGroupChat(term, newRoom)
			err = newRoom.WriteMessageToChat(term, user)
			if err != nil {
				return //problem with the terminal
			}
			continue
		case "JR":
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			room := h.RoomManager.JoinRoom(user, term)
			if room == nil { // the user what to exit join room option
				continue
			}
			h.RoomManager.GetIntoAGroupChat(term, room)
			room.WriteMessageToChat(term, user)
			continue
		default:
			utils.ClearUserTerminal(term)
			term.Write([]byte("-> Command was invalid, try again!\n\n"))
		}
	}
}

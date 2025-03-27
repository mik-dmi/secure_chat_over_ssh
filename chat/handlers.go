package chat

import (
	"fmt"

	"secure_chat_over_ssh/utils"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

type SSHHandler struct {
	RoomManager  *RoomManager
	UsersManager *UsersManager
	Room         *Room // just for the general  room ??
}

func NewSSHHandler(room *Room) *SSHHandler {
	return &SSHHandler{
		RoomManager:  NewRoomManager(),
		UsersManager: NewUsersManager(),
		Room:         room,
	}
}

var banner = ` 
 ____ ____  _   _   ____  _____ ____ _   _ ____  _____ 
/ ___/ ___|| | | | / ___|| ____/ ___| | | |  _ \| ____|
\___ \___ \| |_| | \___ \|  _|| |   | | | | |_) |  _|  
 ___) |__) |  _  |  ___) | |__| |___| |_| |  _ <| |___ 
|____/____/|_| |_| |____/|_____\____|\___/|_| \_\_____|
 / ___| | | |  / \|_   _|                              
| |   | |_| | / _ \ | |                                
| |___|  _  |/ ___ \| |                                
 \____|_| |_/_/   \_\_|                                
`

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

	introductionText := fmt.Sprintf("Welcome to secure chat!!!%s\nWhat's your User Tag?\n", banner)
	term.Write([]byte(introductionText))

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

			h.RoomManager.WriteMessageToChat(user, "0001", h) //0001 --> is the roomManagerMapID
		case "CR":
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			newRoom, roomManagerMapID := h.RoomManager.CreateRoom(user) // ------->> CHANGE --> rooms must be pass as a pointer
			if newRoom == nil {                                         // the user what to exit creating room option
				continue
			}
			user.CurrentRoomName = newRoom.RoomName
			h.RoomManager.GetIntoAGroupChat(user.Term, newRoom)
			err = h.RoomManager.WriteMessageToChat(user, roomManagerMapID, h)
			if err != nil {
				return //problem with the terminal
			}
			continue
		case "JR":
			//fmt.Println("Local Addr : ", session.LocalAddr().String())
			roomManagerMapID, room := h.RoomManager.JoinRoom(user, term)
			if room == nil { // the user what to exit join room option
				continue
			}
			h.RoomManager.GetIntoAGroupChat(term, room)
			h.RoomManager.WriteMessageToChat(user, roomManagerMapID, h)
			continue
		default:
			utils.ClearUserTerminal(term)
			term.Write([]byte("-> Command was invalid, try again!\n\n"))
		}
	}
}

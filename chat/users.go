package chat

import (
	"fmt"
	"secure_chat_over_ssh/utils"

	"sync"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

type UserMessage struct {
	Message string
	Time    string
	UserTag string
}
type User struct {
	Session         ssh.Session
	UserTag         string //unique id
	CurrentRoomName string
	Term            *term.Terminal
}
type UsersManager struct {
	Users sync.Map // key = string (user tag), value = *User
}

func NewUsersManager() *UsersManager {
	return &UsersManager{}
}
func (rm *RoomManager) GetIntoAGroupChat(term *term.Terminal, room *Room) {
	utils.ClearUserTerminal(term)
	roomManagerMapID, ok := rm.GetRoomIDByRoomObject(room)
	if !ok {
		term.Write([]byte(fmt.Sprintf("Error getting room ID %v from room map\n", roomManagerMapID)))
		return
	}

	term.Write([]byte(fmt.Sprintf("You just joinned: %s (Room ID: %s)", room.RoomName, roomManagerMapID)))
	term.Write([]byte("\n"))
	term.Write([]byte(fmt.Sprintf("%s \n", room.RoomName)))

	room.messagesMu.Lock()
	if len(room.MessageHistory) == 0 {
		room.messagesMu.Unlock()
		return
	}

	messages := append([]*UserMessage{}, room.MessageHistory...)
	room.messagesMu.Unlock()

	for _, message := range messages {

		displayedMessage := fmt.Sprintf("%s at %s: %s\n", message.UserTag, message.Time, message.Message)
		term.Write([]byte(displayedMessage))
	}

}

func (user *UsersManager) NewUser(session ssh.Session, term *term.Terminal) (*User, error) {
	var userTag string
	for {
		term.Write([]byte("Enter a unique user tag:\n"))
		line, err := term.ReadLine()
		if err != nil {
			err = fmt.Errorf("error reading user tag: %s", err)
			return nil, err
		}
		userTag = line
		//log.Printf(" 1 - userTag is: %s", userTag)

		if len(userTag) == 0 {
			term.Write([]byte("User tag cannot be empty. Try again:\n"))
			continue
		}
		if len(userTag) > 40 {
			term.Write([]byte("User tag must be 20 characters or less. Try again:\n"))
			continue
		}

		// Check if userTag already exists
		_, ok := user.Users.Load(userTag)
		if ok {
			term.Write([]byte(fmt.Sprintf("The user tag '%s' is already in use. Try again:\n", userTag)))
		} else {
			break
		}
	}
	utils.ClearUserTerminal(term)

	//log.Printf(" 2 - userTag is: %s", userTag)

	newUser := &User{
		Session:         session,
		UserTag:         userTag,
		CurrentRoomName: "Waiting room",
		Term:            term,
	}
	user.Users.Store(userTag, newUser)

	return newUser, nil

}

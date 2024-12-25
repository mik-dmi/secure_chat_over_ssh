package chat

import (
	"fmt"
	"log"
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
	UserTag         string
	CurrentRoomName string
	Term            *term.Terminal
}

type UsersManager struct {
	Users sync.Map // key = string (Room ID), value = *Room
}

func NewUsersManager() *UsersManager {
	return &UsersManager{}
}

func (rm *RoomManager) GetIntoAGroupChat(term *term.Terminal, room *Room) {
	term.Write([]byte("\033[H\033[2J"))
	roomID, ok := rm.GetRoomIDByRoomObject(room)
	if !ok {
		term.Write([]byte(fmt.Sprintf("Error getting room ID %v from room map ", roomID)))
		return
	}
	term.Write([]byte(fmt.Sprintf("You just joinned: %s (Room ID: %s)", room.RoomName, roomID)))
	term.Write([]byte("\n"))

	for _, message := range room.MessageHistory {
		term.Write([]byte(room.RoomName))
		displayedMessage := fmt.Sprintf("%s at %s: %s\n ", message.UserTag, message.Time, message.Message)
		term.Write([]byte(displayedMessage))
	}

}

func WriteMessageToChat(term *term.Terminal, room *Room, userTag string) {
	message, err := term.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	if message == "exit" {
		term.Write([]byte("Left the room"))
		// STILL NEED to take out the user of the room map and so on
	}
	for {
		userMessage, err := term.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		room.UpdateRoomChat(userMessage, userTag)
	}
}

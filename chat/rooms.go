package chat

import (
	"fmt"
	"log"
	"secure_chat_over_ssh/utils"

	"sync"
	"time"

	"github.com/teris-io/shortid"
	"golang.org/x/term"
)

type Room struct {
	RoomName       string
	Users          sync.Map //  key = string ( UserTag), value = *User
	MessageHistory []*UserMessage
	messagesMu     sync.Mutex
}

type RoomManager struct {
	Rooms sync.Map // key = string (Room ID == the id is the name of the room), value = *Room
}

func NewRoomManager() *RoomManager {
	return &RoomManager{}
}

func (rm *RoomManager) GetRoomIDByRoomObject(room *Room) (string, bool) {
	var roomID string
	found := false
	rm.Rooms.Range(func(key, value interface{}) bool {
		if value == room {
			roomID = key.(string)
			found = true
			return false
		}
		return true
	})

	return roomID, found
}

func (rm *RoomManager) GetRoomByID(id string) (*Room, bool) {
	value, ok := rm.Rooms.Load(id)
	if !ok {
		return nil, false
	}
	return value.(*Room), true
}
func (rm *RoomManager) GetRoomByName(name string) (*Room, bool) {
	var foundRoom *Room
	var found bool

	rm.Rooms.Range(func(key, value any) bool {
		room, ok := value.(*Room)
		if !ok {
			return true
		}
		if room.RoomName == name {
			foundRoom = room
			found = true
			return false
		}
		return true
	})

	return foundRoom, found
}

func (rm *RoomManager) DeleteRoom(id string) {
	rm.Rooms.Delete(id)
}

func (rm *RoomManager) ListRooms() []*Room {
	rooms := make([]*Room, 0)
	rm.Rooms.Range(func(key, value interface{}) bool {
		if room, ok := value.(*Room); ok {
			rooms = append(rooms, room)
		}
		return true
	})
	return rooms
}

func (rm *RoomManager) JoinRoom(user *User, term *term.Terminal) *Room {
	for {
		term.Write([]byte("What's the ID of the room you want to join?\n"))
		idOfRoom, err := term.ReadLine()
		if err != nil {
			log.Printf("Error reading input: %v", err)
			return nil
		}
		if idOfRoom == "exit" {
			utils.ClearUserTerminal(term)
			return nil
		}
		room, ok := rm.GetRoomByID(idOfRoom)
		if !ok {
			term.Write([]byte(fmt.Sprintf("The room ID %v does not exist\n", idOfRoom)))
			continue
		}
		room.Users.Store(user.UserTag, user)
		user.CurrentRoomName = room.RoomName
		msg := fmt.Sprintf("Room found! Welcome to %s room (Room ID: %s)\n", room.RoomName, idOfRoom)
		term.Write([]byte(msg))
		return room
	}
}

func (rm *RoomManager) CreateRoom(user *User, term *term.Terminal) *Room { // --> CHANGE rooms []Room   --< pass as pointer cause its needed to add something to the slice
	for {
		term.Write([]byte("What is the name of the room you want to create?"))
		nameOfRoom, err := term.ReadLine()
		if err != nil {
			term.Write([]byte("Error reading terminal\n"))
			return nil
		}
		//if user wants to go back to the initial menu
		if nameOfRoom == "exit" {
			utils.ClearUserTerminal(term)
			return nil
		}

		_, ok := rm.GetRoomByName(nameOfRoom)

		if ok {
			term.Write([]byte("Fail creating room, a room that name already exists\n"))
		} else {

			RoomID, err := shortid.Generate()
			if err != nil {
				fmt.Println("Failed to generate a room id")
				term.Write([]byte("Fail to generate a room id, try again\n"))
				continue
			}

			room := &Room{
				RoomName:       nameOfRoom,
				Users:          sync.Map{},
				MessageHistory: []*UserMessage{},
			}

			room.Users.Store(user.UserTag, user)

			rm.Rooms.Store(RoomID, room)
			return room
		}
	}
}

func (r *Room) ShowAllUserInRoom(userTerminal *term.Terminal) {
	userTerminal.Write([]byte("List of all the Users in the room:\n"))
	r.Users.Range(func(_, value interface{}) bool {
		user := value.(*User)
		userTerminal.Write([]byte(fmt.Sprintf("- %s\n", user.UserTag)))
		return true
	})
}

func (r *Room) UpdateRoomChat(userMessage string, userTag string) {

	currentTime := time.Now()
	formattedMessageTime := currentTime.Format("15:04")
	r.Users.Range(func(_, value interface{}) bool {
		user := value.(*User) // Type assert the value to *User

		if user.UserTag != userTag {
			user.Term.Write([]byte(fmt.Sprintf("%s at %s: %s\n", userTag, formattedMessageTime, userMessage)))
		}
		return true
	})
	r.messagesMu.Lock()
	defer r.messagesMu.Unlock()
	r.MessageHistory = append(r.MessageHistory, &UserMessage{
		Message: userMessage,
		Time:    formattedMessageTime,
		UserTag: userTag,
	})
}

func (room *Room) WriteMessageToChat(term *term.Terminal, user *User) error {

	for {
		userMessage, err := term.ReadLine()
		if err != nil {
			fmt.Println("Failed to generate a room id")
			fmt.Println("Failed to read from terminal")
			return fmt.Errorf("failed to read from terminal: %v", err)
		}
		if userMessage == "exit" {
			term.Write([]byte("---- Left the room -------- \n"))
			utils.ClearUserTerminal(term)
			// STILL NEED to take out the user of the room map and so on
			userMessage = fmt.Sprintf("user %s : left the room", user.UserTag)
			room.Users.Delete(user.UserTag)

			user.CurrentRoomName = ""
			room.UpdateRoomChat(userMessage, user.UserTag)
			return nil

		}
		if userMessage == "show_all_chatroom_users" {
			room.ShowAllUserInRoom(user.Term)
			continue
		}

		room.UpdateRoomChat(userMessage, user.UserTag)
	}
}

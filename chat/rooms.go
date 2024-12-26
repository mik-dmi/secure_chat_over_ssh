package chat

import (
	"fmt"
	"log"

	"sync"
	"time"

	"github.com/teris-io/shortid"
	"golang.org/x/term"
)

type Room struct {
	RoomName       string
	Users          []*User
	MessageHistory []*UserMessage
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

	// Iterate over the sync.Map to find the matching room
	rm.Rooms.Range(func(key, value interface{}) bool {
		if value == room {
			roomID = key.(string)
			found = true
			return false // Stop iteration
		}
		return true // Continue iteration
	})

	return roomID, found
}

// GetRoom retrieves a Room pointer by its ID.
func (rm *RoomManager) GetRoomByID(id string) (*Room, bool) {
	value, ok := rm.Rooms.Load(id)
	if !ok {
		return nil, false
	}
	return value.(*Room), true
}

// GetRoom retrieves a Room pointer by its Name
// GetRoomByName does a linear search of the sync.Map to find a room by its RoomName.
func (rm *RoomManager) GetRoomByName(name string) (*Room, bool) {
	var foundRoom *Room
	var found bool

	rm.Rooms.Range(func(key, value any) bool {
		room, ok := value.(*Room)
		if !ok {
			// Continue iterating if for some reason the value isn’t a *Room
			return true
		}
		if room.RoomName == name {
			foundRoom = room
			found = true
			// Return false to break out of the Range loop
			return false
		}
		// Continue iterating
		return true
	})

	return foundRoom, found
}

// DeleteRoom removes a Room from the sync.Map by its ID.
func (rm *RoomManager) DeleteRoom(id string) {
	rm.Rooms.Delete(id)
}

// ListRooms returns a slice of all rooms.
// sync.Map supports a Range method to iterate over entries.
func (rm *RoomManager) ListRooms() []*Room {
	rooms := make([]*Room, 0)
	rm.Rooms.Range(func(key, value interface{}) bool {
		if room, ok := value.(*Room); ok {
			rooms = append(rooms, room)
		}
		return true // keep iterating
	})
	return rooms
}

func (rm *RoomManager) JoinRoom(user *User, term *term.Terminal) *Room {
	for {
		term.Write([]byte("What's the ID of the room you want to join?\n"))
		idOfRoom, err := term.ReadLine()
		if err != nil {
			log.Printf("Error reading input: %v", err)
			continue
		}
		room, ok := rm.GetRoomByID(idOfRoom)
		if !ok {
			term.Write([]byte(fmt.Sprintf("The room ID %v does not exist\n", idOfRoom)))
			continue
		}
		// Add the user to the room's user list
		room.Users = append(room.Users, user)
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
			term.Write([]byte("Wrong command, try gain:\n"))
		}
		_, ok := rm.GetRoomByName(nameOfRoom)

		if ok {
			term.Write([]byte("Fail creating room, a room that name already exists\n"))
		} else {

			RoomID, err := shortid.Generate()
			if err != nil {
				log.Fatal(err)
			}

			room := &Room{
				RoomName: nameOfRoom,

				Users:          []*User{user},
				MessageHistory: []*UserMessage{},
			}

			rm.Rooms.Store(RoomID, room)
			return room
		}
	}
}

func (r *Room) UpdateRoomChat(userMessage string, userTag string) {
	currentTime := time.Now()
	formattedMessageTime := currentTime.Format("15:04")
	for _, user := range r.Users {
		if user.UserTag != userTag {
			user.Term.Write([]byte(fmt.Sprintf("%s at %s: %s\n", userTag, formattedMessageTime, userMessage)))
			//add the message  to the chat history obj
		}
	}

	r.MessageHistory = append(r.MessageHistory, &UserMessage{
		Message: userMessage,
		Time:    formattedMessageTime,
		UserTag: userTag,
	})
}

func (rm *RoomManager) WriteMessageToChat(term *term.Terminal, room *Room, userTag string) {

	for {
		userMessage, err := term.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		if userMessage == "exit" {
			term.Write([]byte("Left the room \n\n\n"))
			// STILL NEED to take out the user of the room map and so on
		}
		fmt.Println("Print: \n", userTag)
		room.UpdateRoomChat(userMessage, userTag)

	}
}

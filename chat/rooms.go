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
	Owner          string // who created the room and the only allow to delete it
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

func (rm *RoomManager) DeleteRoom(room *Room, h *SSHHandler) {
	rm.Rooms.Delete(room.RoomName)
	//find room object related to the "General Room"
	waitingRoom, ok := rm.GetRoomByID("0000")
	if !ok {
		log.Panicln("error getting Waiting Room room object from room manager")
	}

	room.Users.Range(func(key, value interface{}) bool {
		if u, ok := value.(*User); ok {

			utils.ClearUserTerminal(u.Term)
			u.CurrentRoomName = "Waiting room"
			waitingRoom.Users.Store(u.UserTag, u)
			u.Term.Write([]byte("---- Room was deleted --------\n"))
			if u.UserTag != room.Owner {
				u.Term.Write([]byte("> Press Enter To Continue\n"))

			}
		}

		return true
	})
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

func (rm *RoomManager) JoinRoom(user *User, term *term.Terminal) (string, *Room) {
	for {
		term.Write([]byte("What's the ID of the room you want to join?\n"))
		roomManagerMapID, err := term.ReadLine()
		if err != nil {
			log.Printf("Error reading input: %v", err)
			return "", nil
		}
		if roomManagerMapID == "exit" {
			utils.ClearUserTerminal(term)
			return "", nil
		}
		room, ok := rm.GetRoomByID(roomManagerMapID)
		if !ok {
			term.Write([]byte(fmt.Sprintf("The room ID %v does not exist\n", roomManagerMapID)))
			continue
		}
		user.CurrentRoomName = room.RoomName
		room.Users.Store(user.UserTag, user)

		msg := fmt.Sprintf("Room found! Welcome to %s room (Room ID: %s)\n", room.RoomName, roomManagerMapID)
		term.Write([]byte(msg))
		return roomManagerMapID, room
	}
}

/*
func (rm *RoomManager) AddUserToGeneralRoom(user *User, room *Room) error {

	room.Users.Store(user.UserTag, user)
	user.CurrentRoomName = "General Room"
	msg := fmt.Sprintf("Welcome to %s room\n", room.RoomName)
	user.Term.Write([]byte(msg))
	rm.GetIntoAGroupChat(user.Term, room)

	return nil
}
*/

func (rm *RoomManager) CreateRoom(user *User) (*Room, string) { // --> CHANGE rooms []Room   --< pass as pointer cause its needed to add something to the slice
	for {
		user.Term.Write([]byte("What is the name of the room you want to create?"))
		nameOfRoom, err := user.Term.ReadLine()
		if err != nil {
			user.Term.Write([]byte("Error reading terminal\n"))
			return nil, ""
		}
		//if user wants to go back to the initial menu
		if nameOfRoom == "exit" {
			utils.ClearUserTerminal(user.Term)
			return nil, ""
		}

		_, ok := rm.GetRoomByName(nameOfRoom)

		if ok {
			user.Term.Write([]byte("Fail creating room, a room that name already exists\n"))
		} else {

			roomManagerMapID, err := shortid.Generate()
			if err != nil {
				fmt.Println("Failed to generate a room id")
				user.Term.Write([]byte("Fail to generate a room id, try again\n"))
				continue
			}

			room := &Room{
				RoomName:       nameOfRoom,
				Users:          sync.Map{},
				MessageHistory: []*UserMessage{},
				Owner:          user.UserTag,
			}

			room.Users.Store(user.UserTag, user)

			rm.Rooms.Store(roomManagerMapID, room)
			return room, roomManagerMapID
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

func (rm *RoomManager) WriteMessageToChat(user *User, roomManagerMapID string, h *SSHHandler) error {

	room, ok := rm.GetRoomByID(roomManagerMapID)
	if !ok {
		log.Panicln("error getting General Room room object from room manager (WriteMessageToChat);  ", roomManagerMapID)
	}

	for {
		userMessage, err := user.Term.ReadLine()
		if err != nil {
			fmt.Println("Failed to read from terminal")
			return fmt.Errorf("failed to read from terminal: %v", err)
		}
		if userMessage == "exit" {
			user.Term.Write([]byte("---- Left the room -------- \n"))
			utils.ClearUserTerminal(user.Term)
			// STILL NEED to take out the user of the room map and so on
			userMessage = fmt.Sprintf("user %s : left the room", user.UserTag)
			room.Users.Delete(user.UserTag)

			user.CurrentRoomName = "Waiting room"
			room.UpdateRoomChat(userMessage, user.UserTag)
			return nil

		}
		if userMessage == "show_all_chatroom_users" {
			room.ShowAllUserInRoom(user.Term)
			continue
		}
		if userMessage == "delete_room" {
			if user.UserTag == room.Owner {
				// need to run throught all the users and change their currentRoom
				rm.DeleteRoom(room, h)

			} else {
				user.Term.Write([]byte("Can not delete room, you are not the owner\n"))

				continue
			}
			fmt.Println("Debug: userCurrent room after delete: ", user.CurrentRoomName)
			return nil

		}

		//working on this is a meesconst
		/*roomFromManager, ok := h.RoomManager.Rooms.Load(room.RoomName)
		if !ok {
			// Room not found
			return fmt.Errorf("error finding room in RoomManager")
		}

		roomFromManagerFinal, ok := roomFromManager.(*Room)
		if !ok {
			// The value stored is not of type *Room
			return fmt.Errorf("error asserting value of  RoomManager")
		}

		if roomFromManagerFinal.RoomName == "" {
			return nil
		}*/

		if user.CurrentRoomName == "Waiting room" {
			return nil
		}

		room.UpdateRoomChat(userMessage, user.UserTag)
	}
}

package utils

import (
	"fmt"
	"log"
	"secure_chat_over_ssh/chat"

	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
)

func NewUser(session ssh.Session, UsersManager *chat.UsersManager, term *term.Terminal) (*chat.User, error) {
	var userTag string
	for {
		term.Write([]byte("Enter a unique user tag:\n"))
		line, err := term.ReadLine()

		if err != nil {
			err = fmt.Errorf("error reading user tag: %s", err)
			return nil, err
		}
		userTag = line
		log.Printf(" 1 - userTag is: %s", userTag)

		// Validate userTag (non-empty, alphanumeric, etc.)
		if len(userTag) == 0 {
			term.Write([]byte("User tag cannot be empty. Try again:\n"))
			continue
		}
		if len(userTag) > 40 { // Example max length check
			term.Write([]byte("User tag must be 20 characters or less. Try again:\n"))
			continue
		}

		// Check if userTag already exists
		_, ok := UsersManager.Users.Load(userTag)
		if ok {
			term.Write([]byte(fmt.Sprintf("The user tag '%s' is already in use. Try again:\n", userTag)))
		} else {
			break
		}
	}
	log.Printf(" 2 - userTag is: %s", userTag)
	user := &chat.User{
		Session:         session,
		UserTag:         userTag,
		CurrentRoomName: "General Room",
		Term:            term,
	}

	return user, nil
}

// funtion to help populate a room for test purposes
/*
func PopulateRoom(session ssh.Session, roomMap *RoomManager) {
	// Create random users
	for i := 0; i < rand.Intn(10)+1; i++ { // Generate 1 to 10 users
		userTag := fmt.Sprintf("#%04d", rand.Intn(10000)) // Random 4-digit tag
		user := &users.User{Session: session, UserTag: userTag, CurrentRoom: room.RoomId}
		room.users = append(room.users, user)
	}
	// Create random messages
	for i := 0; i < rand.Intn(20)+5; i++ { // Generate 5 to 20 messages
		user := room.users[rand.Intn(len(room.users))] // Pick a random user
		message := fmt.Sprintf("Message %d content", i+1)
		timestamp := time.Now().Add(time.Duration(-rand.Intn(3600)) * time.Second).Format("15:04:05") // Random timestamp in last hour
		room.messageHistory = append(room.messageHistory, UserMessage{
			message: message,
			time:    timestamp,
			userTag: user.userTag,
		})
	}
}
*/
/*
func RoomExists(rooms []rRoom, selectedRoom rooms.Room) bool {

	for _, room := range rooms {
		if room.RoomId == selectedRoom.RoomId {
			return true
		}
	}
	return false
}
*/

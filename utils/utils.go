package utils

import "golang.org/x/term"

// ClearScreen clears the terminal screen in a cross-platform way.
func ClearUserTerminal(term *term.Terminal) {
	// ANSI escape code to clear the screen and reset the cursor
	term.Write([]byte("\033[H\033[2J"))
	term.Write([]byte("\033[3J")) // Clear scrollback buffer (optional)
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

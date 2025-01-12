--> check the * in the newRoom := Room{users: []*User{}...}
---> simplify the CreateRoom with a costructor maybe


----> in this :  populateRoom(session, &generalRoom)   --- check if session sobr be passed as a pointer 



-----> WriteMessageToChat NEEDS a return when the user writes Exit in the handlers it need to be known


------>> I dont thing func (rm *RoomManager) WriteMessageToChat(term *term.Terminal, room *Room, userTag string) {
      ------> needs to be a (rm *RoomManager)   just did it to try figure out something but it need to be thought again 
      
---> AMybe u don't need to pass the room and just pass the string witht eh ID and then u load from the sync map ( not sure but might)



-----> nEED TO MAKE THE MUX FOR THE CROUP CHAT WORK : 	MessageHistory []*UserMessage
	messagesMu     sync.Mutex
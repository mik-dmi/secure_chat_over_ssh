--> check the * in the newRoom := Room{users: []*User{}...}
---> simplify the CreateRoom with a costructor maybe


----> in this :  populateRoom(session, &generalRoom)   --- check if session sobr be passed as a pointer 



-----> WriteMessageToChat NEEDS a return when the user writes Exit in the handlers it need to be known


------>> I dont thing func (rm *RoomManager) WriteMessageToChat(term *term.Terminal, room *Room, userTag string) {
      ------> needs to be a (rm *RoomManager)   just did it to try figure out something but it need to be thought again 
      
---> AMybe u don't need to pass the room and just pass the string witht eh ID and then u load from the sync map ( not sure but might)



--------->>> DOOOO !!!!  If MessageHistory grows indefinitely, you may want to consider a size limit or implement batch processing for displaying older messages.



---> FIX this BuG that appears in the terminal after CRL + C:
You just joinned: General Room (Room ID: 0000)
General Room 
Hereeee
qwdscv at xx:xx: wdescv
> wdascv 
Hereeee
Failed to read from terminal
Choose an option by typing the corresponding command (cmd):
- Join General Chat Room (cmd: JGR)
- Create A Chat Room (cmd: CR)
- Join a Chat Room (cmd: JR)
> Connection to localhost closed.
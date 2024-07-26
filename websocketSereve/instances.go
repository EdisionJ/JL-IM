package websocketSereve

import (
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map/v2"
	"strconv"
)

type Room struct {
	RoomID      int64
	ActiveUsers int
	Users       map[int64]*websocket.Conn
}

type User struct {
	ID      int64
	Rooms   []int64
	Channel *websocket.Conn
}

// mapping from str(uid) to user instance
var userMap = cmap.New[*User]()

// mapping from str(rmid) to room instance
var roomMap = cmap.New[*Room]()

/**********************  Room  **********************/
func RoomMapSet(roomId int64, room *Room) {
	strUid := GetStrId(roomId)
	roomMap.Set(strUid, room)
}
func RoomMapGet(roomId int64) (room *Room, ok bool) {
	strUid := GetStrId(roomId)
	return roomMap.Get(strUid)
}
func RoomMapDel(roomId int64) {
	strUid := GetStrId(roomId)
	roomMap.Remove(strUid)
}
func JoinRoom(roomID, uid int64, conn *websocket.Conn) {
	room, ok := RoomMapGet(roomID)
	if !ok {
		room = &Room{
			RoomID: roomID,
			Users:  make(map[int64]*websocket.Conn, 0),
		}
		room.Users[uid] = conn
		room.ActiveUsers++
		RoomMapSet(roomID, room)
	} else {
		room.Users[uid] = conn
		room.ActiveUsers++
	}
}
func ExitFromRoom(roomID, uid int64) {
	room, _ := RoomMapGet(roomID)
	delete(room.Users, uid)
	room.ActiveUsers--
	if room.ActiveUsers == 0 {
		//清理不活跃的房间
		RoomMapDel(roomID)
	}
}

/**********************  User  **********************/
func UserMapSet(user *User) {
	strUid := GetStrId(user.ID)
	userMap.Set(strUid, user)
}
func UserMapGet(uid int64) (*User, bool) {
	strUid := GetStrId(uid)
	return userMap.Get(strUid)
}
func UserMapDel(uid int64) {
	strUid := GetStrId(uid)
	userMap.Remove(strUid)
}

func GetStrId(uid int64) string {
	return strconv.FormatInt(uid, 10)
}

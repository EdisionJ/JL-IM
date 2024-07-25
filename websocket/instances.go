package websocket

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"golang.org/x/net/websocket"
	"strconv"
)

type Room struct {
	RoomID int64
	//RoomInfo
	//Users
}

type User struct {
	ID      int64
	Info    RspModels.UserInfo
	Channel *websocket.Conn
}

// mapping from str(uid) to user instance
var userMap = cmap.New[*User]()

// mapping from str(rmid) to room instance
var roomMap = cmap.New[*Room]()

/**********************  Room  **********************/
func RoomMapAdd(roomId int64, room *Room) {
	strUid := GetStrId(roomId)
	roomMap.Set(strUid, room)
}
func RoomMapGet(roomId int64) (*Room, bool) {
	strUid := GetStrId(roomId)
	return roomMap.Get(strUid)
}
func RoomMapDel(roomId int64) {
	strUid := GetStrId(roomId)
	roomMap.Remove(strUid)
}

/**********************  User  **********************/
func UserMapAdd(user *User) {
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

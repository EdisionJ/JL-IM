package responseModels

type RoomInfo struct {
	ID     int64  `json:"id"`   // 房间id
	Type   int32  `json:"type"` // 房间类型  私聊：0  群聊：1
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

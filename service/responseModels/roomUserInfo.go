package responseModels

type RoomUerInfo struct {
	ID       int64  `json:"id"` // 用户id
	GroupID  int64  `json:"group_id"`
	Nickname string `json:"nickname"` // 聊天室内昵称
	Role     int32  `json:"role"`     // -1：已退出 0：普通成员  1：管理员 2：群主
	Ban      int32  `json:"ban"`      // 0：正常发言 1：被禁言
}

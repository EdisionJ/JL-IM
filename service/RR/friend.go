package RR

type AddFriend struct {
	Uid      int64  `json:"uid"`
	FriendId int64  `json:"friend_id"`
	ReqMsg   string `json:"req_msg"`
	Flag     int    `json:"flag"`
}

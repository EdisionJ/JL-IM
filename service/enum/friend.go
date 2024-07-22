package enum

const (
	FriendShipNormal  = 0 //正常好友关系
	FriendShipBaned   = 1 //黑名单
	FriendShipDeleted = 2 //已删除
)

const (
	FriendReqRefuse       = -1
	FriendReqNotYetAgreed = 0
	FriendReqAgree        = 1
)

const (
	UserStatusOffline = 0
	UserStatusOnline  = 1
)

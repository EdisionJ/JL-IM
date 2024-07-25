package enum

import (
	"IM/globle"
	"time"
)

const (
	CacheTime  = 3 * time.Hour
	User       = globle.Project + "-user::"
	UserFriend = globle.Project + "-userFriend::"
	UserApply  = globle.Project + "-userApply::"
	RoomFriend = globle.Project + "-roomFriend::"
	Contact    = globle.Project + "-contact::"
	Room       = globle.Project + "-room::"
)
const (
	// 房间缓存
	RoomCacheByID = Room + "%d"

	// 好友房间缓存
	RoomFriendCacheByRoomID          = RoomFriend + "%d"
	RoomFriendCacheByUidAndFriendUid = RoomFriend + "%d_%d"

	// 用户缓存
	UserCacheByID   = User + "%d"
	UserCacheByName = User + "%s"

	// 用户好友缓存
	UserFriendCacheByUidAndFriendUid = UserFriend + "%d_%d"

	// 好友申请缓存
	UserApplyCacheByUidAndFriendUid = UserApply + "%d_%d"
	UserApplyCacheByFriendUid       = UserApply + "%d"

	// 会话缓存
	ContactCacheById = Contact + "%d"
)

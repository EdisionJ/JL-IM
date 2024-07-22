package enum

// event
const (
	UserLogSignUp  = "UserLogSignUp"
	UserLogInOrOut = "UserLogInOrOut"
	UserFriendReq  = "UserFriendReq"
	UserFriendAdd  = "UserFriendAdd"
)

// group
const (
	UserLogSignUpGroup     = UserLogSignUp + "-G"
	UserLogInOrOutGroup    = "UserLogInOrOut" + "-G"
	UserFriendServiceGroup = "UserFriendService" + "-G"
)

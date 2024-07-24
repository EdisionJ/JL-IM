package globle

type UserInfo struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	SelfDescribe string `json:"self_describe"`
	PhoneNumber  string `json:"phone_number"`
	Email        string `json:"email"`
}

package responseModels

type UserInfo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

package RR

type UserSingUp struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email       string `json:"email" valid:"email"`
	PassWD      string `json:"passwd" gorm:"size:64"`
	RePassWD    string `json:"re_passwd"`
}

type UserLogIn struct {
	PhoneNumber string `json:"phone_number" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Email       string `json:"email" valid:"email"`
	PassWD      string `json:"passwd" gorm:"size:64"`
}

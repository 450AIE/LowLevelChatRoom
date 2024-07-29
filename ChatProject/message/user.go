package message

type User struct {
	UserID     string `json:"userID"`
	UserPwd    string `json:"userPwd"`
	UserStatus int    `json:"userStatus"`
}

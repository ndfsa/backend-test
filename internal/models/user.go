package models

type User struct {
    UserId   uint64 `json:"id"`
	Fullname string `json:"fullname"`
	Username string `json:"username"`
}

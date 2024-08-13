package models

type User struct {
	Username     string
	Password     string
	PasswordHint string
	TOTPSecret   string
	CreatedDate  string
	BlabName     string
	RealName     string
	LastLogin    string
}

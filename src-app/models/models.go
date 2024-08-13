package models

import (
	"log"
	"time"
	"verademo-go/src-app/shared/utils"

	"github.com/pquerna/otp/totp"
)

func Timestamp(postDate string) string {
	date, err := time.Parse(time.RFC3339, postDate)
	timestamp := date.Format("Jan 2, 2006")
	if err != nil {
		timestamp = "Error getting date:\n" + err.Error()
	}
	return timestamp
}

func CreateUser(username string, blabName string, realName string) User {
	password := username
	dateCreated := time.Now().Format(time.RFC3339)
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "VeraDemo",
		AccountName: username,
	})
	if err != nil {
		log.Println("Failed to generate TOTP secret.")
		return User{Username: username, Password: utils.GetMD5Hash(password), PasswordHint: username, CreatedDate: dateCreated, BlabName: blabName, RealName: realName}
	}
	return User{Username: username, Password: utils.GetMD5Hash(password), PasswordHint: username, TOTPSecret: secret.Secret(), CreatedDate: dateCreated, BlabName: blabName, RealName: realName}
}

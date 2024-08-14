package models

import "verademo-go/src-app/shared/utils"

type Blabber struct {
	Id              int
	Username        string `db:"username"`
	RealName        string `db:"real_name"`
	BlabName        string `db:"blab_name"`
	CreatedDate     string `db:"created_at"`
	NumberListeners int
	NumberListening int
}

func (blabber *Blabber) GetProfileImageFromUsername() string {
	return utils.GetProfileImageFromUsername(blabber.Username)
}

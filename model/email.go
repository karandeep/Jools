package model

import (
	"../lib"
)

type Email struct {
	EncUserId string
	Network   int
	Emails    string
	Names     string
}

func GetEmailData(network int, encUserId string) (Email, error) {
	conn := lib.GetDBConnection()
	row := conn.QueryRow("SELECT encUserId,network,emails,names FROM Emails WHERE encUserId = ? AND network = ?", encUserId, network)
	var email Email
	err := row.Scan(
		&email.EncUserId,
		&email.Network,
		&email.Emails,
		&email.Names,
	)
	return email, err
}

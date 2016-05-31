package model

import (
	"../lib"
	"log"
)

const MAX_ADDRESSES_PER_USER = 10

type Address struct {
	Id      int64
	UserId  int64
	Email   string
	Name    string
	Address string
	City    string
	State   string
	Pincode int
	Mobile  string
	Active  int //To allow user to remove addresses. Address can't be deleted from sql table itself coz an order can map to old addresses which the user may have removed.
}

func (address Address) Store() (int64, error) {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT INTO Address (userId, email, name, address, city, state, pincode, mobile) VALUES(?,?,?,?,?,?,?,?)")
	res, err := stmt.Exec(address.UserId, address.Email, address.Name, address.Address, address.City, address.State, address.Pincode, address.Mobile)
	lastId, err := res.LastInsertId()
	return lastId, err
}

func GetAddressFromId(addressId int64) (Address, error) {
	var address Address
	conn := lib.GetDBConnection()
	row := conn.QueryRow("SELECT * FROM Address WHERE id = ?", addressId)
	err := row.Scan(
		&address.Id,
		&address.UserId,
		&address.Email,
		&address.Name,
		&address.Address,
		&address.City,
		&address.State,
		&address.Pincode,
		&address.Mobile,
		&address.Active,
	)
	return address, err
}

func GetAddressesForUser(email string, userId int64) ([MAX_ADDRESSES_PER_USER]Address, error) {
	var addresses [MAX_ADDRESSES_PER_USER]Address
	conn := lib.GetDBConnection()
	rows, err := conn.Query("SELECT * FROM Address WHERE (email = ? OR userId = ?) AND active = 1 ORDER by id DESC LIMIT ?", email, userId, MAX_ADDRESSES_PER_USER)
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&addresses[index].Id,
			&addresses[index].UserId,
			&addresses[index].Email,
			&addresses[index].Name,
			&addresses[index].Address,
			&addresses[index].City,
			&addresses[index].State,
			&addresses[index].Pincode,
			&addresses[index].Mobile,
			&addresses[index].Active,
		); err != nil {
			return addresses, err
		}
		index++
	}
	return addresses, err
}

func GetAllAddressesForUser(email string, userId int64) ([MAX_ADDRESSES_PER_USER]Address, error) {
	var addresses [MAX_ADDRESSES_PER_USER]Address
	conn := lib.GetDBConnection()
	rows, err := conn.Query("SELECT * FROM Address WHERE email = ? OR userId = ? ORDER by id DESC LIMIT ?", email, userId, MAX_ADDRESSES_PER_USER)
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&addresses[index].Id,
			&addresses[index].UserId,
			&addresses[index].Email,
			&addresses[index].Name,
			&addresses[index].Address,
			&addresses[index].City,
			&addresses[index].State,
			&addresses[index].Pincode,
			&addresses[index].Mobile,
			&addresses[index].Active,
		); err != nil {
			return addresses, err
		}
		index++
	}
	return addresses, err
}

func RemoveAddressForUser(addressId int64, email string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Address SET active = 0 WHERE id = ? AND email = ? LIMIT 1")
	if err != nil {
		log.Println(err)
	}
	_, err = stmt.Exec(addressId, email)
	if err != nil {
		log.Println(err)
	}

	return err
}

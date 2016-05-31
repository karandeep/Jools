package model

import (
	"../lib"
	"time"
)

const (
	PAYMENT_CASH = iota
	PAYMENT_CREDIT_CARD
	PAYMENT_DEBIT_CARD
	PAYMENT_NETBANKING
)

const (
	ORDER_STARTED = iota
	ORDER_RECEIVED
	ORDER_SHIPPED
	ORDER_COMPLETED
)

const DELIVERY_LAG = 86400 * 14
const MAX_ORDERS_TO_SHOW = 10

type Order struct {
	Id            int64
	UserId        int64
	Email         string
	AddressId     int64
	Products      string
	PaymentMethod int
	Created       int32
	DeliveryYear  int
	DeliveryMonth string //Expected date
	DeliveryDay   int
	ShippedOn     int32
	CompletedOn   int32
	Cost          float64
	Paid          float64
	Status        int
}

func (order Order) Begin() (Order, error) {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT INTO Orders (userId, email, addressId, products, paymentMethod, created, cost, paid, status) VALUES(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return order, err
	}
	order.Created = lib.GetCurrentTimestamp()
	order.Status = ORDER_STARTED

	res, err := stmt.Exec(order.UserId, order.Email, order.AddressId, order.Products, order.PaymentMethod, order.Created, order.Cost, order.Paid, order.Status)
	if err != nil {
		return order, err
	}

	lastId, err := res.LastInsertId()
	if err != nil {
		return order, err
	}
	order.Id = lastId
	return order, nil
}

func GetOrder(orderId int64) (Order, Address, error) {
	var order Order
	var address Address
	conn := lib.GetDBConnection()
	row := conn.QueryRow("SELECT * FROM Orders WHERE id = ?", orderId)
	err := row.Scan(
		&order.Id,
		&order.UserId,
		&order.Email,
		&order.AddressId,
		&order.Products,
		&order.PaymentMethod,
		&order.Created,
		&order.ShippedOn,
		&order.CompletedOn,
		&order.Cost,
		&order.Paid,
		&order.Status,
	)
	if err == nil {
		order.DeliveryYear, order.DeliveryMonth, order.DeliveryDay = GetExpectedDelivery(order.Created)
		address, err = GetAddressFromId(order.AddressId)
	}
	return order, address, err
}

func GetExpectedDelivery(from int32) (int, string, int) {
	var startTime time.Time
	if from < 0 {
		startTime = time.Now()
	} else {
		startTime = lib.GetTimeFromTimestamp(from)
	}
	return lib.TimeAfterDuration(DELIVERY_LAG, startTime)
}

func GetOrdersForUser(email string, userId int64) ([MAX_ORDERS_TO_SHOW]Order, error) {
	var orders [MAX_ORDERS_TO_SHOW]Order
	conn := lib.GetDBConnection()
	rows, err := conn.Query("SELECT * FROM Orders WHERE email = ? OR userId = ? ORDER by id DESC LIMIT ?", email, userId, MAX_ORDERS_TO_SHOW)
	if err != nil {
		return orders, err
	}
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&orders[index].Id,
			&orders[index].UserId,
			&orders[index].Email,
			&orders[index].AddressId,
			&orders[index].Products,
			&orders[index].PaymentMethod,
			&orders[index].Created,
			&orders[index].ShippedOn,
			&orders[index].CompletedOn,
			&orders[index].Cost,
			&orders[index].Paid,
			&orders[index].Status,
		); err != nil {
			return orders, err
		}
		index++
	}
	return orders, err
}

package main

import (
	"../lib"
	"log"
	"strconv"
)

func main() {
	conn := lib.GetDBConnection()
    var maxId int64
    row := conn.QueryRow("SELECT MAX(id) AS maxId FROM Product")
    err := row.Scan(&maxId)
    if err != nil {
        log.Println(err)
        return
    }
    
	stmt, err := conn.Prepare("UPDATE Product SET encId = ?, created = ? WHERE id = ? AND encId = ''")
    if err != nil {
        log.Println(err)
        return
    }
    
	var i int64
    curTimestamp := lib.GetCurrentTimestamp()
	for i = 1; i <= maxId; i++ {
        encrypted, _ := lib.Encrypt(strconv.FormatInt(i, 10))
		stmt.Exec(encrypted, curTimestamp, i)
		if i%100 == 0 {
			log.Println("Processed", i)
		}
	}
}

package main

import (
	"../lib"
	"log"
	"strconv"
)

func main() {
	conn := lib.GetDBConnection()
	stmt, _ := conn.Prepare("UPDATE Inspiration SET encId = ? WHERE id = ?")
	var i int64
	for i = 1; i < 856; i++ {
        encrypted, _ := lib.Encrypt(strconv.FormatInt(i, 10))
		stmt.Exec(encrypted, i)
		if i%100 == 0 {
			log.Println("Processed", i)
		}
	}
}

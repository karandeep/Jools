package main

import (
	"../lib"
	"log"
)

func main() {
	conn := lib.GetDBConnection()
	
	stmt, _ := conn.Prepare("UPDATE Inspiration SET viewedCount = viewedCount + ?, favoritedCount = favoritedCount + ? WHERE id = ?")
	var i int64
	for i = 1; i < 856; i++ {
		randomIntOne := lib.GetRandomInt()
		viewedCount := 20 + randomIntOne % 100
		randomIntTwo := lib.GetRandomInt()
		favCount := 5 + randomIntTwo % 15
		stmt.Exec(viewedCount, favCount, i)
		if i%100 == 0 {
			log.Println("Processed", i)
		}
	}
}


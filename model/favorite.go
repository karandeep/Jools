package model

import (
	"../lib"
	"database/sql"
	"strings"
)

type Favorite struct {
	UserId       int64
	Inspirations string
}

func GetFavsForUser(userId int64) (string, error) {
	conn := lib.GetDBConnection()
	row := conn.QueryRow("SELECT inspirations FROM Favorite WHERE userId = ?", userId)
	var inspirations string
	err := row.Scan(&inspirations)
	return inspirations, err
}

func updateFavsWithRemoval(userId int64, inspirations string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Favorite SET inspirations = ? WHERE userId = ? LIMIT 1")
	_, err = stmt.Exec(inspirations, userId)
	return err
}

func updateFavs(userId int64, inspirations string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Favorite SET inspirations = CONCAT(inspirations, ?) WHERE userId = ? LIMIT 1")
	_, err = stmt.Exec(inspirations, userId)
	return err
}

func insertFavsForUser(userId int64, inspirations string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT INTO Favorite(userId, inspirations) VALUES (?,?)")
	_, err = stmt.Exec(userId, inspirations)
	return err
}

func FetchFavoriteInspirations(userId int64) (string, error) {
	curFavs, err := GetFavsForUser(userId)
	if err != nil {
		return "", err
	}
	return curFavs, nil
}

func SyncFavoriteInspirations(userId int64, inspirations string, removeInspirations string) (string, error) {
	var finalFavs string
	curFavs, err := GetFavsForUser(userId)
	if err != nil {
		if err == sql.ErrNoRows {
			err = insertFavsForUser(userId, inspirations)
			if err == nil {
				encIds := lib.ReformatEncIdString(inspirations)
				IncrementInspirationFavorites(encIds)
			}
		}
		return inspirations, err
	}

	if removeInspirations != "" {
		removedList := strings.Split(removeInspirations, ";")
		removedListLen := len(removedList)
		removedMap := make(map[string]int)
		for i := 0; i < removedListLen; i++ {
			if removedList[i] == "" {
				continue
			}
			removedMap[removedList[i]] = 1
		}

		curList := strings.Split(curFavs, ";")
		curListLen := len(curList)
		for i := 0; i < curListLen; i++ {
			if curList[i] == "" {
				continue
			}
			_, found := removedMap[curList[i]]
			if !found {
				finalFavs += finalFavs + ";"
			}
		}
		finalFavs += inspirations
		err = updateFavsWithRemoval(userId, finalFavs)
		if err == nil {
			encIds := lib.ReformatEncIdString(inspirations)
			IncrementInspirationFavorites(encIds)

			encIds = lib.ReformatEncIdString(removeInspirations)
			DecrementInspirationFavorites(encIds)
		}
		return finalFavs, err
	} else {
		err = updateFavs(userId, inspirations)
		if err == nil {
			encIds := lib.ReformatEncIdString(inspirations)
			IncrementInspirationFavorites(encIds)
			finalFavs = curFavs + ";" + inspirations
		}
	}

	return finalFavs, err
}

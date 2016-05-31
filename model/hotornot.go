package model

import (
	"../lib"
)

const RANDOM_IMAGE_COUNT = 2

type HotBattles struct {
	Id     int
	Winner int
	Loser  int
}

type HotImages struct {
	Id       int
	EncId    string
	Filename string
	Score    int
	Wins     int
	Losses   int
}

func GetRandomImages() ([]HotImages, error) {
	conn := lib.GetHotOrNotDBConnection()
	rows, err := conn.Query("SELECT encId,filename,score,wins,losses FROM images ORDER BY RAND() LIMIT 0,?", RANDOM_IMAGE_COUNT)
	randomImages := make([]HotImages, 2)

	if err != nil {
		return randomImages, err
	}

	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&randomImages[index].EncId,
			&randomImages[index].Filename,
			&randomImages[index].Score,
			&randomImages[index].Wins,
			&randomImages[index].Losses,
		); err != nil {
			return randomImages, err
		}
		index++
	}
	if err = rows.Err(); err != nil {
		return randomImages, err
	}
	return randomImages, nil
}

func GetTopRated(howMany int) ([]HotImages, error) {
	conn := lib.GetHotOrNotDBConnection()
	rows, err := conn.Query("SELECT encId,filename,score,wins,losses FROM images ORDER BY score DESC LIMIT 0,?", howMany)
	randomImages := make([]HotImages, howMany)

	if err != nil {
		return randomImages, err
	}

	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&randomImages[index].EncId,
			&randomImages[index].Filename,
			&randomImages[index].Score,
			&randomImages[index].Wins,
			&randomImages[index].Losses,
		); err != nil {
			return randomImages, err
		}
		index++
	}
	if err = rows.Err(); err != nil {
		return randomImages, err
	}
	return randomImages, nil
}

func Rate(winner, loser int) error {
	conn := lib.GetHotOrNotDBConnection()
	row := conn.QueryRow("SELECT id,encId,filename,score,wins,losses FROM images WHERE id = ?", winner)

	var winnerInfo, loserInfo HotImages
	err := row.Scan(
		&winnerInfo.Id,
		&winnerInfo.EncId,
		&winnerInfo.Filename,
		&winnerInfo.Score,
		&winnerInfo.Wins,
		&winnerInfo.Losses,
	)

	row = conn.QueryRow("SELECT id,encId,filename,score,wins,losses FROM images WHERE id = ?", loser)

	err = row.Scan(
		&loserInfo.Id,
		&loserInfo.EncId,
		&loserInfo.Filename,
		&loserInfo.Score,
		&loserInfo.Wins,
		&loserInfo.Losses,
	)

	winnerExpected := lib.Expected(float64(loserInfo.Score), float64(winnerInfo.Score))
	winnerNewScore := lib.Win(float64(winnerInfo.Score), winnerExpected)

	loserExpected := lib.Expected(float64(winnerInfo.Score), float64(loserInfo.Score))
	loserNewScore := lib.Loss(float64(loserInfo.Score), loserExpected)

	stmt, err := conn.Prepare("UPDATE images SET score = ?, wins = wins + 1 WHERE id = ?")
	_, err = stmt.Exec(winnerNewScore, winner)
	if err != nil {
		return err
	}
	stmt, err = conn.Prepare("UPDATE images SET score = ?, losses = losses + 1 WHERE id = ?")
	_, err = stmt.Exec(loserNewScore, loser)
	if err != nil {
		return err
	}

	stmt, err = conn.Prepare("INSERT INTO battles (winner, loser) VALUES (?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(winner, loser)
	return err
}

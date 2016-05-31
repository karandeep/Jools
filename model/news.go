package model

import (
	"../lib"
	"database/sql"
)

type News struct {
	Id      int64
	LangId  int64
	PaperId int64
	Title   string
	News    string
}

func scanAllNewsFields(rows *sql.Rows, news []News) error {
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&news[index].Id,
			&news[index].LangId,
			&news[index].PaperId,
			&news[index].Title,
			&news[index].News,
		); err != nil {
			return err
		}
		index++
	}
	return nil
}

func GetNews() ([]News, error) {
	conn := lib.GetNewsDBConnection()

	var news []News
	rows, err := conn.Query("SELECT id,langId,paperId,title,news FROM news")
	if err != nil {
		return news, err
	}
	news = make([]News, 10)
	err = scanAllNewsFields(rows, news)

	return news, err
}

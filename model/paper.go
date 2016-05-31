package model

import (
	"../lib"
	"database/sql"
)

const PAPER_COUNT = 5

type Paper struct {
	Id     int64
	LangId int64
	Name   string
	NameEn string
}

func scanAllPaperFields(rows *sql.Rows, papers []Paper) error {
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&papers[index].Id,
			&papers[index].LangId,
			&papers[index].Name,
			&papers[index].NameEn,
		); err != nil {
			return err
		}
		index++
	}
	return nil
}

func GetPapers() ([]Paper, error) {
	conn := lib.GetNewsDBConnection()

	var papers []Paper
	rows, err := conn.Query("SELECT id,langId,name,nameEn FROM paper")
	if err != nil {
		return papers, err
	}
	papers = make([]Paper, PAPER_COUNT)
	err = scanAllPaperFields(rows, papers)

	return papers, err
}

package model

import (
	"../lib"
	"database/sql"
)

const LANG_COUNT = 3

type Language struct {
	Id   int64
	Name string
}

func scanAllLanguageFields(rows *sql.Rows, languages []Language) error {
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&languages[index].Id,
			&languages[index].Name,
		); err != nil {
			return err
		}
		index++
	}
	return nil
}

func GetLanguages() ([]Language, error) {
	conn := lib.GetNewsDBConnection()

	var languages []Language
	rows, err := conn.Query("SELECT id,name FROM language")
	if err != nil {
		return languages, err
	}
	languages = make([]Language, LANG_COUNT)
	err = scanAllLanguageFields(rows, languages)

	return languages, err
}

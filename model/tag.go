package model

import (
	"../lib"
	"errors"
	"log"
	"strings"
)

type Tag struct {
	Id       int
	Category string
	Tag      string
}

func GetAllTags() ([]Tag, error) {
	tagList := make([]Tag, 1)
	conn := lib.GetDBConnection()
	rows, err := conn.Query("SELECT * FROM Tag")
	if err != nil {
		log.Println(err)
		return tagList, err
	}
	var curTag Tag
	for rows.Next() {
		if err := rows.Scan(
			&curTag.Id,
			&curTag.Category,
			&curTag.Tag,
		); err != nil {
			log.Println(err)
			return tagList, err
		}
		tagList = append(tagList, curTag)
	}
	return tagList, err
}

func AddTag(category, tag string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT INTO Tag(category, tag) VALUES (?,?)")
	if err != nil {
		return err
	}
	category = SanitizeForTagging(category)
	tag = SanitizeForTagging(tag)
	if category == "" || tag == "" {
		return errors.New("Invalid value for category or tag")
	}
	_, err = stmt.Exec(category, tag)
	return err
}

func SanitizeForTagging(tag string) string {
	return strings.ToLower(strings.Trim(tag, " ,"))
}

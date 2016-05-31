package model

import (
	"../lib"
)

type Comment struct {
	Id          int
	UserEncId   string
	UserName    string
	SubjectId   int
	SubjectType int
	Comment     string
	Created     int32
}

const MAX_COMMENT_COUNT = 20

const COMMENT_TYPE_INSPIRATION = 1

func AddComment(comment Comment) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT INTO Comment (userEncId, userName, subjectId, subjectType, comment, created) VALUES (?,?,?,?,?,?)")
	_, err = stmt.Exec(
		comment.UserEncId,
		comment.UserName,
		comment.SubjectId,
		comment.SubjectType,
		comment.Comment,
		comment.Created,
	)
	return err
}

func GetComments(subjectType int, subjectId int) ([]Comment, error) {
	conn := lib.GetDBConnection()
	rows, err := conn.Query("SELECT userEncId,userName,comment,subjectId,subjectType,created FROM Comment WHERE subjectType = ? AND subjectId = ? ORDER BY created ASC LIMIT ?", subjectType, subjectId, MAX_COMMENT_COUNT)
	comments := make([]Comment, MAX_COMMENT_COUNT)
	index := 0
	if err != nil {
		return comments, err
	}
	for rows.Next() {
		if err := rows.Scan(
			&comments[index].UserEncId,
			&comments[index].UserName,
			&comments[index].Comment,
			&comments[index].SubjectId,
			&comments[index].SubjectType,
			&comments[index].Created,
		); err != nil {
			return comments, err
		}
		index++
	}
	return comments, err
}

package model

import (
	"../lib"
	"bytes"
	"database/sql"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Inspiration struct {
	Id                int
	EncId             string
	Name              string
	Description       string
	ImageName         string
	UploaderEncId     string
	UploaderName      string
	UploadedAt        int
	FavoritedCount    int
	ViewedCount       int
	CommentCount      int
	LastComment       string
	Tags              string
	LastCommentorId   string
	LastCommentorName string
	HotEncId          string
}

type Upload struct {
	Id               int
	EncUserId        string
	ImageName        string
	UploadSuccessful int
	Approved         int
	Created          int32
	AdminActionDate  int
}

const INITIAL_INSPIRATION_COUNT = 10
const INSPIRATION_LOT_SIZE = 15
const UPLOAD_LOT_SIZE = 30
const INSPIRATION_FETCH_FIELDS = "encId,name,description,imageName,uploaderEncId,uploaderName,uploadedAt,favoritedCount,viewedCount,commentCount,lastComment,lastCommentorId,lastCommentorName,tags,hotEncId"
const USER_INITIAL_INSPIRATION_LOT = 50
const INSPIRATION_APPROVED = 1
const INSPIRATION_REJECTED = -1

const (
	LATEST_INSPIRATIONS = iota
	RANDOM_INSPIRATIONS
)

func scanAllInspirationFieldsFromRow(row *sql.Row, inspiration *Inspiration) error {
	if err := row.Scan(
		&inspiration.EncId,
		&inspiration.Name,
		&inspiration.Description,
		&inspiration.ImageName,
		&inspiration.UploaderEncId,
		&inspiration.UploaderName,
		&inspiration.UploadedAt,
		&inspiration.FavoritedCount,
		&inspiration.ViewedCount,
		&inspiration.CommentCount,
		&inspiration.LastComment,
		&inspiration.LastCommentorId,
		&inspiration.LastCommentorName,
		&inspiration.Tags,
		&inspiration.HotEncId,
	); err != nil {
		return err
	}
	return nil
}

func scanAllInspirationFields(rows *sql.Rows, inspirations []Inspiration) error {
	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&inspirations[index].EncId,
			&inspirations[index].Name,
			&inspirations[index].Description,
			&inspirations[index].ImageName,
			&inspirations[index].UploaderEncId,
			&inspirations[index].UploaderName,
			&inspirations[index].UploadedAt,
			&inspirations[index].FavoritedCount,
			&inspirations[index].ViewedCount,
			&inspirations[index].CommentCount,
			&inspirations[index].LastComment,
			&inspirations[index].LastCommentorId,
			&inspirations[index].LastCommentorName,
			&inspirations[index].Tags,
			&inspirations[index].HotEncId,
		); err != nil {
			return err
		}
		index++
	}
	return nil
}

func GetInitialInspirationList(listType int) ([]Inspiration, error) {
	conn := lib.GetDBConnection()
	inspirations := make([]Inspiration, INITIAL_INSPIRATION_COUNT)

	row := conn.QueryRow("SELECT MAX(id) AS id FROM Inspiration")
	var maxId, startId int
	err := row.Scan(&maxId)
	if err != nil {
		return inspirations, err
	}
	if listType == LATEST_INSPIRATIONS {
		startId = maxId
	} else if listType == RANDOM_INSPIRATIONS {
		startId = lib.GetRandomIntInRange(maxId/2, maxId)
	}
	rows, err := conn.Query("SELECT "+INSPIRATION_FETCH_FIELDS+" FROM Inspiration WHERE id <= ? ORDER BY id DESC LIMIT ?", startId, INITIAL_INSPIRATION_COUNT)
	if err != nil {
		return inspirations, err
	}

	err = scanAllInspirationFields(rows, inspirations)
	return inspirations, err
}

func GetNextInspirationLot(lastId int) ([]Inspiration, error) {
	limitId := lastId - 50 //Since ids of inspirations uploaded may be missing values
	conn := lib.GetDBConnection()
	inspirations := make([]Inspiration, INSPIRATION_LOT_SIZE)
	rows, err := conn.Query("SELECT "+INSPIRATION_FETCH_FIELDS+" FROM Inspiration WHERE id < ? AND id > ? ORDER BY id DESC LIMIT ?", lastId, limitId, INSPIRATION_LOT_SIZE)
	if err != nil {
		return inspirations, err
	}

	err = scanAllInspirationFields(rows, inspirations)
	if err != nil {
		return inspirations, err
	}

	if inspirations[0].EncId == "" {
		//No matches were found. Try once more without limiting - in case there are 50 consecutive ids missing
		rows, err = conn.Query("SELECT "+INSPIRATION_FETCH_FIELDS+" FROM Inspiration WHERE id < ? ORDER BY id DESC", lastId)
		if err != nil {
			return inspirations, err
		}
		err = scanAllInspirationFields(rows, inspirations)
	}
	return inspirations, err
}

func GetInspirationsForUser(uploaderId int) ([]Inspiration, error) {
	conn := lib.GetDBConnection()
	inspirations := make([]Inspiration, USER_INITIAL_INSPIRATION_LOT)
	rows, err := conn.Query("SELECT "+INSPIRATION_FETCH_FIELDS+" FROM Inspiration WHERE uploaderId = ? ORDER BY id DESC LIMIT ?", uploaderId, USER_INITIAL_INSPIRATION_LOT)
	if err != nil {
		return inspirations, err
	}

	err = scanAllInspirationFields(rows, inspirations)
	return inspirations, err
}

func GetInspiration(id int) (Inspiration, error) {
	conn := lib.GetDBConnection()
	var inspiration Inspiration
	row := conn.QueryRow("SELECT "+INSPIRATION_FETCH_FIELDS+" FROM Inspiration WHERE id = ?", id)

	err := scanAllInspirationFieldsFromRow(row, &inspiration)
	return inspiration, err
}

func IncrementInspirationViews(encIds string) (string, error) {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Inspiration SET viewedCount = viewedCount + 1 WHERE encId IN (" + encIds + ")")
	_, err = stmt.Exec()
	return "", err
}

func IncrementInspirationFavorites(encIds string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Inspiration SET favoritedCount = favoritedCount + 1 WHERE encId IN (" + encIds + ")")
	_, err = stmt.Exec()
	return err
}

func DecrementInspirationFavorites(encIds string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Inspiration SET favoritedCount = favoritedCount - 1 WHERE encId IN (" + encIds + ") AND favoritedCount > 0")
	_, err = stmt.Exec()
	return err
}

func UpdateInspirationLastComment(comment Comment) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("UPDATE Inspiration SET lastCommentorId = ?, lastCommentorName = ?, lastComment = ?, commentCount = commentCount + 1 WHERE id = ? LIMIT 1")
	_, err = stmt.Exec(comment.UserEncId, comment.UserName, comment.Comment, comment.SubjectId)
	return err
}

func GetUploadsForApproval() ([]Upload, error) {
	conn := lib.GetUploadDBConnection()
	uploads := make([]Upload, UPLOAD_LOT_SIZE)

	rows, err := conn.Query("SELECT * FROM Uploads WHERE approved = 0 LIMIT ?", UPLOAD_LOT_SIZE)
	if err != nil {
		return uploads, err
	}

	index := 0
	for rows.Next() {
		if err := rows.Scan(
			&uploads[index].Id,
			&uploads[index].EncUserId,
			&uploads[index].ImageName,
			&uploads[index].UploadSuccessful,
			&uploads[index].Approved,
			&uploads[index].Created,
			&uploads[index].AdminActionDate,
		); err != nil {
			return uploads, err
		}
		index++
	}
	return uploads, err
}

func MarkUpload(id int, approveStatus int) error {
	conn := lib.GetUploadDBConnection()
	row := conn.QueryRow("SELECT * FROM Uploads WHERE id = ?", id)
	var upload Upload
	err := row.Scan(
		&upload.Id,
		&upload.EncUserId,
		&upload.ImageName,
		&upload.UploadSuccessful,
		&upload.Approved,
		&upload.Created,
		&upload.AdminActionDate,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	finalDestination := "images/inspirations/" + upload.ImageName
	smallDestination := "images/inspirations/small/" + upload.ImageName
	if approveStatus == INSPIRATION_REJECTED {
		finalDestination = "imageServer/deleted/" + upload.ImageName
	}
	var out bytes.Buffer
	cmd := exec.Command("mv", "imageServer/uploads/"+upload.ImageName, finalDestination)
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		log.Println(out.String(), err)
		return err
	}
	if approveStatus == INSPIRATION_APPROVED {
		cmd = exec.Command("cp", finalDestination, smallDestination)
		cmd.Stderr = &out
		err = cmd.Run()
		if err != nil {
			log.Println(out.String(), err)
			return err
		}
		cmd = exec.Command("mogrify", "-resize", "320", smallDestination)
		cmd.Stderr = &out
		err = cmd.Run()
		if err != nil {
			log.Println(out.String(), err)
			return err
		}
	}

	curDate := lib.GetTrackingDate()
	stmt, err := conn.Prepare("UPDATE Uploads SET adminActionDate = ?, approved = ? WHERE id = ? LIMIT 1")
	if err != nil {
		log.Println(err)
		return err
	}
	_, err = stmt.Exec(curDate, approveStatus, id)

	if approveStatus == INSPIRATION_REJECTED || err != nil {
		log.Println(err)
		return err
	}
	connHot := lib.GetHotOrNotDBConnection()
	stmt, err = connHot.Prepare("INSERT INTO images (filename) VALUES (?)")
	res, err := stmt.Exec(upload.ImageName)
	if err != nil {
		log.Println(err)
		return err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	hotEncId, _ := lib.Encrypt(strconv.FormatInt(lastId, 10))
	stmt, err = connHot.Prepare("UPDATE images SET encId = ? WHERE id = ? LIMIT 1")
	res, err = stmt.Exec(hotEncId, lastId)
	if err != nil {
		log.Println(err)
		return err
	}

	connMain := lib.GetDBConnection()
	userData := GetUserByEncId(upload.EncUserId)
	stmt, err = connMain.Prepare("INSERT INTO Inspiration (imageName, uploaderEncId, uploaderName, uploadedAt, hotEncId) VALUES (?,?,?,?,?)")
	res, err = stmt.Exec(upload.ImageName, userData.EncId, userData.Name, upload.Created, hotEncId)
	if err != nil {
		log.Println(err)
		return err
	}
	lastId, err = res.LastInsertId()
	if err != nil {
		log.Println(err)
		return err
	}
	inspirationEncId, _ := lib.Encrypt(strconv.FormatInt(lastId, 10))
	stmt, err = connMain.Prepare("UPDATE Inspiration SET encId = ? WHERE id = ? LIMIT 1")
	res, err = stmt.Exec(inspirationEncId, lastId)

	return err
}

func GetInspirationsForTagging() ([]Inspiration, error) {
	conn := lib.GetDBConnection()
	inspirations := make([]Inspiration, INITIAL_INSPIRATION_COUNT)

	rows, err := conn.Query("SELECT "+INSPIRATION_FETCH_FIELDS+" FROM Inspiration WHERE tags = '' ORDER BY id DESC LIMIT ?", INITIAL_INSPIRATION_COUNT)
	if err != nil {
		return inspirations, err
	}

	err = scanAllInspirationFields(rows, inspirations)
	return inspirations, err
}

func TagInspiration(encId, tags, tagIds string) error {
	var tagsAdded string
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT IGNORE INTO InspirationTag(tagId, inspirationEncId) VALUES (?,?)")
	tagList := strings.Split(tags, ",")
	tagIdList := strings.Split(tagIds, ",")
	tagIdListLen := len(tagIdList)
	for i := 0; i < tagIdListLen; i++ {
		if tagIdList[i] == "" {
			continue
		}
		tagId, err := strconv.Atoi(tagIdList[i])
		if err != nil {
			log.Println(err)
			continue
		}
		tag := SanitizeForTagging(tagList[i])
		if tag == "" {
			log.Println("Empty tag during tagging")
			continue
		}
		res, err := stmt.Exec(tagId, encId)
		affectedRows, err := res.RowsAffected()
		if err != nil {
			log.Println(err)
			return err
		}
		if affectedRows == 1 {
			tagsAdded += tag + ","
		}
	}
	if tagsAdded != "" {
		stmt, err = conn.Prepare("UPDATE Inspiration SET tags = CONCAT(tags, ?) WHERE encId = ? LIMIT 1")
		_, err = stmt.Exec(tagsAdded, encId)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return err
}

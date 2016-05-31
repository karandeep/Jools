package controller

import (
	"../lib"
	"../model"
	"encoding/json"
	"github.com/gorilla/mux"
	"html"
	"net/http"
	"strconv"
)

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	if action == "add" {
		AddComment(w, r)
	} else if action == "get" {
		GetComments(w, r)
	}
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	var response JsonResponse
	session, err := store.Get(r, "user-session")
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	if session.Values["user"] == nil {
		w.Write(GetErrorJson("User is not logged in. Session may have expired."))
		return
	}

	subjectTypeStr := r.FormValue("type")
	encSubjectId := r.FormValue("subjectId")
	commentStr := r.FormValue("comment")

	if subjectTypeStr == "" || encSubjectId == "" || commentStr == "" {
		w.Write(GetErrorJson("Invalid params when adding comment"))
		return
	}
	subjectType, err := strconv.Atoi(subjectTypeStr)
	if err != nil {
		w.Write(GetErrorJson("Invalid subject type passed:" + subjectTypeStr))
		return
	}

	decryptedSubjectId, err := lib.Decrypt(encSubjectId)
	if err != nil {
		w.Write(GetErrorJson("Invalid subject id passed" + encSubjectId))
		return
	}
	subjectId, err := strconv.Atoi(decryptedSubjectId)
	if err != nil {
		w.Write(GetErrorJson("Invalid subject id passed" + encSubjectId))
		return
	}

	userData := GetUserDataFromSession(session)
	var comment model.Comment
	comment.UserName = userData.Name
	comment.UserEncId = userData.EncId
	comment.SubjectId = subjectId
	comment.SubjectType = subjectType
	comment.Comment = html.EscapeString(commentStr)
	comment.Created = lib.GetCurrentTimestamp()

	err = model.AddComment(comment)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	response.Success = true
	response.Data = map[string]string{
		"comment":     comment.Comment,
		"commentator": userData.Name,
	}
	data := GetJson(response)
	w.Write(data)

	if subjectType == model.COMMENT_TYPE_INSPIRATION {
		go model.UpdateInspirationLastComment(comment)
	}
}

func GetComments(w http.ResponseWriter, r *http.Request) {
	subjectTypeStr := r.FormValue("type")
	encSubjectId := r.FormValue("subjectId")

	if subjectTypeStr == "" || encSubjectId == "" {
		w.Write(GetErrorJson("Invalid type or subject id passed"))
		return
	}

	subjectType, err := strconv.Atoi(subjectTypeStr)
	if err != nil {
		w.Write(GetErrorJson("Invalid subject type passed:" + subjectTypeStr))
		return
	}

	decSubjectId, err := lib.Decrypt(encSubjectId)
	if err != nil {
		w.Write(GetErrorJson("Invalid subject id passed:" + encSubjectId))
		return
	}
	subjectId, err := strconv.Atoi(decSubjectId)
	if err != nil {
		w.Write(GetErrorJson("Invalid subject id passed" + encSubjectId))
		return
	}

	comments, err := model.GetComments(subjectType, subjectId)
	if err != nil {
		w.Write(GetErrorJson("Error fetching comments:" + err.Error()))
		return
	}
	commentsData, err := json.Marshal(comments)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	var response JsonResponse
	response.Success = true
	response.Data = map[string]string{
		"comments": string(commentsData),
	}
	data := GetJson(response)
	w.Write(data)
}

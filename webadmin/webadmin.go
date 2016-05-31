package webadmin

import (
	"../controller"
	"../model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func WebAdminHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	if action == "home" {
		Home(w, r)
	} else if action == "approveUploads" {
		ApproveUploads(w, r)
	} else if action == "markApproved" {
		MarkApproved(w, r)
	} else if action == "markRejected" {
		MarkRejected(w, r)
	} else if action == "getTaggingList" {
		GetTaggingList(w, r)
	} else if action == "tagInspiration" {
		TagInspiration(w, r)
	} else if action == "addTag" {
		AddTag(w, r)
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	response := controller.Response{}
	view := "webadmin/view/home.html"
	pageInfo := controller.PageInfo{
		Title: "Admin Home",
	}
	page := controller.GetAdminPage(view, response, pageInfo, r)
	w.Write(page)
}

func ApproveUploads(w http.ResponseWriter, r *http.Request) {
	uploads, err := model.GetUploadsForApproval()
	if err != nil {
		log.Println("Admin tool:", err)
		Home(w, r)
		return
	}
	uploadsJson, err := json.Marshal(uploads)
	if err != nil {
		log.Println("Admin tool:", err)
		Home(w, r)
		return
	}

	response := controller.Response{
		Data: map[string]string{
			"inspirations": string(uploadsJson),
		},
	}
	view := "webadmin/view/approveUploads.html"
	pageInfo := controller.PageInfo{
		Title: "Approve Uploads",
	}
	page := controller.GetAdminPage(view, response, pageInfo, r)
	w.Write(page)
}

func MarkApproved(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}
	model.MarkUpload(id, model.INSPIRATION_APPROVED)
}

func MarkRejected(w http.ResponseWriter, r *http.Request) {
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return
	}
	model.MarkUpload(id, model.INSPIRATION_REJECTED)
}

func GetTaggingList(w http.ResponseWriter, r *http.Request) {
	inspirations, err := model.GetInspirationsForTagging()
	if err != nil {
		log.Println("Admin tool:", err)
		Home(w, r)
		return
	}
	inspirationsJson, err := json.Marshal(inspirations)
	if err != nil {
		log.Println("Admin tool:", err)
		Home(w, r)
		return
	}

	allTags, err := model.GetAllTags()
	if err != nil {
		log.Println("Admin tool:", err)
		Home(w, r)
		return
	}
	allTagsJson, err := json.Marshal(allTags)
	if err != nil {
		log.Println("Admin tool:", err)
		Home(w, r)
		return
	}

	response := controller.Response{
		Data: map[string]string{
			"inspirations": string(inspirationsJson),
			"allTags":      string(allTagsJson),
		},
	}
	view := "webadmin/view/tagInspirations.html"
	pageInfo := controller.PageInfo{
		Title: "Tag Inspirations",
	}
	page := controller.GetAdminPage(view, response, pageInfo, r)
	w.Write(page)
}

func AddTag(w http.ResponseWriter, r *http.Request) {
	category := r.FormValue("category")
	tag := r.FormValue("tag")
	model.AddTag(category, tag)
	GetTaggingList(w, r)
}

func TagInspiration(w http.ResponseWriter, r *http.Request) {
	encId := r.FormValue("encId")
	tags := r.FormValue("tags")
	tagIds := r.FormValue("tagIds")
	if encId == "" || tagIds == "" || tags == "" {
		return
	}
	model.TagInspiration(encId, tags, tagIds)
}

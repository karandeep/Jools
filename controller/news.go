package controller

import (
	"../model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func NewsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	if action == "home" {
		Home(w, r)
	} else if action == "get" {
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "home"
	//isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	isMobileBrowser := true
	view := "view/mobile/newsHome.html"
	referrer := r.FormValue("jref")
	session, _ := store.Get(r, "user-session")
	pageInfo := PageInfo{
		Title:       "News",
		Description: "News",
		Keywords:    "News",
	}

	languages, err := model.GetLanguages()
	if err != nil {
		log.Println(err)
		return
	}
	papers, err := model.GetPapers()
	if err != nil {
		log.Println(err)
		return
	}
	news, err := model.GetNews()
	if err != nil {
		log.Println(err)
		return
	}

	languagesJson, err := json.Marshal(languages)
	if err != nil {
		log.Println(err)
		return
	}
	papersJson, err := json.Marshal(papers)
	if err != nil {
		log.Println(err)
		return
	}
	newsJson, err := json.Marshal(news)
	if err != nil {
		log.Println(err)
		return
	}
	response := Response{
		Data: map[string]string{
			"referrer":     referrer,
			"trackingName": pageTrackingName,
			"languages":    string(languagesJson),
			"papers":       string(papersJson),
			"news":         string(newsJson),
		},
	}

	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

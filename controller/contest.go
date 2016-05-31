package controller

import (
	"../lib"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ContestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	if action == "winPendant" {
		WinPendant(w, r)
	}
}

func WinPendant(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	pageTrackingName := "win_pendant"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/winPendant.html"
	if isMobileBrowser {
		view = "view/mobile/winPendant.html"
	}
	pageInfo := PageInfo{
		Title:       "Guess the price and win a Gold & Diamond Infiniti pendant",
		Description: "Jools.in - Guess the right price and win a Gold & Diamond Infiniti pendant, contest open till 31st March, 2014.",
		Keywords:    "Jewellery contest, Diamond pendant, Infiniti pendant, Guess the price, free jewellery, jewellery promotions",
		OG_TITLE:    "Win a Gold & Diamond pendant quiz @ Jools",
		OG_TYPE:     "my_jools:quiz",
	}

	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

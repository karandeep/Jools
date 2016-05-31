package controller

import (
	"../lib"
	"log"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	if !isMobileBrowser {
		//Shop(w, r)
		ListProducts(w, r)
		return
	}
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	view := "view/mobile/home.html"
	pageTrackingName := "home"
	pageInfo := PageInfo{
		Title:       "Online jewellery shopping store India | Jools.in",
		Description: "Online jewellery shopping store India | Jools.in",
	}
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

package controller

import (
	"../lib"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	if action == "aboutUs" {
		AboutUs(w, r)
	} else if action == "privacy" {
		PrivacyPolicy(w, r)
	} else if action == "shipping" {
		Shipping(w, r)
	} else if action == "tou" {
		TermsOfUse(w, r)
	}
}

func AboutUs(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	pageTrackingName := "aboutUs"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/aboutUs.html"
	if isMobileBrowser {
		//view = "view/mobile/privacyPolicy.html"
	}
	pageInfo := PageInfo{
		Title: "About Us",
	}

	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func PrivacyPolicy(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	pageTrackingName := "privacy"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/privacyPolicy.html"
	if isMobileBrowser {
		//view = "view/mobile/privacyPolicy.html"
	}
	pageInfo := PageInfo{
		Title: "Privacy Policy",
	}
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}

	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func Shipping(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	pageTrackingName := "shipping"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/shipping.html"
	if isMobileBrowser {
		//view = "view/mobile/privacyPolicy.html"
	}
	pageInfo := PageInfo{
		Title: "Shipping and Returns",
	}
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}

	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func TermsOfUse(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	pageTrackingName := "tou"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/termsOfUse.html"
	if isMobileBrowser {
		//view = "view/mobile/termsOfUse.html"
	}
	pageInfo := PageInfo{
		Title: "Terms of use",
	}
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}

	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

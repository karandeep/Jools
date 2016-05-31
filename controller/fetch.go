package controller

import (
	"../model"
	"errors"
	"github.com/gorilla/mux"
	"html"
	"net/http"
)

func FetchHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	if session.Values["user"] == nil {
		w.Write(GetErrorJson("User is not logged in. Session may have expired."))
		return
	}

	userData := GetUserDataFromSession(session)
	vars := mux.Vars(r)
	item := vars["item"]

	index := html.EscapeString(r.FormValue("index"))
	if index == "" {
		w.Write(GetErrorJson("Invalid params passed when fetching"))
		return
	}

	var syncedData string
	if item == "favorite" {
		syncedData, err = FetchFavorite(userData.Id)
	} else if item == "cart" {
		syncedData, err = FetchCart(userData.Id)
	} else {
		err = errors.New("No match found for fetch action")
	}

	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}

	var response JsonResponse
	response.Success = true
	response.Data = map[string]string{
		"index":      index,
		"syncedData": syncedData,
	}
	data := GetJson(response)
	w.Write(data)
}

func FetchFavorite(userId int64) (string, error) {
	return model.FetchFavoriteInspirations(userId)
}

func FetchCart(userId int64) (string, error) {
	return model.FetchCart(userId)
}

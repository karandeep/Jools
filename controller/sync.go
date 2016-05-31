package controller

import (
	"../lib"
	"../model"
	"errors"
	"github.com/gorilla/mux"
	"html"
	"net/http"
)

func SyncHandler(w http.ResponseWriter, r *http.Request) {
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
	localData := r.FormValue("localData")
	removedData := r.FormValue("removedData")
	if index == "" || (localData == "" && removedData == "") {
		w.Write(GetErrorJson("Invalid params passed when syncing"))
		return
	}

	var syncedData string
	if item == "view" {
		syncedData, err = SyncView(localData)
	} else if item == "favorite" {
		syncedData, err = SyncFavorite(userData.Id, localData, removedData)
	} else if item == "experiment" {
		syncedData, err = SyncExperiment(userData.Id, localData)
	} else if item == "cart" {
		syncedData, err = SyncCart(userData.Id, localData, removedData)
	} else {
		err = errors.New("No match found for sync action")
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

func SyncView(localData string) (string, error) {
	encIds := lib.ReformatEncIdString(localData)
	return model.IncrementInspirationViews(encIds)
}

func SyncFavorite(userId int64, localData string, removedData string) (string, error) {
	return model.SyncFavoriteInspirations(userId, localData, removedData)
}

func SyncExperiment(userId int64, localData string) (string, error) {
	return "", nil
}

func SyncCart(userId int64, localData string, removedData string) (string, error) {
	return model.SyncCart(userId, localData, removedData)
}

package controller

import (
	"../lib"
	"../model"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func RatingHandler(w http.ResponseWriter, r *http.Request) {
	var response JsonResponse
	response.Success = false

	encWinner := r.FormValue("winner")
	encLoser := r.FormValue("loser")
	if encWinner == "" || encLoser == "" {
		response.Data = map[string]string{"Message": "Empty winner/loser id passed"}
		data := GetJson(response)
		w.Write(data)
		return
	}

	decWinner, err := lib.Decrypt(encWinner)
	winner, err := strconv.Atoi(decWinner)
	if err != nil {
		response.Data = map[string]string{"Message": "Invalid winner id passed"}
		data := GetJson(response)
		w.Write(data)
		return
	}
	decLoser, err := lib.Decrypt(encLoser)
	loser, err := strconv.Atoi(decLoser)
	if err != nil {
		response.Data = map[string]string{"Message": "Invalid loser id passed"}
		data := GetJson(response)
		w.Write(data)
		return
	}

	err = model.Rate(winner, loser)
	if err != nil {
		log.Println("Rating error:", err)
	}

	randomImages, err := model.GetRandomImages()
	if err != nil {
		response.Data = map[string]string{"Message": err.Error()}
		data := GetJson(response)
		w.Write(data)
		return
	}

	response.Success = true
	hotImagesJson, err := json.Marshal(randomImages)
	response.Data = map[string]string{
		"HotImages": string(hotImagesJson),
	}
	data := GetJson(response)
	w.Write(data)
}

func RatingSkipHandler(w http.ResponseWriter, r *http.Request) {
	var response JsonResponse
	response.Success = false

	randomImages, err := model.GetRandomImages()
	if err != nil {
		response.Data = map[string]string{"Message": err.Error()}
		data := GetJson(response)
		w.Write(data)
		return
	}

	response.Success = true
	hotImagesJson, err := json.Marshal(randomImages)
	response.Data = map[string]string{
		"HotImages": string(hotImagesJson),
	}
	data := GetJson(response)
	w.Write(data)
}

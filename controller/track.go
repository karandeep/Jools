package controller

import (
	"../config"
	"../lib"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func TrackingHandler(w http.ResponseWriter, r *http.Request) {
	var trackingData lib.TrackData

	userId := r.FormValue("userId")
	if userId == "" {
		log.Println("User id cannot be empty for tracking")
		return
	}
	var err error
	if userId != "0" && strings.Contains(userId, "rand_") == false {
		userId, err = lib.Decrypt(userId)
		if err != nil {
			log.Println("Error in tracking:", err)
			return
		}
	}
	configData := config.GetConfig()
	if configData.ENV != config.PRODUCTION {
		return
	}
	trackingData.UserId = userId
	trackingData.Kingdom = r.FormValue("kingdom")
	trackingData.Phylum = r.FormValue("phylum")
	trackingData.Class = r.FormValue("class")
	trackingData.Family = r.FormValue("family")
	trackingData.Genus = r.FormValue("genus")
	trackingData.Species = r.FormValue("species")
	trackingData.Order = r.FormValue("order")
	trackingData.OS = r.FormValue("os")
	trackingData.Browser = r.FormValue("browser")
	trackingData.Tester = r.FormValue("tester")
	trackingData.Width = r.FormValue("width")
	trackingData.Height = r.FormValue("height")
	trackingData.IP = r.Header["X-Real-Ip"][0]
	trackingData.FwdIP = r.Header["Cf-Connecting-Ip"][0]
	incrementBy := r.FormValue("incrementBy")

	if trackingData.Kingdom == "" || trackingData.Phylum == "" {
		log.Println("Atleast kingdom and phylum needed for tracking")
		return
	}
	incrementCount, err := strconv.Atoi(incrementBy)
	if err != nil {
		log.Println("Invalid value passed for increment (tracking)", r)
		return
	}

	trackingData.Date = lib.GetTrackingDate()
	go lib.TrackCounter(trackingData, incrementCount)

	if trackingData.Kingdom == "debug" && trackingData.Phylum == "temp_id" {
		go lib.UpdateId(trackingData.Class, trackingData.UserId, trackingData.Date)
	}

	if trackingData.Kingdom == "experiment" {
		var expData lib.ExperimentData
		expData.UserId = trackingData.UserId
		expData.Experiment = trackingData.Phylum
		expData.Variant = trackingData.Class
		expData.Created = lib.GetCurrentTimestamp()
		go lib.TrackExperiment(expData)
	}
	return
}

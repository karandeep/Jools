package controller

import (
	"../config"
	"../lib"
	"../model"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func InspirationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	if page == "hotTrends" {
		HotTrends(w, r)
	} else if page == "getInitialInfo" {
		GetInitialInfoForInspiration(w, r)
	} else if page == "getNextLot" {
		GetNextInspirationLot(w, r)
	} else if page == "playHot" {
		PlayHotOrNot(w, r)
	} else if page == "upload" {
		UploadInspirations(w, r)
	}
}

func HotTrends(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	inspirations, err := model.GetInitialInspirationList(model.RANDOM_INSPIRATIONS)
	if err != nil {
		log.Println("Error showing inspirations:", err)
		return
	}

	inspirationJson, err := json.Marshal(inspirations)
	if err != nil {
		log.Println("Unable to convert inspirations to json")
		return
	}
	configData := config.GetConfig()
	title := "Jewellery inspirations from across the world, Get featured by uploading design images"
	description := "Jools.in - Check out exclusive jewellery designs online in India with latest designs - Soon you will be able to select an image and buy a similar design."
	keywords := "Jewellery inspirations, jewellery designs, occasion designs, fashion jewellery in India, browse jewellery trends,play hot or not, earn credits"
	canonical := configData.BASE_URL + r.URL.Path
	og_title := title
	og_description := description
	og_url := canonical
	og_image := ""
	og_type := "product"
	encId := r.FormValue("id")
	if encId != "" {
		decId, err := lib.Decrypt(encId)
		if err == nil {
			og_url += "?id=" + encId
			title = "Browse Trends | Jools's jewellery lovers community"
			og_title = "Online jewellery shopping, Love, Customize, Share, Shop @ Jools!"
			description = "Join Jools's jewellery lovers community and share style ideas, inspirations, and photos with members across the globe"
			decIdInt, err := strconv.Atoi(decId)
			if err == nil {
				inspiration, err := model.GetInspiration(decIdInt)
				if err == nil {
					og_image = configData.INSPIRATIONS_URL + "/" + inspiration.ImageName
				}
			}
		}
	}
	pageTrackingName := "hot_trends"
	response := Response{
		Data: map[string]string{
			"inspirations":        string(inspirationJson),
			"showInspirationInfo": encId,
			"trackingName":        pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/hotTrends.html"
	if isMobileBrowser {
		view = "view/mobile/hotTrends.html"
	}
	pageInfo := PageInfo{
		Title:       title,
		Description: description,
		Keywords:    keywords,
		Canonical:   canonical,
		OG_TITLE:    og_title,
		OG_DESC:     og_description,
		OG_TYPE:     og_type,
		OG_URL:      og_url,
		OG_IMAGE:    og_image,
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func PlayHotOrNot(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	hotImages, err := model.GetRandomImages()
	var hotImagesJson, topRatedImagesJson []byte
	if err == nil {
		hotImagesJson, err = json.Marshal(hotImages)
	}
	topRatedImages, err := model.GetTopRated(5)
	if err == nil {
		topRatedImagesJson, err = json.Marshal(topRatedImages)
	}
	pageTrackingName := "play_hot"
	response := Response{
		Data: map[string]string{
			"hotImages":    string(hotImagesJson),
			"topRated":     string(topRatedImagesJson),
			"trackingName": pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/playHot.html"
	if isMobileBrowser {
		view = "view/mobile/playHot.html"
	}
	configData := config.GetConfig()
	canonical := configData.BASE_URL + r.URL.Path
	og_title := "HOT or NOT design @ Jools"
	og_type := "my_jools:photo"
	og_url := canonical
	cogs := r.FormValue("cogs")
	if cogs == "quiz" {
		og_title = "Be the Stylist quiz @ Jools"
		og_type = "my_jools:hot_or_not"
		og_url += "?cogs=quiz"
	}
	pageInfo := PageInfo{
		Title:       "Play Hot or Not, you get to pick and upload great jewellery designs",
		Description: "Jools.in - Play Hot or NOT, you get to pick and upload great jewellery designs. Soon you will be able to buy all the top rated designs at our website.",
		Keywords:    "Jewellery inspirations, play hot or not, earn credits,jewellery websites, diamond jewellery india, gold jewellery online, Indian diamond jewellery, indian jewellery",
		OG_TITLE:    og_title,
		OG_TYPE:     og_type,
		OG_URL:      og_url,
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func UploadInspirations(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	if session.Values["user"] == nil {
		return
	}
	var response Response
	view := "view/popup/uploadInspirations.html"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	if isMobileBrowser {
		view = "view/mobile/uploadInspirations.html"
		pageTrackingName := "uploadInspirations"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Upload your designs",
			Description: "Upload your designs",
		}
		response = Response{
			Data: map[string]string{
				"trackingName": pageTrackingName,
			},
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func GetNextInspirationLot(w http.ResponseWriter, r *http.Request) {
	var response JsonResponse
	var data []byte

	encLastId := r.FormValue("lastId")
	if encLastId == "" {
		w.Write(GetErrorJson("Empty last id passed for inspiration lot"))
		return
	}

	decId, err := lib.Decrypt(encLastId)
	if err != nil {
		w.Write(GetErrorJson("Invalid value for last id"))
		return
	}
	lastId, err := strconv.Atoi(decId)
	if err != nil {
		w.Write(GetErrorJson("Invalid value for last id"))
		return
	}

	nextLot, err := model.GetNextInspirationLot(lastId)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	nextInspirations, err := json.Marshal(nextLot)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	response.Success = true
	response.Data = map[string]string{
		"inspirations": string(nextInspirations),
	}
	data = GetJson(response)
	w.Write(data)
}

func GetInitialInfoForInspiration(w http.ResponseWriter, r *http.Request) {
	encSubjectId := r.FormValue("subjectId")
	decSubjectId, err := lib.Decrypt(encSubjectId)
	if err != nil {
		log.Println("Invalid value for subject id", encSubjectId)
		return
	}
	subjectId, err := strconv.Atoi(decSubjectId)
	if err != nil {
		log.Println("Invalid value for subject id", encSubjectId)
		return
	}
	inspiration, err := model.GetInspiration(subjectId)
	if err != nil {
		log.Println("Error when loading inspiration info", err)
		return
	}

	inspirationJson, err := json.Marshal(inspiration)
	if err != nil {
		log.Println("Unable to convert inspiration data to json", err)
		return
	}
	daysTillOpen := lib.DaysLeftTill(SHOP_OPENS_AT)
	daysTillOpenStr := strconv.Itoa(daysTillOpen)
	var response Response
	response.Data = map[string]string{
		"inspiration":  string(inspirationJson),
		"daysTillOpen": daysTillOpenStr,
	}
	view := "view/popup/inspirationInfo.html"
	data := GetPopup(view, response)
	w.Write(data)
}

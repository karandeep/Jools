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

const SHOP_OPENS_AT int32 = 1396224000

func ProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	action := vars["action"]
	if action == "shop" {
		//Shop(w, r)
		ListProducts(w, r)
	} else if action == "realShop" {
		RealShop(w, r)
	} else if action == "list" {
		ListProducts(w, r)
	} else if action == "getNextList" {
		GetNextProductList(w, r)
	}
}

func ListProducts(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")

	pageTrackingName := "productList"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/productList.html"
	if isMobileBrowser {
		view = "view/mobile/productList.html"
	}
	title := "Product List"
	description := "Product List"
	keywords := "Product list"

	var refinement model.ProductRefinement
	categoryStr := r.FormValue("category")
	if categoryStr != "" {
		refinement.Category, _ = strconv.Atoi(categoryStr)
	}
	gemstoneStr := r.FormValue("gemstone")
	refinement.Gemstone = -1
	if gemstoneStr != "" {
		refinement.Gemstone, _ = strconv.Atoi(gemstoneStr)
	}
	settingStr := r.FormValue("setting")
	refinement.Setting = -1
	if settingStr != "" {
		refinement.Setting, _ = strconv.Atoi(settingStr)
	}
	priMetalStr := r.FormValue("priMetal")
	refinement.PrimaryMetal = -1
	if priMetalStr != "" {
		refinement.PrimaryMetal, _ = strconv.Atoi(priMetalStr)
	}
	priMetalColorStr := r.FormValue("priMetalColor")
	refinement.PrimaryMetalColor = -1
	if priMetalColorStr != "" {
		refinement.PrimaryMetalColor, _ = strconv.Atoi(priMetalColorStr)
	}
	priceMinStr := r.FormValue("priceMin")
	refinement.PriceMin = -1
	if priceMinStr != "" {
		refinement.PriceMin, _ = strconv.Atoi(priceMinStr)
	}
	priceMaxStr := r.FormValue("priceMax")
	refinement.PriceMax = -1
	if priceMaxStr != "" {
		refinement.PriceMax, _ = strconv.Atoi(priceMaxStr)
	}
	sortOrderStr := r.FormValue("sort")
	refinement.Sort = 0
	if sortOrderStr != "" {
		refinement.Sort, _ = strconv.Atoi(sortOrderStr)
	}

	products, err := model.GetInitialProductList(refinement)
	if err != nil {
		log.Println(err)
		HomeHandler(w, r)
		return
	}
	productListJson, err := json.Marshal(products)
	if err != nil {
		log.Println(err)
		HomeHandler(w, r)
		return
	}
	refinementJson, err := json.Marshal(refinement)
	response := Response{
		Data: map[string]string{
			"productList":       string(productListJson),
			"productRefinement": string(refinementJson),
			"trackingName":      pageTrackingName,
		},
	}
	pageInfo := PageInfo{
		Title:       title,
		Description: description,
		Keywords:    keywords,
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func GetNextProductList(w http.ResponseWriter, r *http.Request) {
	var response JsonResponse
	var data []byte

	var refinement model.ProductRefinement
	categoryStr := r.FormValue("category")
	if categoryStr != "" {
		refinement.Category, _ = strconv.Atoi(categoryStr)
	}
	gemstoneStr := r.FormValue("gemstone")
	refinement.Gemstone = -1
	if gemstoneStr != "" {
		refinement.Gemstone, _ = strconv.Atoi(gemstoneStr)
	}
	settingStr := r.FormValue("setting")
	refinement.Setting = -1
	if settingStr != "" {
		refinement.Setting, _ = strconv.Atoi(settingStr)
	}
	priMetalStr := r.FormValue("priMetal")
	refinement.PrimaryMetal = -1
	if priMetalStr != "" {
		refinement.PrimaryMetal, _ = strconv.Atoi(priMetalStr)
	}
	priMetalColorStr := r.FormValue("priMetalColor")
	refinement.PrimaryMetalColor = -1
	if priMetalColorStr != "" {
		refinement.PrimaryMetalColor, _ = strconv.Atoi(priMetalColorStr)
	}
	priceMinStr := r.FormValue("priceMin")
	refinement.PriceMin = -1
	if priceMinStr != "" {
		refinement.PriceMin, _ = strconv.Atoi(priceMinStr)
	}
	priceMaxStr := r.FormValue("priceMax")
	refinement.PriceMax = -1
	if priceMaxStr != "" {
		refinement.PriceMax, _ = strconv.Atoi(priceMaxStr)
	}
	sortOrderStr := r.FormValue("sort")
	refinement.Sort = 0
	if sortOrderStr != "" {
		refinement.Sort, _ = strconv.Atoi(sortOrderStr)
	}

	nextList, err := model.GetInitialProductList(refinement)

	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	nextListJson, err := json.Marshal(nextList)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	refinementJson, err := json.Marshal(refinement)
	response.Success = true
	response.Data = map[string]string{
		"productList":       string(nextListJson),
		"productRefinement": string(refinementJson),
	}
	data = GetJson(response)
	w.Write(data)
}

func RealShop(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	configData := config.GetConfig()
	//TODO: Remove default productId
	productId := 22
	var product model.Product
	productCodeStr := r.FormValue("productCode")
	canonical := configData.BASE_URL + r.URL.Path
	if productCodeStr != "" {
		product, err = model.GetProductFromIndex(productCodeStr)
		if err != nil {
			log.Println(err)
			HomeHandler(w, r)
			return
		}
	} else {
		product, err = model.GetProduct(productId)
		if err != nil {
			log.Println(err)
			HomeHandler(w, r)
			return
		}
	}

	variants, err := product.GetAllCombosForProduct()
	if err != nil {
		log.Println(err)
		HomeHandler(w, r)
		return
	}

	title := product.CreateTitle()
	product.Description = product.CreateDescription()
	description := product.Description + " | " + product.Name
	keywords := product.CreateKeywords()
	imageIndex, uniqueIndex := product.GetIndices()
	canonical += "?productCode=" + uniqueIndex
	og_image := product.CreateImageUrl(imageIndex)
	og_type := "ornament"

	productInfo := model.ProductInfo{
		ProductData:     product,
		ProductVariants: variants,
	}
	productInfoJson, err := json.Marshal(productInfo)
	if err != nil {
		log.Println(err)
		HomeHandler(w, r)
		return
	}
	FBCommentsUrl := canonical
	pageTrackingName := "shop"
	response := Response{
		Data: map[string]string{
			"productInfo":   string(productInfoJson),
			"fbCommentsUrl": FBCommentsUrl,
			"trackingName":  pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/shop.html"
	if isMobileBrowser {
		view = "view/mobile/shop.html"
	}
	pageInfo := PageInfo{
		Title:       title,
		Description: description,
		Keywords:    keywords,
		Canonical:   canonical,
		OG_IMAGE:    og_image,
		OG_TYPE:     og_type,
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func Shop(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	daysTillOpen := lib.DaysLeftTill(SHOP_OPENS_AT)
	daysTillOpenStr := strconv.Itoa(daysTillOpen)
	pageTrackingName := "shop"
	response := Response{
		Data: map[string]string{
			"daysTillOpen": daysTillOpenStr,
			"trackingName": pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/shopPreview.html"
	if isMobileBrowser {
		view = "view/mobile/shopPreview.html"
	}
	configData := config.GetConfig()
	cogs := r.FormValue("cogs")
	canonical := configData.BASE_URL + r.URL.Path
	og_url := canonical
	og_type := "website"
	if cogs == "cash" {
		og_url = canonical + "?cogs=cash"
		og_type = "my_jools:five_hundred_cash"
	}
	pageInfo := PageInfo{
		Title:       "Buy Jewellery Online in India with Exclusive and latest designs for 2014",
		Description: "Jools.in - Buy Gold and Diamond Jewellery Online in India with the latest jewellery designs 2014 from our online jewelry shopping store with COD, 30-day free returns on jewellery, free shipping and a lifetime exchange policy. Best jewellery website for online jewellery shopping in India with Designer Rings, Pendants, Earrings, Mangalsutra, Bangles, Bracelets, Solitaire Diamonds.",
		Keywords:    "Jewellery, buy Jewellery, buy Jewellery  in India, earn credits, custom jewellery,diamond jewellery, gold jewellery, jewellery website, jewellery designs, fashion jewellery, indian jewellery, designer jewellery, diamond Jewellery,  fashion Jewellery, online jewellery shopping, online jewellery shopping india, jewellery websites, diamond jewellery india, gold jewellery online, Indian diamond jewellery",
		OG_TITLE:    "Rs 500 cash giftcard @ Jools!",
		OG_DESC:     "I just got a Rs 500 giftcard on joining Jools. Get yours now @ Jools.in",
		OG_URL:      og_url,
		OG_TYPE:     og_type,
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

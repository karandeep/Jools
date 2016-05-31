package controller

import (
	"../config"
	"../lib"
	"../model"
	"bytes"
	"encoding/json"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PageInfo struct {
	Title       string
	Description string
	Keywords    string
	Canonical   string
	OG_TITLE    string
	OG_DESC     string
	OG_TYPE     string
	OG_URL      string
	OG_IMAGE    string
}

// Contents are of type HTML to prevent html escaping.
type Content struct {
	Config      config.ConfigData
	PageInfo    PageInfo
	InternalCSS template.HTML
	InternalJS  template.HTML
	HeaderHTML  template.HTML
	ContentHTML template.HTML
	FooterHTML  template.HTML
	FooterJS    template.HTML
}

type HeaderInfo struct {
	Config   config.ConfigData
	UserName string
}

type FooterInfo struct {
	Config       config.ConfigData
	TrackingName string
}

type FooterJSInfo struct {
	AccountJustCreated   bool
	Config               config.ConfigData
	ExperimentInfo       string
	IncludeRecaptcha     bool
	IsUserLoggedIn       bool
	Referrer             string
	SignupCompletionReqd bool
	ShowInitialCongrats  bool
	ShowChangePassword   bool
	ShowFBConnect        bool
	UserData             string
	TrackingName         string
	Source               string
	Medium               string
	Content              string
	Campaign             string
	TopRated             string
	HotImages            string
	Inspirations         string
	ShowInviteFriends    string
	AllTags              string
}

type Response struct {
	Config config.ConfigData
	Data   map[string]string
}

type JsonResponse struct {
	Success bool
	Data    map[string]string
}

func parseTemplate(file string, data interface{}) (out []byte, error error) {
	var buf bytes.Buffer
	t, err := template.ParseFiles(file)
	if err != nil {
		return nil, err
	}
	err = t.Execute(&buf, data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GetPage(session *sessions.Session, file string, isMobileBrowser bool, data Response, pageInfo PageInfo, r *http.Request) []byte {
	configData := config.GetConfig()
	internalCSSFile := "view/generated/internal_css.html"
	if isMobileBrowser {
		internalCSSFile = "view/generated/internal_css_mobile.html"
	}
	internalCSS, error := parseTemplate(internalCSSFile, map[string]string{})
	if error != nil {
		print(error.Error())
	}

	internalJS, error := parseTemplate("view/generated/internal_js.html", map[string]string{})
	if error != nil {
		print(error.Error())
	}

	userLoggedIn := false
	signupCompletionReqd := false
	showInitialCongrats := false
	showChangePassword := false
	showFBConnect := false

	var publicUser model.PublicUser
	justCreated := model.AccountJustCreated
	if justCreated {
		model.AccountJustCreated = false
	}
	if session != nil && session.Values["user"] != nil {
		userLoggedIn = true
		userData := GetUserDataFromSession(session)
		publicUser = userData.GetPublicUser()
		if justCreated && userData.Fbid == "" && !isMobileBrowser {
			showFBConnect = true
		} else if userData.Password == "" && (userData.Status == model.VERIFIED || userData.Status == model.VERIFICATION_PENDING) {
			signupCompletionReqd = true
		} else if userData.Status == model.PASSWORD_UPDATE_PENDING {
			showChangePassword = true
		} else if justCreated {
			showInitialCongrats = true
		}
	}
	var header []byte
	headerFile := "view/header.html"
	if isMobileBrowser {
		headerFile = "view/mobile/header.html"
	}
	header, error = parseTemplate(headerFile,
		HeaderInfo{
			Config:   configData,
			UserName: publicUser.Name,
		})
	if error != nil {
		print(error.Error())
	}

	publicUserData, _ := json.Marshal(publicUser)

	var footer []byte
	footerFile := "view/footer.html"
	if isMobileBrowser {
		footerFile = "view/mobile/footer.html"
	}
	footer, error = parseTemplate(footerFile,
		FooterInfo{
			Config:       configData,
			TrackingName: data.Data["trackingName"],
		})
	if error != nil {
		print(error.Error())
	}

	var footerJS []byte
	footerJSFile := "view/footerJS.html"
	if isMobileBrowser {
		footerJSFile = "view/footerJS.html"
	}

	var utm_source, utm_medium, utm_content, utm_campaign string
	utm_source = r.FormValue("utm_source")
	if utm_source == "" {
		refererUrl := r.Referer()
		if refererUrl != "" {
			//Is this request from a non jools referer?
			if !strings.Contains(refererUrl, configData.DOMAIN) {
				utm_source = refererUrl
			}
		}
	} else {
		utm_medium = r.FormValue("utm_medium")
		utm_content = r.FormValue("utm_content")
		utm_campaign = r.FormValue("utm_campaign")
	}
	referrer := r.FormValue("jref")
	if referrer == "" {
		referrer = "-1"
	}
	//For requests coming in to FB App
	requestIds := r.FormValue("request_ids")
	if requestIds != "" {
		requestIdList := strings.Split(requestIds, ",")
		requestId := requestIdList[0]
		encUserId, err := model.GetUserWhoSentRequest(requestId)
		if err == nil && encUserId != "" {
			referrer = encUserId
		}
	}
	footerJS, error = parseTemplate(footerJSFile,
		FooterJSInfo{
			AccountJustCreated:   justCreated,
			Config:               configData,
			ExperimentInfo:       data.Data["experimentInfo"],
			IncludeRecaptcha:     false,
			IsUserLoggedIn:       userLoggedIn,
			Referrer:             referrer,
			SignupCompletionReqd: signupCompletionReqd,
			ShowInitialCongrats:  showInitialCongrats,
			ShowChangePassword:   showChangePassword,
			ShowFBConnect:        showFBConnect,
			UserData:             string(publicUserData),
			TrackingName:         data.Data["trackingName"],
			Source:               utm_source,
			Medium:               utm_medium,
			Content:              utm_content,
			Campaign:             utm_campaign,
			TopRated:             data.Data["topRated"],
			HotImages:            data.Data["hotImages"],
			Inspirations:         data.Data["inspirations"],
			ShowInviteFriends:    data.Data["showInviteFriends"],
		})
	if error != nil {
		print(error.Error())
	}
	data.Config = configData
	content, error := parseTemplate(file, data)
	if error != nil {
		print(error.Error())
	}

	baseView := "view/base.html"
	if isMobileBrowser {
		baseView = "view/mobile/base.html"
	} else {
		pageInfo.Title += " | Jools.in"
	}
	daysTillOpen := lib.DaysLeftTill(SHOP_OPENS_AT)
	daysTillOpenStr := strconv.Itoa(daysTillOpen)
	year, month, day := time.Now().Date()
	pageInfo.Description += "As on " + strconv.Itoa(day) + "-" + month.String() + "-" + strconv.Itoa(year) + ",only " + daysTillOpenStr + " days to launch."
	if pageInfo.Canonical == "" {
		pageInfo.Canonical = configData.BASE_URL + r.URL.Path // + "?" + r.URL.RawQuery ..Use this when you need to have params, like for specific product
	}
	if pageInfo.OG_TITLE == "" {
		pageInfo.OG_TITLE = pageInfo.Title
	}
	if pageInfo.OG_DESC == "" {
		pageInfo.OG_DESC = pageInfo.Description
	}
	if pageInfo.OG_TYPE == "" {
		pageInfo.OG_TYPE = "website"
	}
	if pageInfo.OG_URL == "" {
		pageInfo.OG_URL = pageInfo.Canonical
	}
	if pageInfo.OG_IMAGE == "" {
		pageInfo.OG_IMAGE = configData.STATIC_URL + "/images/guess_price.jpg"
	}
	base, error := parseTemplate(baseView,
		Content{
			Config:      configData,
			PageInfo:    pageInfo,
			InternalCSS: template.HTML(internalCSS),
			InternalJS:  template.HTML(internalJS),
			HeaderHTML:  template.HTML(header),
			FooterHTML:  template.HTML(footer),
			FooterJS:    template.HTML(footerJS),
			ContentHTML: template.HTML(content),
		})
	if error != nil {
		print(error.Error())
		return []byte("Internal server error")
	}
	return base
}

func GetJson(data JsonResponse) []byte {
	responseData, err := json.Marshal(data)
	if err != nil {
		return []byte("Internal server error")
	}
	return responseData
}

func GetErrorJson(message string) []byte {
	var response JsonResponse
	response.Success = false
	response.Data = map[string]string{
		"message": message,
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		return []byte("Internal server error")
	}
	return responseData
}

func GetPopup(file string, data Response) []byte {
	data.Config = config.GetConfig()
	content, error := parseTemplate(file, data)
	if error != nil {
		print(error.Error())
	}
	return content
}

func GetAdminPage(file string, data Response, pageInfo PageInfo, r *http.Request) []byte {
	configData := config.GetConfig()
	internalCSSFile := "view/generated/internal_css.html"
	internalCSS, error := parseTemplate(internalCSSFile, map[string]string{})
	if error != nil {
		print(error.Error())
	}

	internalJS, error := parseTemplate("view/generated/internal_js.html", map[string]string{})
	if error != nil {
		print(error.Error())
	}

	var header []byte
	headerFile := "webadmin/view/header.html"
	header, error = parseTemplate(headerFile,
		HeaderInfo{
			Config: configData,
		})
	if error != nil {
		print(error.Error())
	}

	var footer []byte
	footerFile := "webadmin/view/footer.html"
	footer, error = parseTemplate(footerFile,
		FooterInfo{
			Config: configData,
		})
	if error != nil {
		print(error.Error())
	}

	var footerJS []byte
	footerJSFile := "webadmin/view/footerJS.html"
	footerJS, error = parseTemplate(footerJSFile,
		FooterJSInfo{
			Config:       configData,
			Inspirations: data.Data["inspirations"],
			AllTags:      data.Data["allTags"],
		})
	if error != nil {
		print(error.Error())
	}
	data.Config = configData
	content, error := parseTemplate(file, data)
	if error != nil {
		print(error.Error())
	}

	baseView := "view/base.html"
	base, error := parseTemplate(baseView,
		Content{
			Config:      configData,
			PageInfo:    pageInfo,
			InternalCSS: template.HTML(internalCSS),
			InternalJS:  template.HTML(internalJS),
			HeaderHTML:  template.HTML(header),
			FooterHTML:  template.HTML(footer),
			FooterJS:    template.HTML(footerJS),
			ContentHTML: template.HTML(content),
		})
	if error != nil {
		print(error.Error())
		return []byte("Internal server error")
	}
	return base
}

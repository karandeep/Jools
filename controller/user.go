package controller

import (
	"../config"
	"../lib"
	"../model"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/huandu/facebook"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var store = sessions.NewFilesystemStore("/tmp", []byte("J00lsU$erSess!on$"))

// Struct for holding the parts that merges multiple templates.
func UserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	if page == "account" {
		Account(w, r)
	} else if page == "askEmail" {
		AskEmail(w, r)
	} else if page == "askFBConnect" {
		AskFBConnect(w, r)
	} else if page == "authenticate" {
		Authenticate(w, r)
	} else if page == "changePassword" {
		ChangePassword(w, r)
	} else if page == "completeSignup" {
		CompleteSignup(w, r)
	} else if page == "checkout" {
		Checkout(w, r)
	} else if page == "create" {
		CreateAccount(w, r)
	} else if page == "earnCash" {
		EarnCash(w, r)
	} else if page == "forgotPassword" {
		ForgotPassword(w, r)
	} else if page == "generateAndFetchHash" {
		GenerateAndFetchHash(w, r)
	} else if page == "inviteFBFriends" {
		InviteFBFriends(w, r)
	} else if page == "invitesSent" {
		InvitesSent(w, r)
	} else if page == "login" {
		Login(w, r)
	} else if page == "logout" {
		Logout(w, r)
	} else if page == "markForPasswordUpdate" {
		MarkForPasswordUpdate(w, r)
	} else if page == "purchase" {
		Purchase(w, r)
	} else if page == "purchaseSuccess" {
		PurchaseSuccess(w, r)
	} else if page == "removeAddress" {
		RemoveAddress(w, r)
	} else if page == "sendReferralEmail" {
		SendReferralEmail(w, r)
	} else if page == "showInitialCongrats" {
		ShowInitialCongrats(w, r)
	} else if page == "signInWithFB" {
		SignInWithFB(w, r)
	} else if page == "signup" {
		Signup(w, r)
	} else if page == "signupCompletionReqd" {
		SignupCompletionReqd(w, r)
	} else if page == "showChangePassword" {
		ShowChangePassword(w, r)
	} else if page == "showInviteFriends" {
		InviteFriends(w, r)
	} else if page == "storeFBRequest" {
		StoreFBRequest(w, r)
	} else if page == "verify" {
		VerifyEmail(w, r)
	} else if page == "viewCart" {
		ViewCart(w, r)
	} else {
		log.Println("Unknown action passed to user controller")
	}
	return
}

func Account(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	userData := GetUserDataFromSession(session)
	if userData.Id == 0 {
		HomeHandler(w, r)
		return
	}

	var addresses [model.MAX_ADDRESSES_PER_USER]model.Address
	var addressesStr string
	addresses, err = model.GetAllAddressesForUser(userData.Email, userData.Id)
	if err == nil {
		addressesJson, _ := json.Marshal(addresses)
		addressesStr = string(addressesJson)
	}

	var orders [model.MAX_ORDERS_TO_SHOW]model.Order
	var ordersStr string
	orders, err = model.GetOrdersForUser(userData.Email, userData.Id)
	if err == nil {
		ordersJson, _ := json.Marshal(orders)
		ordersStr = string(ordersJson)
	}

	pageTrackingName := "account"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
			"addresses":    addressesStr,
			"orders":       ordersStr,
		},
	}

	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/account.html"
	if isMobileBrowser {
		view = "view/mobile/account.html"
	}
	pageInfo := PageInfo{
		Title:       "User account",
		Description: "User account",
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func PopulateSessionInfo(userData model.User, session *sessions.Session, w http.ResponseWriter, r *http.Request) {
	session.Values["user"] = userData
	session.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 1800,
	}
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
}

func GetUserDataFromSession(session *sessions.Session) model.User {
	//In some cases the session data is of type *model.User, and in some cases it is model.User
	//This probably happens because GetRegistry avoids duplicate storage for type interface{}, instead using pointers in some cases
	var userData model.User
	var userDataPtr *model.User
	userData, ok := session.Values["user"].(model.User)
	if !ok {
		userDataPtr, ok = session.Values["user"].(*model.User)
		if ok {
			userData = *userDataPtr
		}
	}
	return userData
}

func SignInWithFB(w http.ResponseWriter, r *http.Request) {
	encReferrer := r.FormValue("jref")
	_ = r.FormValue("userID")
	_ = r.FormValue("accessToken")
	signedRequest := r.FormValue("signedRequest")
	configData := config.GetConfig()

	session, err := store.Get(r, "user-session")
	var response JsonResponse
	if err != nil || signedRequest == "" {
		w.Write(GetErrorJson("Invalid params when doing fb sign in"))
		return
	}

	userData := GetUserDataFromSession(session)
	if userData.Fbid == "" {
		globalApp := facebook.New(configData.FB_APP_ID, configData.FB_APP_SECRET)
		fbSession, err := globalApp.SessionFromSignedRequest(signedRequest)
		if err != nil {
			log.Println("Error during FB login:", err)
			w.Write(GetErrorJson("Error occurred during fb sign in"))
			return
		}
		//fbSession := globalApp.Session(globalApp.AppAccessToken())
		res, err := fbSession.Api("/me", "GET", facebook.Params{
			"fields": "id,name,email,gender",
		})
		if err != nil {
			log.Println("Error during FB login:", err)
			w.Write(GetErrorJson("Error occurred during fb sign in"))
			return
		}
		var userEmail string
		res.DecodeField("email", &userEmail)
		if userEmail == "" || err != nil {
			log.Println("Error during FB login:", err)
			w.Write(GetErrorJson("Error occurred during fb sign in"))
			return
		}

		var userData model.User
		userData = model.GetUserByEmail(userEmail)
		if userData.Id != 0 {
			if userData.Fbid == "" {
				res.DecodeField("name", &userData.Name)
				res.DecodeField("id", &userData.Fbid)
				res.DecodeField("gender", &userData.Gender)
				model.ConnectAccountToFB(userEmail, userData.Name, userData.Fbid, userData.Gender)
				go model.StoreFBFriends(fbSession, userData.Fbid)
			} else {
				//Just go ahead and sign in
			}
		} else {
			res.DecodeField("id", &userData.Fbid)
			res.DecodeField("name", &userData.Name)
			res.DecodeField("gender", &userData.Gender)
			userData.Email = userEmail

			referrer := -1
			if encReferrer != "" && encReferrer != "-1" {
				decReferrer, err := lib.Decrypt(encReferrer)
				referrer, err = strconv.Atoi(decReferrer)
				if err != nil {
					referrer = -1
				}
			}
			userData.Referrer = referrer

			userData, err = model.CreateUserAccount(userData)
			if err != nil {
				log.Println("Account creation failure:", err)
				w.Write(GetErrorJson("Account creation failure"))
				return
			}
			go model.StoreFBFriends(fbSession, userData.Fbid)
		}
		userData.Status = model.CREATED_VIA_FB
		PopulateSessionInfo(userData, session, w, r)
	}
	response.Success = true
	data := GetJson(response)
	w.Write(data)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	encReferrer := r.FormValue("jref")
	loginEmail := r.FormValue("login_email")

	var response JsonResponse
	if loginEmail == "" || !lib.IsValidEmail(loginEmail) {
		log.Println("login email found to be invalid", r)
		w.Write(GetErrorJson("login email found to be invalid"))
		return
	}

	var userData model.User
	userData = model.GetUserByEmail(loginEmail)
	if userData.Id != 0 {
		log.Println("User with specified email id already exists", loginEmail)
		w.Write(GetErrorJson("User with specified email id already exists"))
		return
	}

	referrer := -1
	if encReferrer != "" && encReferrer != "-1" {
		decReferrer, err := lib.Decrypt(encReferrer)
		referrer, err = strconv.Atoi(decReferrer)
		if err != nil {
			referrer = -1
		}
	}
	userData.Referrer = referrer
	userData.Status = model.VERIFICATION_PENDING
	userData.Code = lib.GetRandomInt()
	userData.Email = loginEmail
	userData, err := model.CreateUserAccount(userData)
	if err != nil {
		log.Println("Account creation failure:", err)
		w.Write(GetErrorJson("Account creation failure. Try again."))
		return
	}
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println("Unable to create session for user during account creation")
		w.Write(GetErrorJson("Unable to create session."))
		return
	}
	PopulateSessionInfo(userData, session, w, r)
	response.Success = true
	data := GetJson(response)
	w.Write(data)
}

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	code := r.FormValue("code")
	hash := r.FormValue("hash")
	verified, userData, err := model.VerifyUser(email, code, hash)

	if !verified {
		log.Println("Verification failed:", err)
		HomeHandler(w, r)
		return
	}

	session, _ := store.Get(r, "user-session")
	PopulateSessionInfo(userData, session, w, r)
	HomeHandler(w, r)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "user-session")
	session.Options = &sessions.Options{MaxAge: -1, Path: "/"}
	session.Values["user"] = nil
	session.Save(r, w)
	HomeHandler(w, r)
}

func CompleteSignup(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("username")
	password := r.FormValue("password")
	verifyPassword := r.FormValue("verifyPassword")

	var response JsonResponse
	session, err := store.Get(r, "user-session")
	userData := GetUserDataFromSession(session)
	if err != nil || session.Values["user"] == nil ||
		name == "" || password == "" || password != verifyPassword ||
		userData.Name != "" {
		w.Write(GetErrorJson("Invalid params when completing signup"))
		return
	}

	userData.Name = html.EscapeString(name)
	userData.Password = lib.GetPasswordHash(password)

	err = userData.CompleteSignup()
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	PopulateSessionInfo(userData, session, w, r)
	response.Success = true
	data := GetJson(response)
	w.Write(data)
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	verifyPassword := r.FormValue("verifyPassword")

	var response JsonResponse
	session, err := store.Get(r, "user-session")
	userData := GetUserDataFromSession(session)
	if err != nil || session.Values["user"] == nil ||
		password == "" || password != verifyPassword {
		w.Write(GetErrorJson("Entered passwords don't match."))
		return
	}

	userData.Password = lib.GetPasswordHash(password)
	userData.Status = model.VERIFIED

	err = userData.ChangePassword()
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	PopulateSessionInfo(userData, session, w, r)
	response.Success = true
	response.Data = map[string]string{"message": "Your password was successfully changed."}
	data := GetJson(response)
	w.Write(data)
}

func Authenticate(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("login_email")
	password := r.FormValue("login_password")

	var response JsonResponse
	var message string
	var data []byte
	session, err := store.Get(r, "user-session")
	userData, authenticated := model.AuthenticateUser(email, password)
	if !authenticated || err != nil {
		if err != nil {
			message = "Temporarily unable to login. Please try again."
		} else {
			message = "Invalid email or password"
		}
		response.Success = false
		response.Data = map[string]string{
			"message": message,
		}
		data = GetJson(response)
		w.Write(data)
		return
	}

	PopulateSessionInfo(userData, session, w, r)
	response.Success = true
	data = GetJson(response)
	w.Write(data)
}

func SendReferralEmail(w http.ResponseWriter, r *http.Request) {
	emails := r.FormValue("emails")
	session, err := store.Get(r, "user-session")

	var response JsonResponse
	var message string
	var data []byte
	if err != nil || emails == "" || session.Values["user"] == nil {
		response.Success = false
		response.Data = map[string]string{
			"message": "Temporarily unable to send emails. Please try again.",
		}
		data = GetJson(response)
		w.Write(data)
		return
	}

	emailList := strings.Split(emails, ",")
	var validEmails []string
	sentEmails := 0
	for _, email := range emailList {
		if lib.IsValidEmail(email) {
			sentEmails++
			validEmails = append(validEmails, email)
		}
	}

	userData := GetUserDataFromSession(session)
	senderName := userData.Name
	if senderName == "" {
		senderName = userData.Email
	}
	inviteMessage := model.Message{
		Type:         model.INVITE,
		Sender:       userData.Email,
		SenderName:   senderName,
		ReferralLink: model.GetReferralLink(userData.EncId),
		Fbid:         userData.Fbid,
		Recepients:   validEmails,
	}
	inviteMessage.Send()

	if sentEmails == 1 {
		message = "1 email successfully sent"
	} else {
		message = strconv.Itoa(sentEmails) + " emails successfully sent"
	}
	response.Success = true
	response.Data = map[string]string{
		"message": message,
	}
	data = GetJson(response)
	w.Write(data)
}

func ShowInitialCongrats(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "initialCongrats"
	response := Response{
		Data: map[string]string{
			"credit":       "500",
			"trackingName": pageTrackingName,
		},
	}

	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/initialCongrats.html"
	if isMobileBrowser {
		view = "view/mobile/initialCongrats.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Welcome to Jools",
			Description: "Welcome to Jools",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func SignupCompletionReqd(w http.ResponseWriter, r *http.Request) {
	var response Response
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/completeSignup.html"
	if isMobileBrowser {
		view = "view/mobile/completeSignup.html"
		pageTrackingName := "completeSignup"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Welcome to Jools",
			Description: "Welcome to Jools",
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

func Login(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "login"
	referrer := r.FormValue("jref")
	response := Response{
		Data: map[string]string{
			"referrer":     referrer,
			"trackingName": pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/login.html"
	if isMobileBrowser {
		view = "view/mobile/login.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Sign In to Jools",
			Description: "Jools.in - Sign in to avail great offers.",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "signup"
	referrer := r.FormValue("jref")
	response := Response{
		Data: map[string]string{
			"referrer":     referrer,
			"trackingName": pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/signup.html"
	if isMobileBrowser {
		view = "view/mobile/signup.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Sign up for Jools",
			Description: "Jools.in - Sign up to avail great offers.",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "forgotPass"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/forgotPassword.html"
	referrer := r.FormValue("jref")
	response := Response{
		Data: map[string]string{
			"referrer":     referrer,
			"trackingName": pageTrackingName,
		},
	}
	if isMobileBrowser {
		view = "view/mobile/forgotPassword.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Forgot your password?",
			Description: "Forgot your password?",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func ShowChangePassword(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "changePass"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}

	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/changePassword.html"
	if isMobileBrowser {
		view = "view/mobile/changePassword.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Change your password",
			Description: "Change your password",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func AskEmail(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "askEmail"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/askEmail.html"
	referrer := r.FormValue("jref")
	response := Response{
		Data: map[string]string{
			"referrer":     referrer,
			"trackingName": pageTrackingName,
		},
	}
	if isMobileBrowser {
		view = "view/mobile/askEmail.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Create your account",
			Description: "Create your account",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func AskFBConnect(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "fbconnect"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/popup/askFBConnect.html"
	referrer := r.FormValue("jref")
	response := Response{
		Data: map[string]string{
			"referrer":     referrer,
			"trackingName": pageTrackingName,
		},
	}
	if isMobileBrowser {
		view = "view/mobile/askFBConnect.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Connect your account to facebook",
			Description: "Connect your account to facebook",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func MarkForPasswordUpdate(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("login_email")
	var response JsonResponse

	if email == "" || !lib.IsValidEmail(email) {
		response.Success = false
		response.Data = map[string]string{"message": "Please enter valid email address."}
		data := GetJson(response)
		w.Write(data)
		return
	}

	err := model.MarkForPasswordUpdate(email)
	if err != nil {
		response.Success = false
		response.Data = map[string]string{"message": err.Error()}
		data := GetJson(response)
		w.Write(data)
		return
	}

	response.Success = true
	response.Data = map[string]string{"message": "You've been mailed the instructions for resetting your password. Please use the link in the email to update your password."}
	data := GetJson(response)
	w.Write(data)
}

func GenerateAndFetchHash(w http.ResponseWriter, r *http.Request) {
	var response JsonResponse

	session, err := store.Get(r, "user-session")
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	if session.Values["user"] == nil {
		w.Write(GetErrorJson("User is not logged in"))
		return
	}
	userData := GetUserDataFromSession(session)
	curTime := strconv.Itoa(int(lib.GetCurrentTimestamp()))
	hash := lib.GetHashOf(userData.EncId + curTime)
	response.Success = true
	response.Data = map[string]string{
		"userId": userData.EncId,
		"time":   curTime,
		"hash":   hash,
	}
	data := GetJson(response)
	w.Write(data)
}

func StoreFBRequest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	reqId := r.FormValue("reqId")
	recepients := r.Form["recepients[]"]

	session, err := store.Get(r, "user-session")
	if err != nil {
		return
	}
	if session.Values["user"] == nil {
		return
	}
	userData := GetUserDataFromSession(session)
	model.StoreFBRequest(reqId, userData.EncId)
	model.StoreRequestFbids(recepients)
}

func InvitesSent(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "invitesSent"
	inviteCount := r.FormValue("count")
	response := Response{
		Data: map[string]string{
			"inviteCount":  inviteCount,
			"trackingName": pageTrackingName,
		},
	}

	view := "view/popup/invitesSent.html"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	if isMobileBrowser {
		view = "view/mobile/invitesSent.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Invites successfully sent",
			Description: "Invites successfully sent",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func InviteFriends(w http.ResponseWriter, r *http.Request) {
	networkStr := r.FormValue("network")
	network, err := strconv.Atoi(networkStr)
	if err != nil {
		log.Println(err)
		return
	}
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	if session.Values["user"] == nil {
		log.Println("No user session")
		return
	}
	userData := GetUserDataFromSession(session)
	emailData, err := model.GetEmailData(network, userData.EncId)
	if err != nil {
		log.Println(err)
		return
	}
	emailDataJson, err := json.Marshal(emailData)
	if err != nil {
		log.Println(err)
		return
	}
	pageTrackingName := "invitesFriends"
	response := Response{
		Data: map[string]string{
			"network":      networkStr,
			"emailData":    string(emailDataJson),
			"trackingName": pageTrackingName,
		},
	}

	view := "view/popup/inviteFriends.html"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	if isMobileBrowser {
		view = "view/mobile/inviteFriends.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Invite your friends to Jools",
			Description: "Invite your friends to Jools",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func InviteFBFriends(w http.ResponseWriter, r *http.Request) {
	pageTrackingName := "inviteFBFriends"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}

	view := "view/popup/inviteFBFriends.html"
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	if isMobileBrowser {
		view = "view/mobile/inviteFBFriends.html"
		session, _ := store.Get(r, "user-session")
		pageInfo := PageInfo{
			Title:       "Invite your friends to Jools",
			Description: "Invite your friends to Jools",
		}
		page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
		w.Write(page)
	} else {
		data := GetPopup(view, response)
		w.Write(data)
	}
}

func EarnCash(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	userData := GetUserDataFromSession(session)
	configData := config.GetConfig()
	referrer := r.FormValue("jref")
	importedContacts := r.FormValue("importedContacts")
	var showInviteFriends string
	if importedContacts != "" {
		network, err := strconv.Atoi(importedContacts)
		if err == nil {
			if (network == configData.GOOGLE && userData.GmailImport == 0) || (network == configData.YAHOO && userData.YahooImport == 0) {
				//Repopulate session info to reflect correct status
				userData = model.GetUserByEmail(userData.Email)
				PopulateSessionInfo(userData, session, w, r)
			}
			if (network == configData.GOOGLE && userData.GmailImport == 1) || (network == configData.YAHOO && userData.YahooImport == 1) {
				showInviteFriends = importedContacts
			}
		}
	}
	referralLink := configData.REFERRAL_BASE + "?jref=" + userData.EncId
	pageTrackingName := "earn_cash"
	response := Response{
		Data: map[string]string{
			"referrer":          referrer,
			"referralLink":      referralLink,
			"showInviteFriends": showInviteFriends,
			"trackingName":      pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/earnCash.html"
	if isMobileBrowser {
		view = "view/mobile/earnCash.html"
	}
	pageInfo := PageInfo{
		Title:       "Earn cash credits to buy jewellery",
		Description: "Jools.in - Earn cash credits by inviting friends. You earn 500 and they earn 500.",
		Keywords:    "Invite friends, Earn credits, rupees, earn indian rupees, jewellery promotion",
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func ViewCart(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}
	pageTrackingName := "viewCart"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/cart.html"
	if isMobileBrowser {
		view = "view/mobile/cart.html"
	}
	pageInfo := PageInfo{
		Title:       "Your shopping cart",
		Description: "Your shopping cart",
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func Checkout(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	userData := GetUserDataFromSession(session)
	var addresses [model.MAX_ADDRESSES_PER_USER]model.Address
	var addressesStr string
	if userData.Id != 0 {
		addresses, err = model.GetAddressesForUser(userData.Email, userData.Id)
		if err == nil {
			addressesJson, _ := json.Marshal(addresses)
			addressesStr = string(addressesJson)
		} else {
			log.Println(err)
		}
	}
	pageTrackingName := "checkout"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
			"addresses":    addressesStr,
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/checkout.html"
	if isMobileBrowser {
		view = "view/mobile/checkout.html"
	}
	pageInfo := PageInfo{
		Title:       "Checkout",
		Description: "Checkout",
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func Purchase(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	userData := GetUserDataFromSession(session)

	email := r.FormValue("email")
	name := r.FormValue("name")
	address := r.FormValue("address")
	city := r.FormValue("city")
	state := r.FormValue("state")
	pincode := r.FormValue("pincode")
	mobile := r.FormValue("mobile")
	products := r.FormValue("products")
	paymentMethod := r.FormValue("paymentMethod")
	if email == "" || name == "" || address == "" || city == "" || state == "" || pincode == "" ||
		mobile == "" || products == "" || paymentMethod == "" {
		w.Write(GetErrorJson("Invalid params passed for purchase"))
		return
	}

	var order model.Order
	var addressInfo model.Address
	var addressId int64
	addressIdStr := r.FormValue("addressId")
	if addressIdStr != "" && addressIdStr != "0" {
		addressId, err = strconv.ParseInt(addressIdStr, 10, 64)
		if err == nil {
			_, err = model.GetAddressFromId(addressId)
			if err != nil {
				w.Write(GetErrorJson("Invalid addressId passed for purchase"))
				return
			}
		}
	} else {
		addressInfo.UserId = userData.Id
		addressInfo.Email = email
		addressInfo.Name = name
		addressInfo.Address = address
		addressInfo.City = city
		addressInfo.State = state
		addressInfo.Pincode, err = strconv.Atoi(pincode)
		if err != nil {
			w.Write(GetErrorJson("Invalid pincode passed for purchase"))
			return
		}
		addressInfo.Mobile = mobile
		addressId, err = addressInfo.Store()
		if err != nil {
			w.Write(GetErrorJson(err.Error()))
			return
		}
	}
	order.UserId = userData.Id
	order.Email = email
	order.AddressId = addressId
	order.Products = products //This is the json data for cart, as visible to the user at this point
	order.PaymentMethod, err = strconv.Atoi(paymentMethod)
	if err != nil {
		w.Write(GetErrorJson("Invalid payment method passed for purchase"))
		return
	}
	order.Cost, err = model.GetCostOfCartProducts(products)
	if err != nil {
		log.Println(err)
		w.Write(GetErrorJson("Invalid products in cart"))
		return
	}

	order, err = order.Begin()
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	orderJson, err := json.Marshal(order)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}
	addressJson, _ := json.Marshal(address)

	var response JsonResponse
	response.Success = true
	response.Data = map[string]string{
		"order":   string(orderJson),
		"address": string(addressJson),
	}
	data := GetJson(response)
	w.Write(data)
}

func PurchaseSuccess(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		log.Println(err)
		return
	}

	orderIdStr := r.FormValue("orderId")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		log.Println("Invalid order id passed")
		HomeHandler(w, r)
		return
	}
	order, address, err := model.GetOrder(orderId)
	if err != nil {
		log.Println("Order not found for order id:", orderId, err)
		HomeHandler(w, r)
		return
	}

	orderJson, _ := json.Marshal(order)
	addressJson, _ := json.Marshal(address)
	pageTrackingName := "purchaseSuccess"
	response := Response{
		Data: map[string]string{
			"trackingName": pageTrackingName,
			"order":        string(orderJson),
			"address":      string(addressJson),
		},
	}
	isMobileBrowser := lib.IsMobileBrowser(r.UserAgent())
	view := "view/purchaseSuccess.html"
	if isMobileBrowser {
		view = "view/mobile/purchaseSuccess.html"
	}
	pageInfo := PageInfo{
		Title:       "Congratulations! Your order has been placed.",
		Description: "Congratulations! Your order has been placed.",
	}
	page := GetPage(session, view, isMobileBrowser, response, pageInfo, r)
	w.Write(page)
}

func RemoveAddress(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "user-session")
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}

	userData := GetUserDataFromSession(session)
	if userData.Id == 0 {
		w.Write(GetErrorJson("User is not logged in"))
		return
	}
	addressIdStr := r.FormValue("addressId")
	addressId, err := strconv.ParseInt(addressIdStr, 10, 64)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}

	err = model.RemoveAddressForUser(addressId, userData.Email)
	if err != nil {
		w.Write(GetErrorJson(err.Error()))
		return
	}

	var response JsonResponse
	response.Success = true
	data := GetJson(response)
	w.Write(data)
}

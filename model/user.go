package model

import (
	"../config"
	"../lib"
	"bytes"
	"errors"
	"github.com/huandu/facebook"
	"html/template"
	"log"
	"strconv"
	textTemplate "text/template"
)

type User struct {
	Id            int64
	Name          string
	Email         string
	Password      string
	Gender        string
	Fbid          string
	EncId         string
	ReferralCount int
	Referrer      int
	Status        int
	Code          int64
	Created       int32
	GmailImport   int
	YahooImport   int
	EarnedCredits int
	BoughtCredits int
}

type PublicUser struct {
	Name          string
	Email         string
	Gender        string
	Fbid          string
	EncId         string
	ReferralCount int
	GmailImport   int
	YahooImport   int
	EarnedCredits int
	BoughtCredits int
}

type FBFriendData struct {
	Id string
}

var AccountJustCreated bool = false

const (
	VERIFICATION_PENDING    int = 0
	VERIFIED                int = 1
	PASSWORD_UPDATE_PENDING int = 2
	CREATED_VIA_FB          int = 3
)

func (userData User) GetPublicUser() PublicUser {
	var publicUser PublicUser

	publicUser.Name = userData.Name
	publicUser.Email = userData.Email
	publicUser.Gender = userData.Gender
	publicUser.Fbid = userData.Fbid
	publicUser.EncId = userData.EncId
	publicUser.ReferralCount = userData.ReferralCount
	publicUser.GmailImport = userData.GmailImport
	publicUser.YahooImport = userData.YahooImport
	publicUser.EarnedCredits = userData.EarnedCredits
	publicUser.BoughtCredits = userData.BoughtCredits

	return publicUser
}

func GetUserByEncId(encId string) User {
	conn := lib.GetDBConnection()
	var userData User

	row := conn.QueryRow("SELECT id,name,email,password,gender,fbid,encId,referralCount,referrer,status,code,created,gmailImport,yahooImport,earnedCredits,boughtCredits FROM User WHERE encId = ?", encId)

	_ = row.Scan(
		&userData.Id,
		&userData.Name,
		&userData.Email,
		&userData.Password,
		&userData.Gender,
		&userData.Fbid,
		&userData.EncId,
		&userData.ReferralCount,
		&userData.Referrer,
		&userData.Status,
		&userData.Code,
		&userData.Created,
		&userData.GmailImport,
		&userData.YahooImport,
		&userData.EarnedCredits,
		&userData.BoughtCredits,
	)
	return userData
}

func GetUserByEmail(email string) User {
	conn := lib.GetDBConnection()
	var userData User

	row := conn.QueryRow("SELECT id,name,email,password,gender,fbid,encId,referralCount,referrer,status,code,created,gmailImport,yahooImport,earnedCredits,boughtCredits FROM User WHERE email = ?", email)

	_ = row.Scan(
		&userData.Id,
		&userData.Name,
		&userData.Email,
		&userData.Password,
		&userData.Gender,
		&userData.Fbid,
		&userData.EncId,
		&userData.ReferralCount,
		&userData.Referrer,
		&userData.Status,
		&userData.Code,
		&userData.Created,
		&userData.GmailImport,
		&userData.YahooImport,
		&userData.EarnedCredits,
		&userData.BoughtCredits,
	)
	return userData
}

func CreateUserAccount(userData User) (User, error) {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT INTO User (name, email, password, gender, fbid, encId, referralCount, referrer, status, code, created, earnedCredits) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return userData, err
	}

	userData.Created = lib.GetCurrentTimestamp()
	//Not expecting the password to have a value at this point
	res, err := stmt.Exec(userData.Name, userData.Email, "", userData.Gender, userData.Fbid, userData.EncId, userData.ReferralCount, userData.Referrer, userData.Status, userData.Code, userData.Created, userData.EarnedCredits)
	if err != nil {
		return userData, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return userData, err
	}
	userData.Id = lastId

	encId, _ := lib.Encrypt(strconv.FormatInt(lastId, 10))
	stmt, err = conn.Prepare("UPDATE User SET encId = ? WHERE id = ? LIMIT 1")
	res, err = stmt.Exec(encId, lastId)
	/*rowCnt, err := res.RowsAffected()*/
	if err != nil {
		return userData, err
	}
	userData.EncId = encId

	if userData.Referrer != -1 && userData.Status == CREATED_VIA_FB {
		incrementReferralCount(userData.Referrer)
	}

	if userData.Status == VERIFICATION_PENDING {
		configData := config.GetConfig()
		codeString := strconv.FormatInt(userData.Code, 10)
		verificationHash := GetVerificationHash(userData.Email, codeString)
		if configData.ENV != config.DEV {
			sendVerificationEmail(verificationHash, userData.Email, codeString, encId)
		} else {
			log.Println("Verification Pending: Code - ", userData.Code, "; Hash - ", verificationHash)
		}
	}
	AccountJustCreated = true
	return userData, nil
}

func ConnectAccountToFB(email, name, fbid, gender string) {
	conn := lib.GetDBConnection()
	stmt, _ := conn.Prepare("UPDATE User SET name = ?, fbid = ?, gender = ?, status = ? WHERE email = ? LIMIT 1")
	_, _ = stmt.Exec(name, fbid, gender, CREATED_VIA_FB, email)
}

func incrementReferralCount(referrer int) {
	//TODO: Also increment earned credits by the applicable amount
	conn := lib.GetDBConnection()
	stmt, _ := conn.Prepare("UPDATE User SET referralCount = referralCount + 1 WHERE id = ?")
	_, _ = stmt.Exec(referrer)
}

func GetVerificationHash(email string, code string) string {
	configData := config.GetConfig()
	hashedString := lib.GetHashOf(configData.SALT + email + code + configData.SALT)
	return hashedString
}

func sendVerificationEmail(verificationHash string, email string, code string, encId string) {
	configData := config.GetConfig()
	verificationLink := configData.BASE_URL + "/user/verify?email=" + email + "&code=" + code + "&hash=" + verificationHash
	referralLink := GetReferralLink(encId)
	to := email
	subject := "Welcome to Jools, confirm your account"

	var buf bytes.Buffer
	var bufText bytes.Buffer
	data := map[string]string{"VerificationLink": verificationLink, "ReferralLink": referralLink}
	t, _ := template.ParseFiles("email/welcome.html")
	_ = t.Execute(&buf, data)
	htmlMessage := string(buf.Bytes())

	txt, _ := textTemplate.ParseFiles("email/welcome.txt")
	_ = txt.Execute(&bufText, data)
	txtMessage := string(bufText.Bytes())
	lib.SendEmail(htmlMessage, txtMessage, to, subject)
}

func sendPasswordChangeEmail(verificationHash string, email string, code string) {
	configData := config.GetConfig()
	changePasswordLink := configData.BASE_URL + "/user/verify?email=" + email + "&code=" + code + "&hash=" + verificationHash
	to := email
	subject := "Forgot your password on Jools?"

	var buf bytes.Buffer
	var bufText bytes.Buffer
	data := map[string]string{"ChangePasswordLink": changePasswordLink}
	t, _ := template.ParseFiles("email/forgot_password.html")
	_ = t.Execute(&buf, data)
	htmlMessage := string(buf.Bytes())

	txt, _ := textTemplate.ParseFiles("email/forgot_password.txt")
	_ = txt.Execute(&bufText, data)
	txtMessage := string(bufText.Bytes())
	lib.SendEmail(htmlMessage, txtMessage, to, subject)
}

func GetReferralLink(encId string) string {
	configData := config.GetConfig()
	return configData.REFERRAL_BASE + "?jref=" + encId
}

func VerifyUser(email string, code string, hash string) (bool, User, error) {
	var userData User
	if email == "" || !lib.IsValidEmail(email) || code == "" || hash == "" {
		return false, userData, errors.New("Invalid params passed")
	}

	userData = GetUserByEmail(email)

	if userData.Id == 0 {
		return false, userData, errors.New("No user account with that email address")
	}

	if userData.Status != VERIFICATION_PENDING && userData.Status != PASSWORD_UPDATE_PENDING {
		return false, userData, errors.New("User account is not in a state to be verified")
	}

	codeString := strconv.FormatInt(userData.Code, 10)
	if codeString != code {
		return false, userData, errors.New("Invalid verification code")
	}

	computedHash := GetVerificationHash(email, code)
	if computedHash != hash {
		return false, userData, errors.New("Hash verification failed")
	}

	if userData.Status == PASSWORD_UPDATE_PENDING {
		return true, userData, nil
	}

	conn := lib.GetDBConnection()
	stmt, _ := conn.Prepare("UPDATE User SET status = ?, code = 0 WHERE id = ?")
	_, err := stmt.Exec(VERIFIED, userData.Id)
	if err != nil {
		return false, userData, err
	}
	userData.Status = VERIFIED
	if userData.Referrer != -1 {
		incrementReferralCount(userData.Referrer)
	}
	return true, userData, nil
}

func (userData User) CompleteSignup() error {
	conn := lib.GetDBConnection()
	stmt, _ := conn.Prepare("UPDATE User SET name = ?, password = ? WHERE id = ?")
	_, err := stmt.Exec(userData.Name, userData.Password, userData.Id)
	if err != nil {
		return err
	}
	return nil
}

func (userData User) ChangePassword() error {
	conn := lib.GetDBConnection()
	stmt, _ := conn.Prepare("UPDATE User SET password = ?, status = ? WHERE id = ?")
	_, err := stmt.Exec(userData.Password, userData.Status, userData.Id)
	if err != nil {
		return err
	}
	return nil
}

func MarkForPasswordUpdate(email string) error {
	code := strconv.FormatInt(lib.GetRandomInt(), 10)
	conn := lib.GetDBConnection()
	verificationHash := GetVerificationHash(email, code)

	stmt, _ := conn.Prepare("UPDATE User SET status = ?, code = ? WHERE email = ? LIMIT 1")
	res, err := stmt.Exec(PASSWORD_UPDATE_PENDING, code, email)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowCnt == 0 {
		return errors.New("No user found for given email id")
	}

	configData := config.GetConfig()
	if configData.ENV != config.DEV {
		sendPasswordChangeEmail(verificationHash, email, code)
	} else {
		log.Println("Verification Pending: Code - ", code, "; Hash - ", verificationHash)
	}
	return nil
}

func (userData User) UpdatePassword() error {
	return nil
}

func AuthenticateUser(email string, password string) (User, bool) {
	userData := GetUserByEmail(email)
	if userData.Id == 0 {
		return userData, false
	}

	passwordHash := lib.GetPasswordHash(password)
	if passwordHash != userData.Password {
		return userData, false
	}
	return userData, true
}

func StoreFBRequest(reqId string, userEncId string) error {
	conn := lib.GetDBConnection()
	curTime := lib.GetCurrentTimestamp()
	stmt, err := conn.Prepare("INSERT INTO FBRequests (reqId, encUserId,created) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(reqId, userEncId, curTime)
	return err
}

func StoreFBFriends(fbSession *facebook.Session, userFbid string) {
	res, err := fbSession.Api("/me/friends", "GET", facebook.Params{
		"fields": "id",
	})
	if err != nil {
		log.Println("Error during FB friends fetch:", err)
		return
	}
	var friendData FBFriendData
	allFriends := make([]string, 5000)
	for i := 0; i < 5000; i++ {
		err = res.DecodeField("data."+strconv.Itoa(i), &friendData)
		if err != nil {
			break
		} else {
			allFriends[i] = friendData.Id
		}
	}
	//TODO: Make this into a single insert query
	//https://code.google.com/p/go/issues/detail?id=5171
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT IGNORE INTO FBFriends (userFbid,friendFbid) VALUES (?,?)")
	if err != nil {
		log.Println(err)
		return
	}
	for _, friendFbid := range allFriends {
		if friendFbid == "" {
			break
		}
		_, err = stmt.Exec(userFbid, friendFbid)
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func StoreRequestFbids(recepients []string) error {
	conn := lib.GetDBConnection()
	stmt, err := conn.Prepare("INSERT IGNORE INTO FBRequestsSentTo (fbid) VALUES (?)")
	for _, fbid := range recepients {
		_, err = stmt.Exec(fbid)
	}
	return err
}

func GetUserWhoSentRequest(reqId string) (string, error) {
	var encUserId string
	conn := lib.GetDBConnection()
	row := conn.QueryRow("SELECT encUserId FROM FBRequests WHERE reqId = ?", reqId)

	err := row.Scan(&encUserId)
	return encUserId, err
}

package lib

import (
	"../config"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net/smtp"
	"regexp"
	"strings"
	"time"
)

func Decrypt(text string) (string, error) {
	configData := config.GetConfig()
	key := []byte(configData.KEY)
	ciphertext, _ := base64.URLEncoding.DecodeString(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		log.Println("Cipher Text too short:", ciphertext, text)
		return "", errors.New("unknown key type")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return removePadding(string(ciphertext)), nil
}

func Encrypt(text string) (string, error) {
	configData := config.GetConfig()
	key := []byte(configData.KEY)
	paddedText := addPaddingForAES(text)
	plaintext := []byte(paddedText)
	if len(plaintext)%aes.BlockSize != 0 {
		return "", errors.New("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func addPaddingForAES(text string) string {
	currentLen := len(text)
	paddingNeeded := aes.BlockSize - (currentLen % aes.BlockSize)

	for i := 0; i < paddingNeeded; i++ {
		text += "?"
	}
	return text
}

func removePadding(text string) string {
	return strings.TrimRight(text, "?")
}

func GetRandomInt() int64 {
	maxRandom := big.NewInt(100000)
	randomInt, _ := rand.Int(rand.Reader, maxRandom)
	return randomInt.Int64()
}

func GetRandomIntInRange(min int, max int) int {
	diff := max - min
	if diff <= 0 {
		return max
	}
	randomInt := GetRandomInt()
	randomIntInRange := min + (int(randomInt) % diff)
	return randomIntInRange
}

func GetHashOf(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func GetSha256HashOf(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func GetPasswordHash(pwd string) string {
	configData := config.GetConfig()
	return GetHashOf(configData.PASSWORD_SALT + pwd + configData.PASSWORD_SALT)
}

func IsValidEmail(email string) bool {
	re := regexp.MustCompile("^[_a-z0-9-]+(.[_a-z0-9-]+)*@[a-z0-9-]+(.[a-z0-9-]+)*(.[a-z]{2,3})$")
	return re.MatchString(email)
}

func GetDKIMSignature(canonicalForm string) (string, error) {
	keyFile := "../lib/dkim_pvt_key.txt"

	// Read the private key
	pemData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Println("read key file:", err)
		return "", err
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Println("bad key data:", "not PEM-encoded")
		return "", errors.New("bad key data")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Println("unknown key type ", got, "want ", got, want)
		return "", errors.New("unknown key type")
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Println("bad private key: ", err)
		return "", err
	}

	var out []byte
	label := "dkim-sig"
	out, err = rsa.EncryptOAEP(sha1.New(), rand.Reader, &priv.PublicKey, []byte(canonicalForm), []byte(label))
	if err != nil {
		log.Println("encrypt: ", err)
		return "", err
	}

	hashedString := base64.URLEncoding.EncodeToString(out)
	return hashedString, nil
}

func GetDKIMHeader(msg []byte, to string, mime string, replyTo string, contentType string) (string, error) {
	bodyHash := GetSha256HashOf(msg)
	canonicalForm := "to:" + to + "\nmime-version:" + mime + "\nreply-to:" + replyTo +
		"\ncontent-type:" + contentType + "\n"
	signature, err := GetDKIMSignature(canonicalForm)
	if err != nil {
		return "", err
	}

	header := "dkim-signature:v=1;a=rsa-sha256;bh=" + bodyHash + ";c=relaxed;d=www.jools.in;" +
		"h=to:mime-version:reply-to:content-type;s=1384941434.jools;b=" + signature

	return header, nil
}

func SendEmail(htmlMessage string, txtMessage string, to string, subject string) {
	c, err := smtp.Dial("127.0.0.1:25")
	if err != nil {
		log.Println(err)
		return
	}
	// Set the sender and recipient.
	from := "Jools <team@jools.in>"
	c.Mail(from)
	c.Rcpt(to)
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Println(err)
		return
	}
	defer wc.Close()
	toHeader := "To: " + to + "\n"
	fromHeader := "From : " + from + "\n"
	replyTo := "Reply-to: " + from + "\n"
	subject = "Subject: " + subject + "\n"
	xPriority := "X-Priority: 3\n"
	mime := "MIME-version: 1.0;\n"
	msgBoundary := "b1_9fcc68b7415b7327f44d7c7c4a749855"
	contentType := "Content-Type: multipart/alternative;boundary=\"" + msgBoundary + "\"\n"
	txtMsgContentType := "Content-Type: text/plain; charset = \"UTF-8\"\n"
	htmlMsgContentType := "Content-Type: text/html; charset=\"UTF-8\"\n"
	msgEncoding := "Content-Transfer-Encoding: 8bit\n\n"
	msg := toHeader + fromHeader + replyTo + xPriority + subject + mime + contentType +
		"--" + msgBoundary + "\n" + txtMsgContentType + msgEncoding + txtMessage +
		"\n\n--" + msgBoundary + "\n" + htmlMsgContentType + msgEncoding + htmlMessage +
		"\n\n--" + msgBoundary + "--"

	buf := bytes.NewBufferString(msg)
	//dkim := GetDKIMHeader(buf, to, "1.0", from)
	if _, err = buf.WriteTo(wc); err != nil {
		log.Println(err)
		return
	}
}

func IsMobileBrowser(userAgent string) bool {
	//return true
	if strings.Contains(userAgent, "Mobile") {
		return true
	} else if strings.Contains(userAgent, "Phone") {
		return true
	} else if strings.Contains(userAgent, "Android") {
		return true
	}
	return false
}

func GetCurrentTimestamp() int32 {
	return int32(time.Now().Unix())
}

func DaysLeftTill(endTimestamp int32) int {
	curTimestamp := GetCurrentTimestamp()
	daysLeft := int((endTimestamp - curTimestamp) / 86400)
	return daysLeft
}

func ReformatEncIdString(encIds string) string {
	return "'" + strings.TrimRight(strings.Replace(encIds, ";", "','", -1), ",'") + "'"
}

func GetSanitizedName(curName string) string {
	lowerCase := strings.ToLower(curName)
	return strings.Replace(lowerCase, " ", "-", -1)
}

func GetTimeFromTimestamp(timestamp int32) time.Time {
	return time.Unix(int64(timestamp), 0)
}
func TimeAfterDuration(seconds int, startTime time.Time) (int, string, int) {
	duration := time.Duration(seconds) * time.Second
	year, month, day := startTime.Add(duration).Date()
	monthStr := month.String()
	return year, monthStr, day
}

func MonthAbbr(month string) string {
	return month[0:3]
}

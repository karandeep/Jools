package lib

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestEncryption(t *testing.T) {
	Convey("Encrypt a string", t, func() {
		var encryptedStr, decryptedStr string
		originalStr := "This string will be encrypted"
		encryptedStr, err := Encrypt(originalStr)
		Convey("Encryption works without errors", func() {
			So(err, ShouldEqual, nil)
		})
		Convey("Encrypted string is different from original string", func() {
			So(originalStr, ShouldNotEqual, encryptedStr)
		})
		Convey("Decryption works without errors", func() {
			decryptedStr, err = Decrypt(encryptedStr)
			So(err, ShouldEqual, nil)
		})
		Convey("Decrypting the encrypted string results in the original string", func() {
			So(decryptedStr, ShouldEqual, originalStr)
		})
		Convey("Encrypting the same string again gives a different result", func() {
			secondEncryption, err := Encrypt(originalStr)
			So(err, ShouldEqual, nil)
			So(secondEncryption, ShouldNotEqual, encryptedStr)
		})
	})
}

func TestGetRandomInt(t *testing.T) {
	Convey("Random int generator is fairly random", t, func() {
		var randomInts [10]int64
		for i := 0; i < 10; i++ {
			randomInts[i] = GetRandomInt()
			for j := 0; j < i; j++ {
				So(randomInts[i], ShouldNotEqual, randomInts[j])
			}
		}
	})
}

func TestIsValidEmail(t *testing.T) {
	Convey("Valid emails are not invalidated", t, func() {
		var emailList []string = []string{
			"abc@yahoo.com",
			"absuqj@gmail.com",
		}
		emailListLen := len(emailList)
		for i := 0; i < emailListLen; i++ {
			isValid := IsValidEmail(emailList[i])
			So(isValid, ShouldEqual, true)
		}
	})
	Convey("Invalid emails are not passed", t, func() {
		var emailList []string = []string{
			//"abc@yahoo",
			"absuqjgmail.com",
		}
		emailListLen := len(emailList)
		for i := 0; i < emailListLen; i++ {
			isValid := IsValidEmail(emailList[i])
			So(isValid, ShouldEqual, false)
		}
	})
}

func TestIsMobileBrowser(t *testing.T) {
	Convey("Mobile browsers are detected correctly", t, func() {
		var userAgents []string = []string{
			"Mozilla/5.0 (iPad; CPU OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B554a Safari/9537.53",
			"Mozilla/5.0 (iPod touch; CPU iPhone OS 7_0_4 like Mac OS X) AppleWebKit/537.51.1 (KHTML, like Gecko) Version/7.0 Mobile/11B554a Safari/9537.53",
			"Mozilla/5.0 (Linux; U; Android 2.3.6; en-gb; SCH-I589 Build/GINGERBREAD) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
			"Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; NOKIA; Lumia 710)",
		}
		userAgentsLen := len(userAgents)
		for i := 0; i < userAgentsLen; i++ {
			isMobile := IsMobileBrowser(userAgents[i])
			So(isMobile, ShouldEqual, true)
		}
	})
	Convey("Web browsers are detected correctly", t, func() {
		var userAgents []string = []string{
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36",
			"Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:26.0) Gecko/20100101 Firefox/26.0",
		}
		userAgentsLen := len(userAgents)
		for i := 0; i < userAgentsLen; i++ {
			isMobile := IsMobileBrowser(userAgents[i])
			So(isMobile, ShouldEqual, false)
		}
	})
}

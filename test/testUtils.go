package main

import (
	"../lib"
	//	"../model"
	"bytes"
	"fmt"
	"html/template"
)

func main() {
	randomInt := lib.GetRandomInt()
	fmt.Println("Got random int as ", randomInt)

	toHash := "Get hash for this string"
	hashedString := lib.GetHashOf(toHash)
	fmt.Println("Hashed String:", hashedString)

	encryptStr := "This string is to be encrypted"
	fmt.Println("About to encrypt:", encryptStr)
	encryptedStr,_ := lib.Encrypt(encryptStr)
	fmt.Println("Encrypted string:", encryptedStr)
	decryptedStr,_ := lib.Decrypt(encryptedStr)
	fmt.Println("Decrypted string:", decryptedStr)

	fmt.Println("Valid email test:", lib.IsValidEmail("abc@def.com"))
	fmt.Println("Valid email test:", lib.IsValidEmail("9ab-99c@def.com"))
	fmt.Println("Valid email test:", lib.IsValidEmail("9ab-99cdef.com"))

	var buf bytes.Buffer
	t, _ := template.ParseFiles("../email/promo.html")
	_ = t.Execute(&buf, map[string]string{"Test": "active"})
	//fmt.Println("What does the email template have? : ", string(buf.Bytes()))

	//TestLib()
	/*	var validEmails []string
		validEmails = append(validEmails, "email1@email.com")
		validEmails = append(validEmails, "email2@email.com")
		inviteMessage := model.Message{
	        Type:       model.INVITE,
	        Sender:     "abc@def.com",
	        Recepients: validEmails,
	    }
	    inviteMessage.Send()
	*/
	encrypted, err := lib.GetDKIMSignature("Helwowo etrwer")
	fmt.Println("Encrypted:", encrypted, "Err:", err)
}

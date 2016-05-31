package model

import (
	"../lib"
	"bytes"
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"html/template"
	"log"
	"strconv"
	textTemplate "text/template"
	"time"
)

const (
	_ = iota
	INVITE
)

var (
	SenderLimits map[int]int = map[int]int{
		INVITE: 500,
	}
	ReceiverLimits map[int]int = map[int]int{
		INVITE: 1,
	}
)

type Message struct {
	Type         int
	Sender       string
	SenderName   string
	ReferralLink string
	Fbid         string
	Recepients   []string
}

func (message Message) Send() {
	encodedData, err := json.Marshal(message)
	if err != nil {
		log.Println("Unable to send message", err)
		return
	}
	lib.Enqueue(message.Type, encodedData)
}

func (message Message) ProcessAndDeliver() {
	if message.Type == INVITE {
		userLimit := SenderLimits[message.Type]
		userSentCount := message.getUserSentCount()
		if userSentCount >= userLimit {
			log.Println(message.Sender, " reached daily limit. Message not sent.", message)
			return
		}

		emailCount := len(message.Recepients)
		emailsSent := 0
		var notSentList []string
		notSentCount := 0
		currentReceiverCounts := message.getCurrentReceiverCounts()
		receiverLimit := ReceiverLimits[message.Type]
		index := 0

		for ; userSentCount < userLimit && index < emailCount; index++ {
			//TODO: Move this SQL query outside and get data for all recepients who are not members
			receiver := GetUserByEmail(message.Recepients[index])
			if receiver.Id != 0 {
				log.Println(receiver.Email, "Already a member. Message not sent.")
				continue
			}

			if currentReceiverCounts[index] >= receiverLimit {
				notSentList = append(notSentList, message.Recepients[index])
				notSentCount++
				continue
			}

			SendInviteEmail(message.Recepients[index], message.SenderName, message.ReferralLink, message.Fbid)

			currentReceiverCounts[index]++
			userSentCount++
			emailsSent++

			//One email every x seconds
			time.Sleep(time.Duration(2) * time.Second)
		}

		for ; index < emailCount; index++ {
			notSentList = append(notSentList, message.Recepients[index])
			notSentCount++
		}

		if notSentCount > 0 {
			log.Println("Sender:", message.Sender, " limit reached. Emails not sent to:", notSentList)
		}
		if emailsSent > 0 {
			message.updateReceiverCounts(currentReceiverCounts)
			message.updateSentCount(userSentCount)
		}
	}
}

func (message Message) getSenderCacheKey() string {
	msgType := strconv.Itoa(message.Type)
	return lib.GetPrefixedKey("sent_" + message.Sender + "_" + msgType)
}

func (message Message) getAllReceiverCacheKeys() []string {
	receiverCount := len(message.Recepients)
	cacheKeys := make([]string, receiverCount)
	msgType := strconv.Itoa(message.Type)
	for index := 0; index < receiverCount; index++ {
		cacheKeys[index] = lib.GetPrefixedKey("rcvd_" + message.Recepients[index] + "_" + msgType)
	}
	return cacheKeys
}

func (message Message) getUserSentCount() int {
	key := message.getSenderCacheKey()
	cache := lib.GetMCConnection()
	data, err := cache.Get(key)
	if err != nil {
		return 0
	}

	sentCountStr := string(data.Value)
	sentCount, _ := strconv.Atoi(sentCountStr)
	return sentCount
}

func (message Message) getCurrentReceiverCounts() []int {
	numReceivers := len(message.Recepients)
	receiverCounts := make([]int, numReceivers)
	cacheKeys := message.getAllReceiverCacheKeys()
	cache := lib.GetMCConnection()
	multiData, err := cache.GetMulti(cacheKeys)
	if err != nil {
		return receiverCounts
	}
	for index := 0; index < numReceivers; index++ {
		if multiData[cacheKeys[index]] != nil {
			data := multiData[cacheKeys[index]]
			tempCount := string(data.Value)
			receiverCounts[index], _ = strconv.Atoi(tempCount)
		}
	}
	return receiverCounts
}

func (message Message) updateReceiverCounts(receivedCount []int) {
	numReceivers := len(message.Recepients)
	cacheKeys := message.getAllReceiverCacheKeys()
	cache := lib.GetMCConnection()
	//No SetMulti in the bradfitz gomemcache library
	for index := 0; index < numReceivers; index++ {
		rcvdCountStr := strconv.Itoa(receivedCount[index])
		data := memcache.Item{
			Key:        cacheKeys[index],
			Value:      []byte(rcvdCountStr),
			Expiration: 86400,
		}
		_ = cache.Set(&data)
	}
}

func (message Message) updateSentCount(sentCount int) {
	key := message.getSenderCacheKey()
	cache := lib.GetMCConnection()
	sentCountStr := strconv.Itoa(sentCount)
	data := memcache.Item{
		Key:        key,
		Value:      []byte(sentCountStr),
		Expiration: 86400,
	}
	err := cache.Set(&data)
	if err != nil {
		log.Println(err)
	}
}

func SendInviteEmail(email string, senderName string, referralLink string, fbid string) {
	to := email
	subject := "I want you to checkout Jools"

	var buf bytes.Buffer
	var bufText bytes.Buffer
	t, err := template.ParseFiles("../email/invite_friend.html")
	if err != nil {
		log.Println(err)
		return
	}
	data := map[string]string{"SenderName": senderName, "ReferralLink": referralLink}
	err = t.Execute(&buf, data)
	if err != nil {
		log.Println(err)
		return
	}
	htmlMessage := string(buf.Bytes())

	txt, _ := textTemplate.ParseFiles("../email/invite_friend.txt")
	_ = txt.Execute(&bufText, data)
	txtMessage := string(bufText.Bytes())
	lib.SendEmail(htmlMessage, txtMessage, to, subject)
}

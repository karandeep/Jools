package lib

import (
	"../config"
	"labix.org/v2/mgo"
	"log"
	"strconv"
	"time"
)

type TrackData struct {
	UserId  string
	Kingdom string
	Phylum  string
	Class   string
	Order   string
	Family  string
	Genus   string
	Species string
	Date    int
	OS      string
	Browser string
	Tester  string
	Width   string
	Height  string
	IP      string
	FwdIP   string
}

type ExperimentData struct {
	UserId     string
	Experiment string
	Variant    string
	Created    int32
}

var trackingCollection *mgo.Collection
var experimentCollection *mgo.Collection

func getTrackingCollection() (*mgo.Collection, error) {
	if trackingCollection == nil {
		mongoDbConn, err := GetMongoDbConnection()
		if err != nil {
			return trackingCollection, err
		}
		configData := config.GetConfig()
		trackingCollection = mongoDbConn.DB(configData.TRACKING_DB).C(configData.TRACKING_COLLECTION_COUNTERS)
		index := mgo.Index{
			Key:        []string{"UserId", "Kingdom", "Phylum", "Class", "Date"},
			DropDups:   false,
			Background: true,
		}
		_ = trackingCollection.EnsureIndex(index)
	}
	return trackingCollection, nil
}

func getExperimentCollection() (*mgo.Collection, error) {
	if experimentCollection == nil {
		mongoDbConn, err := GetMongoDbConnection()
		if err != nil {
			return experimentCollection, err
		}
		configData := config.GetConfig()
		experimentCollection = mongoDbConn.DB(configData.TRACKING_DB).C(configData.TRACKING_COLLECTION_EXPERIMENTS)
		index := mgo.Index{
			Key:        []string{"userid", "experiment", "variant"},
			DropDups:   false,
			Background: true,
		}
		_ = experimentCollection.EnsureIndex(index)
	}
	return experimentCollection, nil
}

func GetTrackingDate() int {
	year, month, day := time.Now().Date()
	dayPadding := ""
	if day < 10 {
		dayPadding = "0"
	}
	monthPadding := ""
	if month < time.October {
		monthPadding = "0"
	}

	dateString := strconv.Itoa(year) + monthPadding + strconv.Itoa(int(month)) + dayPadding + strconv.Itoa(day)
	date, _ := strconv.Atoi(dateString)
	return date
}

func TrackCounter(trackingData TrackData, incrementBy int) (*mgo.ChangeInfo, error) {
	var info *mgo.ChangeInfo
	collection, err := getTrackingCollection()
	if err != nil {
		log.Println("Tracking failed:", trackingData, err)
		return info, err
	}

	type M map[string]interface{}
	curTimestamp := GetCurrentTimestamp()
	info, err = collection.Upsert(&trackingData, M{"$inc": M{"count": incrementBy}, "$setOnInsert": M{"created": curTimestamp}})
	if err != nil {
		log.Println("Tracking failed:", trackingData, err)
		return info, err
	}
	return info, nil
}

func UpdateId(randId string, userId string, curDate int) {
	collection, err := getTrackingCollection()
	if err != nil {
		log.Println("Tracking Id Swap Failed:", randId, userId, curDate)
		return
	}
	expCollection, err := getExperimentCollection()
	if err != nil {
		log.Println("Tracking Id Swap Failed(Exp):", randId, userId, curDate)
		return
	}

	type M map[string]interface{}
	_, err = collection.UpdateAll(M{"userid": randId, "date": curDate}, M{"$set": M{"userid": userId}})
	if err != nil {
		log.Println("Tracking Id Swap Failed:", randId, userId, curDate)
		return
	}

	_, err = expCollection.UpdateAll(M{"userid": randId}, M{"$set": M{"userid": userId}})
	if err != nil {
		log.Println("Tracking Id Swap Failed(Exp):", randId, userId, curDate)
		return
	}

}

func TrackExperiment(expData ExperimentData) {
	collection, err := getExperimentCollection()
	if err != nil {
		log.Println("Experiment tracking failed:", expData)
		return
	}

	type M map[string]interface{}
	err = collection.Insert(&expData)
	if err != nil {
		log.Println("Experiment tracking failed:", expData, err)
		return
	}
}

package lib

import (
	"../config"
	"labix.org/v2/mgo"
	"log"
)

var mongoDbConn *mgo.Session

func GetMongoDbConnection() (*mgo.Session, error) {
	if mongoDbConn == nil {
		configData := config.GetConfig()

		var info mgo.DialInfo
		info.Addrs = []string{configData.MONGO_DB_HOST + ":" + configData.MONGO_DB_PORT}
		info.Password = configData.MONGO_DB_PASSWORD
		info.Username = configData.MONGO_DB_USER
		info.Database = configData.TRACKING_DB

		var err error
		mongoDbConn, err = mgo.DialWithInfo(&info)
		if err != nil {
			log.Println("Oh no, unable to open mongo db connection")
			return mongoDbConn, err
		}

		mongoDbConn.SetMode(mgo.Monotonic, true)
	}
	return mongoDbConn, nil
}

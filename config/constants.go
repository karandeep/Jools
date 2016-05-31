package config

const (
	DEV = iota
	STAGING
	PROD_STAGE
	PRODUCTION
)

type ConfigData struct {
	GOOGLE   int
	FACEBOOK int
	YAHOO    int

	REFERRAL_BASE      string
	SOLR_BASE          string
	PARAM_FETCH_FAILED int

	SALT          string
	KEY           string
	PASSWORD_SALT string
	WEBSITE_NAME  string

	MONGO_DB_HOST     string
	MONGO_DB_PORT     string
	MONGO_DB_USER     string
	MONGO_DB_PASSWORD string

	RABBITMQ_USER     string
	RABBITMQ_PASSWORD string

	TRACK_DOMAIN                    string
	TRACKING_DB                     string
	TRACKING_COLLECTION_COUNTERS    string
	TRACKING_COLLECTION_EXPERIMENTS string

	FIVE_MIN              int
	ONE_DAY               int
	ONE_HOUR              int
	MC_KEY_PREFIX         string
	MC_DEFAULT_EXPIRATION int

	CAPTCHA_ENABLED_FOR_SIGNUP bool

	FB_PAGE          string
	FB_PAGE_ID       string
	TWTR_PAGE        string
	GOOGLE_PLUS_PAGE string
	BLOG_URL         string
	EMAIL_CONTACT    string
	MOBILE_CONTACT   string
	COMPANY_ADDRESS  string

	//Environment dependent config
	ENV               int
	BASE_URL          string
	INSPIRATIONS_URL  string
	DESIGNS_URL       string
	STATIC_URL        string
	DOMAIN            string
	FB_APP_ID         string
	FB_APP_NAME       string
	FB_APP_SECRET     string
	RECAPTCHA_PUBLIC  string
	RECAPTCHA_PRIVATE string
	DB_HOST           string
	MC_1              string
	RABBITMQ_HOST     string
	UPLOAD_URL        string
}

func GetConfig() ConfigData {
	return ConfigData{
		GOOGLE:   1,
		FACEBOOK: 2,
		YAHOO:    3,

		REFERRAL_BASE:      "https://www.jools.in",
		SOLR_BASE:          "http://ec2-54-213-159-133.us-west-2.compute.amazonaws.com:8983/solr/collection1/",
		PARAM_FETCH_FAILED: -1,

		SALT:          "Encryt!onS@ltPh0en!x",
		KEY:           "Et!onS@ltPh0en!x",
		PASSWORD_SALT: "$@lTf0rPa$$wo^d",
		WEBSITE_NAME:  "www.jools.in",

		MONGO_DB_HOST:     "ds041198.mongolab.com",
		MONGO_DB_PORT:     "41198",
		MONGO_DB_USER:     "jools-tracking",
		MONGO_DB_PASSWORD: "T^@c^Ever%th!ng0nE@rth",

		RABBITMQ_USER:     "jools-hare",
		RABBITMQ_PASSWORD: "NOt2DxxfiCuLT",

		TRACK_DOMAIN:                    "1",
		TRACKING_DB:                     "tracking",
		TRACKING_COLLECTION_COUNTERS:    "counters",
		TRACKING_COLLECTION_EXPERIMENTS: "experiments",

		FIVE_MIN:              300,
		ONE_DAY:               86400,
		ONE_HOUR:              3600,
		MC_KEY_PREFIX:         "jools_",
		MC_DEFAULT_EXPIRATION: 86400,

		CAPTCHA_ENABLED_FOR_SIGNUP: false,

		FB_PAGE:          "https://www.facebook.com/jools.in",
		FB_PAGE_ID:       "302755203197500",
		TWTR_PAGE:        "https://twitter.com/joolsandyou",
		GOOGLE_PLUS_PAGE: "https://plus.google.com/116343654856461907684",
		BLOG_URL:         "http://blog.jools.in",
		EMAIL_CONTACT:    "team@jools.in",
		MOBILE_CONTACT:   "+91-8088074745",
		COMPANY_ADDRESS:  "Phoenix Jewels Ventures Pvt. Ltd. 137, Jor Bagh, New Delhi - 110003, India",
		//PRODUCTION CONFIG
			ENV:               PRODUCTION,
			BASE_URL:          "https://www.jools.in",
			INSPIRATIONS_URL:  "https://www.jools.in/images/inspirations",
			DESIGNS_URL:       "https://www.jools.in/images/designs",
			STATIC_URL:        "https://www.jools.in",
			DOMAIN:            "www.jools.in",
			FB_APP_ID:         "372311212871698",
			FB_APP_NAME:       "my_jools",
			FB_APP_SECRET:     "0c73de9d0d7fd9b74664f3e8adedbcc0",
			RECAPTCHA_PUBLIC:  "6LdcsOYSAAAAAOQuRGzfwxTsNEB4sklQ85aSrn2E",
			RECAPTCHA_PRIVATE: "6LdcsOYSAAAAAGNPNizPP96jgW1tjhIohJmfP2w3",
			DB_HOST:           "db-1.c2tyhpuxaerw.us-west-2.rds.amazonaws.com",
			MC_1:              "mc-cluster-1.hmgjtp.cfg.usw2.cache.amazonaws.com:11211",
			RABBITMQ_HOST:     "ec2-54-213-159-133.us-west-2.compute.amazonaws.com",
			UPLOAD_URL:        "https://www.jools.in/imageServer/upload.php",
		//STAGING CONFIG
		/*
			ENV:               STAGING,
			BASE_URL:          "https://www.jools.in/staging",
			INSPIRATIONS_URL:  "https://www.jools.in/images/inspirations",
			DESIGNS_URL:       "https://www.jools.in/images/designs",
			STATIC_URL:        "https://www.jools.in/staging",
			DOMAIN:            "www.jools.in/staging",
			FB_APP_ID:         "372311212871698",
			FB_APP_NAME:       "my_jools",
			FB_APP_SECRET:     "0c73de9d0d7fd9b74664f3e8adedbcc0",
			RECAPTCHA_PUBLIC:  "6LdcsOYSAAAAAOQuRGzfwxTsNEB4sklQ85aSrn2E",
			RECAPTCHA_PRIVATE: "6LdcsOYSAAAAAGNPNizPP96jgW1tjhIohJmfP2w3",
			DB_HOST:           "db-1.c2tyhpuxaerw.us-west-2.rds.amazonaws.com",
			MC_1:              "mc-cluster-1.hmgjtp.cfg.usw2.cache.amazonaws.com:11211",
			RABBITMQ_HOST:     "ec2-54-213-159-133.us-west-2.compute.amazonaws.com",
			UPLOAD_URL:        "https://www.jools.in/staging/imageServer/upload.php",
		//DEV CONFIG
		ENV:               DEV,
		BASE_URL:          "https://192.168.1.139",
		INSPIRATIONS_URL:  "https://192.168.1.139/images/inspirations",
		DESIGNS_URL:       "https://192.168.1.139/images/designs",
		STATIC_URL:        "https://192.168.1.139",
		DOMAIN:            "192.168.1.139",
		FB_APP_ID:         "631331433558190",
		FB_APP_NAME:       "my_jools_dev",
		FB_APP_SECRET:     "ebc40e2d8f4a94c1a9a1a56418f4753a",
		RECAPTCHA_PUBLIC:  "6LekQeYSAAAAAPXnzscRgKOuKJOhsvzC1l-N97o4",
		RECAPTCHA_PRIVATE: "6LekQeYSAAAAAFRZCo0zY_EC7-7Mbo-MekBDxfZg",
		DB_HOST:           "192.168.1.139",
		MC_1:              "192.168.1.139:11211",
		RABBITMQ_HOST:     "192.168.1.139",
		UPLOAD_URL:        "https://192.168.1.139/imageServer/upload.php",
		*/
	}
}

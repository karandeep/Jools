<?php
define("GOOGLE", 1);
define("FACEBOOK", 2);
define("YAHOO",3);

if(ENV == 3) {
    define("BASE_URL","https://www.jools.in");
    define("STATIC_URL","https://www.jools.in");
    define("DOMAIN", "www.jools.in");
    define("FB_APP_ID", "372311212871698");
    define("FB_APP_SECRET","0c73de9d0d7fd9b74664f3e8adedbcc0");
    define("RECAPTCHA_PUBLIC", "6LdcsOYSAAAAAOQuRGzfwxTsNEB4sklQ85aSrn2E");
    define("RECAPTCHA_PRIVATE", "6LdcsOYSAAAAAGNPNizPP96jgW1tjhIohJmfP2w3");
    define("DB_HOST", "db-1.c2tyhpuxaerw.us-west-2.rds.amazonaws.com");
    define('MC_1','mc-cluster-1.hmgjtp.cfg.usw2.cache.amazonaws.com');
    define('RABBITMQ_HOST','ec2-54-213-159-133.us-west-2.compute.amazonaws.com');
    define('UPLOAD_URL','https://www.surajnal.in/imageServer/upload.php');
} else if (ENV == 1) {
    define("BASE_URL","https://www.jools.in/staging");
    define("STATIC_URL","https://www.jools.in/staging");
    define("DOMAIN", "www.jools.in/staging");
    define("FB_APP_ID", "372311212871698");
    define("FB_APP_SECRET","0c73de9d0d7fd9b74664f3e8adedbcc0");
    define("RECAPTCHA_PUBLIC", "6LdcsOYSAAAAAOQuRGzfwxTsNEB4sklQ85aSrn2E");
    define("RECAPTCHA_PRIVATE", "6LdcsOYSAAAAAGNPNizPP96jgW1tjhIohJmfP2w3");
    define("DB_HOST", "db-1.c2tyhpuxaerw.us-west-2.rds.amazonaws.com");
    define('MC_1','mc-cluster-1.hmgjtp.cfg.usw2.cache.amazonaws.com');
    define('RABBITMQ_HOST','ec2-54-213-159-133.us-west-2.compute.amazonaws.com');
    define('UPLOAD_URL','http://www.surajnal.in/imageServer/upload.php');
} else {
    define("BASE_URL","https://192.168.1.139");
    define("STATIC_URL","https://192.168.1.139");
    define("DOMAIN", "192.168.1.139");
    define("FB_APP_ID", "631331433558190");
    define("FB_APP_SECRET","ebc40e2d8f4a94c1a9a1a56418f4753a");
    define("RECAPTCHA_PUBLIC", "6LekQeYSAAAAAPXnzscRgKOuKJOhsvzC1l-N97o4");
    define("RECAPTCHA_PRIVATE", "6LekQeYSAAAAAFRZCo0zY_EC7-7Mbo-MekBDxfZg");
    define("DB_HOST", "192.168.1.139");
    define('MC_1','192.168.1.139');
    define('RABBITMQ_HOST','192.168.1.139');
    define('UPLOAD_URL','https://192.168.1.139/imageServer/upload.php');
}
define('INSPIRATIONS_URL','https://www.surajnal.in/imageServer/inspirations');
define("REFERRAL_BASE", "https://www.jools.in");
define("SOLR_BASE", "http://ec2-54-213-159-133.us-west-2.compute.amazonaws.com:8983/solr/collection1/");
define("PARAM_FETCH_FAILED", -1);

define('SALT', 'Encryt!onS@ltPh0en!x'); //Has a limit of 24 characters
define('PASSWORD_SALT', '$@lTf0rPa$$wo^d');
define('WEBSITE_NAME', 'jools.in');

define('MONGO_DB_HOST','ds041198.mongolab.com');
define('MONGO_DB_PORT','41198');
define('MONGO_DB_USER', 'jools-tracking');
define('MONGO_DB_PASSWORD', 'T^@c^Ever%th!ng0nE@rth');

define('RABBITMQ_USER','jools-hare');
define('RABBITMQ_PASSWORD','H@reOnTheRun%');

define('TRACK_DOMAIN','1');
define('TRACKING_DB', 'tracking');
define('TRACKING_COLLECTION_COUNTERS', 'counters');
define('COUNTER_TRACK_URL', BASE_URL . '/count');

define('FIVE_MIN', 300);
define('ONE_DAY', 86400);
define('ONE_HOUR', 3600);
define('MC_KEY_PREFIX', 'jools_');
define('MC_DEFAULT_EXPIRATION', ONE_DAY);

define('CAPTCHA_ENABLED_FOR_SIGNUP', false);

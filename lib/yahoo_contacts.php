<?php
$ch = curl_init();
curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
//$oauth_consumer_key = 'dj0yJmk9YlVVZ0ZVRjJVNzlKJmQ9WVdrOVZYcG1RVTlvTXpnbWNHbzlNQS0tJnM9Y29uc3VtZXJzZWNyZXQmeD1jZA--';
$oauth_consumer_key = 'dj0yJmk9Y0hYZTVHRDNVMmplJmQ9WVdrOVlsbHhkMUJFTmpRbWNHbzlPRFE1TmprMU1qWXkmcz1jb25zdW1lcnNlY3JldCZ4PTVl';
$oauth_signature_method = 'plaintext';
$oauth_version = '1.0';
$oauth_timestamp = time();
$mt = microtime();
$rand = mt_rand();
$oauth_nonce = md5($mt . $rand);
//$oauth_signature = '0266c2a30883472bbd35222b9e4e57c0522e58e8%26';
$oauth_signature = 'f1487a709bc061350011beb50856f474911626a2%26';
    
$oauth_token = Util::getParam('oauth_token');
if($oauth_token == PARAM_FETCH_FAILED) {
    //Need to obtain request token before redirecting user
    $xoauth_lang_pref = 'en-us';
    $oauth_callback = REFERRAL_BASE . '/callback/yahoo';
    
    curl_setopt($ch, CURLOPT_URL, "https://api.login.yahoo.com/oauth/v2/get_request_token?"
            . "oauth_nonce=$oauth_nonce&oauth_timestamp=$oauth_timestamp&oauth_consumer_key=$oauth_consumer_key"
            . "&oauth_signature_method=$oauth_signature_method&oauth_signature=$oauth_signature"
            . "&oauth_version=$oauth_version&xoauth_lang_pref=$xoauth_lang_pref&oauth_callback=$oauth_callback");
    $output = curl_exec($ch);
    $result = explode('&', $output);
    $params = array();
    foreach($result as $values) {
        $param = explode('=', $values);
        $params[$param[0]] = $param[1];
    }
    $_SESSION['yahoo_oauth_token'] = $params['oauth_token'];
    $_SESSION['yahoo_oauth_token_secret'] = $params['oauth_token_secret'];
    echo '<script type="text/javascript"> window.location="'. urldecode($params['xoauth_request_auth_url']) .'"; </script>';
} else {
    $oauth_verifier = Util::getParam('oauth_verifier');
    $oauth_signature .= $_SESSION['yahoo_oauth_token_secret'];
    $requestUrl = "https://api.login.yahoo.com/oauth/v2/get_token?"
            . "oauth_consumer_key=$oauth_consumer_key&oauth_signature_method=$oauth_signature_method"
            . "&oauth_version=$oauth_version&oauth_verifier=$oauth_verifier&oauth_token=$oauth_token"
            . "&oauth_timestamp=$oauth_timestamp&oauth_nonce=$oauth_nonce"
            . "&oauth_signature=$oauth_signature";
    curl_setopt($ch, CURLOPT_URL, $requestUrl);
    $output = curl_exec($ch);
    $result = explode('&', $output);
    $params = array();
    foreach($result as $values) {
        $param = explode('=', $values);
        $params[$param[0]] = $param[1];
    }
    
    //Make the Yahoo REST API call to get the contacts
    $oauth_timestamp = time();
    $mt = microtime();
    $rand = mt_rand();
    $oauth_nonce = md5($mt . $rand);
    $oauth_signature_method = 'HMAC-SHA1';
    $oauth_token = $params['oauth_token'];
    $oauth_token_secret = $params['oauth_token_secret'];
    $yahoo_guid = $params['xoauth_yahoo_guid'];
    $apiUrl = "http://social.yahooapis.com/v1/user/$yahoo_guid/contacts";
    $queryParams = 'count=5000&format=json&'
        . 'oauth_consumer_key='.$oauth_consumer_key.'&oauth_nonce='.$oauth_nonce.'&'
        . 'oauth_signature_method='.$oauth_signature_method.'&oauth_timestamp='.$oauth_timestamp.'&'
        . 'oauth_token='.$oauth_token.'&oauth_version='.$oauth_version;
    $parts = array(
      'GET',
      $apiUrl,
      $queryParams
    );
    $parts = Util::urlencode_rfc3986($parts);
    $base_string = implode('&', $parts);
    $key = 'f1487a709bc061350011beb50856f474911626a2'. '&' . $oauth_token_secret;
    $oauth_signature = base64_encode(hash_hmac('sha1', $base_string, $key, true));
    $url = $apiUrl . '?' . $queryParams . '&oauth_signature=' . $oauth_signature;
        
    $ch1 = curl_init();
    curl_setopt($ch1, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($ch1, CURLOPT_URL, $url);
    $output = curl_exec($ch1);
    curl_close($ch1);
    $results = json_decode($output);
    $contacts = $results->contacts->contact;
    $allEmails = '';
    $allNames = '';
    $emailData = array();
    $nameData = array();
    $count = 0;
    foreach($contacts as $contact) {
        $email = $contact->fields[0]->value;
        $nameObject = $contact->fields[1]->value;
        if(!is_string($email) || !is_object($nameObject)) {
            continue;
        }
        $name = $nameObject->givenName . ' ' . $nameObject->middleName . ' ' . $nameObject->familyName;
        $name = Util::sanitizeName($name);
        $emailData[$count] = $email;
        $nameData[$count] = $name;
        $count++;
    }
    asort($nameData);
    foreach($nameData as $index => $name) {
        $allEmails .= "{$emailData[$index]};";
        $allNames .= "$name;";
	}

	/*$goUrl = "http://localhost:8081/user/storeContacts";
	$ch = curl_init();
    curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
    curl_setopt($ch, CURLOPT_URL, $goUrl);
    curl_setopt($ch, CURLOPT_POST, true);
	curl_setopt($ch, CURLOPT_COOKIE, $_COOKIE);
	curl_setopt($ch, CURLOPT_POSTFIELDS, array('network' => YAHOO, 'emails' => $allEmails, 'names' => $allNames));
    $output = curl_exec($ch1);
    curl_close($ch1);*/
    User::storeContactEmails($_SESSION['id'], YAHOO, $allEmails, $allNames);
    print '<script type="text/javascript">'
            . 'window.opener.popupClosed('.YAHOO.');'
            . 'window.close();'
            . '</script>';
}
curl_close($ch);
